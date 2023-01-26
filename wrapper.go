package goinsta

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	TooManyRequestsTimeout = 60 * time.Second
)

type ReqWrapper interface {
	GoInstaWrapper(*ReqWrapperArgs) (body []byte, h http.Header, err error)
}

type ReqWrapperArgs struct {
	insta      *Instagram
	reqOptions *reqOptions

	Body    []byte
	Headers http.Header
	Error   error
}

type Wrapper struct {
	o *ReqWrapperArgs
}

func (w *ReqWrapperArgs) RetryRequest() (body []byte, h http.Header, err error) {
	return w.insta.sendRequest(w.reqOptions)
}

func (w *ReqWrapperArgs) GetWrapperCount() int {
	return w.reqOptions.WrapperCount
}

func (w *ReqWrapperArgs) GetInsta() *Instagram {
	return w.insta
}

func (w *ReqWrapperArgs) GetEndpoint() string {
	return w.reqOptions.Endpoint
}

func (w *ReqWrapperArgs) SetInsta(insta *Instagram) {
	w.insta = insta
}

func (w *ReqWrapperArgs) Ignore429() bool {
	return w.reqOptions.Ignore429
}

func DefaultWrapper() *Wrapper {
	return &Wrapper{}
}

// GoInstaWrapper is a warpper function for goinsta
func (w *Wrapper) GoInstaWrapper(o *ReqWrapperArgs) ([]byte, http.Header, error) {
	// If no errors occured, directly return
	if o.Error == nil {
		return o.Body, o.Headers, o.Error
	}

	// If wrapper called more than threshold, return
	if o.GetWrapperCount() > 3 {
		return o.Body, o.Headers, o.Error
	}

	w.o = o
	insta := o.GetInsta()

	switch true {
	case errors.Is(o.Error, ErrTooManyRequests):
		// Some endpoints often return 429, too many requests, and can be safely ignored.
		if o.Ignore429() {
			return o.Body, o.Headers, nil
		}
			insta.warnHandler("Too many requests, sleeping for %d seconds", TooManyRequestsTimeout)
			time.Sleep(TooManyRequestsTimeout)

	case errors.Is(o.Error, Err2FARequired):
		// Attempt auto 2FA login with TOTP code generation
		err := insta.TwoFactorInfo.Login2FA()
		if err != nil && err != Err2FANoCode {
			return o.Body, o.Headers, err
		} else {
			return o.Body, o.Headers, o.Error
		}

	case errors.Is(o.Error, ErrLoggedOut):
		fallthrough
	case errors.Is(o.Error, ErrLoginRequired):
		return o.Body, o.Headers, o.Error

	case errors.Is(o.Error, ErrCheckpointRequired):
		// Attempt to accecpt cookies using headless browser
		err := insta.Checkpoint.Process()
		if err != nil {
			return o.Body, o.Headers, fmt.Errorf(
				"failed to automatically process status code 400 'checkpoint_required' with checkpoint url '%s', please report this on github. Error provided: %w",
				insta.Checkpoint.URL,
				err,
			)
		}
		insta.infoHandler(
			fmt.Sprintf("Auto solving of checkpoint with url '%s' seems to have gone successful. This is an experimental feature, please let me know if it works! :)\n",
				insta.Checkpoint.URL,
			))

	case errors.Is(o.Error, ErrCheckpointPassed):
		// continue without doing anything, retry request

	case errors.Is(o.Error, ErrChallengeRequired):
		if err := insta.Challenge.Process(); err != nil {
			return o.Body, o.Headers, fmt.Errorf("failed to process challenge automatically with: %w", err)
		}
	default:
		// Unhandeled errors should be passed on
		return o.Body, o.Headers, o.Error
	}

	body, h, err := w.o.RetryRequest()
	return body, h, err
}
