package goinsta

import (
	"encoding/json"
	"strings"
)

type ChallengeStepData struct {
	Choice           string      `json:"choice"`
	FbAccessToken    string      `json:"fb_access_token"`
	BigBlueToken     string      `json:"big_blue_token"`
	GoogleOauthToken string      `json:"google_oauth_token"`
	Email            string      `json:"email"`
	SecurityCode     string      `json:"security_code"`
	ResendDelay      interface{} `json:"resend_delay"`
	ContactPoint     string      `json:"contact_point"`
	FormType         string      `json:"form_type"`
}

// Challenge is a status code 400 error, usually prompting the user to perform
//   some action.
type Challenge struct {
	insta *Instagram

	LoggedInUser *Account `json:"logged_in_user,omitempty"`
	UserID       int64    `json:"user_id"`
	Status       string   `json:"status"`

	Errors []string `json:"errors"`

	URL               string            `json:"url"`
	ApiPath           string            `json:"api_path"`
	Context           *ChallengeContext `json:"challenge_context"`
	FlowRenderType    int               `json:"flow_render_type"`
	HideWebviewHeader bool              `json:"hide_webview_header"`
	Lock              bool              `json:"lock"`
	Logout            bool              `json:"logout"`
	NativeFlow        bool              `json:"native_flow"`

	TwoFactorRequired bool
	TwoFactorInfo     TwoFactorInfo
}

// Checkpoint is just like challenge, a status code 400 error.
// Usually used to prompt the user to accept cookies.
type Checkpoint struct {
	insta *Instagram

	URL            string `json:"checkpoint_url"`
	Lock           bool   `json:"lock"`
	FlowRenderType int    `json:"flow_render_type"`
}

type ChallengeContext struct {
	TypeEnum    string            `json:"challenge_type_enum"`
	IsStateless bool              `json:"is_stateless"`
	Action      string            `json:"action"`
	NonceCode   string            `json:"nonce_code"`
	StepName    string            `json:"step_name"`
	StepData    ChallengeStepData `json:"step_data"`
	UserID      int64             `json:"user_id"`
}

type challengeResp struct {
	*Challenge
}

func newChallenge(insta *Instagram) *Challenge {
	return &Challenge{
		insta: insta,
	}
}

// updateState updates current data from challenge url
func (c *Challenge) updateState() error {
	insta := c.insta

	ctx, err := json.Marshal(c.Context)
	if err != nil {
		return err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: c.insta.challengeURL,
			Query: map[string]string{
				"guid":              insta.uuid,
				"device_id":         insta.dID,
				"challenge_context": string(ctx),
			},
		},
	)
	if err == nil {
		resp := challengeResp{}
		err = json.Unmarshal(body, &resp)
		if err == nil {
			*c = *resp.Challenge
			c.insta = insta
		}
	}
	return err
}

// selectVerifyMethod selects a way and verify it (Phone number = 0, email = 1)
func (challenge *Challenge) selectVerifyMethod(choice string, isReplay ...bool) error {
	insta := challenge.insta

	url := challenge.insta.challengeURL
	if len(isReplay) > 0 && isReplay[0] {
		url = strings.Replace(url, "/challenge/", "/challenge/replay/", -1)
	}

	data, err := json.Marshal(
		map[string]string{
			"choice":    choice,
			"guid":      insta.uuid,
			"device_id": insta.dID,
			"_uuid":     insta.uuid,
			"_uid":      toString(insta.Account.ID),
		})
	if err != nil {
		return err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: url,
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err == nil {
		resp := challengeResp{}
		err = json.Unmarshal(body, &resp)
		if err == nil {
			*challenge = *resp.Challenge
			challenge.insta = insta
		}
	}
	return err
}

// sendSecurityCode sends the code received in the message
func (challenge *Challenge) SendSecurityCode(code string) error {
	insta := challenge.insta
	url := challenge.insta.challengeURL

	data, err := json.Marshal(map[string]string{
		"security_code": code,
		"guid":          insta.uuid,
		"device_id":     insta.dID,
		"_uuid":         insta.uuid,
		"_uid":          toString(insta.Account.ID),
	})
	if err != nil {
		return err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: url,
			IsPost:   true,
			Query:    generateSignature(data),
		},
	)
	if err == nil {
		resp := challengeResp{}
		err = json.Unmarshal(body, &resp)
		if err == nil {
			*challenge = *resp.Challenge
			challenge.insta = insta
			challenge.LoggedInUser.insta = insta
			insta.Account = challenge.LoggedInUser
		}
	}
	return err
}

// deltaLoginReview process with choice (It was me = 0, It wasn't me = 1)
func (c *Challenge) deltaLoginReview() error {
	return c.selectVerifyMethod("0")
}

func (c *Challenge) ProcessOld(apiURL string) error {
	c.insta.challengeURL = apiURL[1:]

	if err := c.updateState(); err != nil {
		return err
	}

	switch c.Context.StepName {
	case "select_verify_method":
		return c.selectVerifyMethod(c.Context.StepData.Choice)
	case "delta_login_review":
		return c.deltaLoginReview()
	}

	return ErrChallengeProcess{StepName: c.Context.StepName}
}

// Process will open up the challenge url in a chromium browser and
//   take a screenshot. Please report the screenshot and printed out struct so challenge
//   automation can be build in.
func (c *Challenge) Process() error {
	insta := c.insta

	insta.warnHandler("Encountered a captcha challenge, goinsta will attempt to open the challenge in a headless chromium browser, and take a screenshot. Please report the details in a github issue.")
	err := insta.openChallenge(c.URL)
	err = checkHeadlessErr(err)

	return err
}

// Process will open up the url passed as a checkpoint response (not a challenge)
//   in a headless browser. This method is experimental, please report if you still
//   get a /privacy/checks/ checkpoint error.
func (c *Checkpoint) Process() error {
	insta := c.insta
	if insta.privacyRequested.Get() {
		panic("Privacy request again, it hus failed, panicing")
	}

	insta.privacyRequested.Set(true)
	err := insta.acceptPrivacyCookies(c.URL)
	err = checkHeadlessErr(err)
	if err != nil {
		return err
	}

	insta.privacyCalled.Set(true)
	return nil
}
