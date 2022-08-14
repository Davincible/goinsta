package goinsta

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type reqOptions struct {
	// Connection is connection header. Default is "close".
	Connection string

	// Endpoint is the request path of instagram api
	Endpoint string

	// Omit API omit the /api/v1/ part of the url
	OmitAPI bool

	// IsPost set to true will send request with POST method.
	//
	// By default this option is false.
	IsPost bool

	// Compress post form data with gzip
	Gzip bool

	// UseV2 is set when API endpoint uses v2 url.
	UseV2 bool

	// Use b.i.instagram.com
	Useb bool

	// Query is the parameters of the request
	//
	// This parameters are independents of the request method (POST|GET)
	Query map[string]string

	// DataBytes can be used to pass raw data to a request, instead of a
	//   form using the Query param. This is used for e.g. photo and vieo uploads.
	DataBytes *bytes.Buffer

	// List of headers to ignore
	IgnoreHeaders []string

	// Extra headers to add
	ExtraHeaders map[string]string

	// Timestamp
	Timestamp string

	// Count the number of times the wrapper has been called.
	WrapperCount int

	// If Status 429 should be ignored, ErrTooManyRequests. This behaviour should be implemented in
	//  the wrapper. Goinsta does nothing directly with this value.
	Ignore429 bool
}

func (insta *Instagram) sendSimpleRequest(uri string, a ...interface{}) (body []byte, err error) {
	body, _, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(uri, a...),
		},
	)
	return body, err
}

func (insta *Instagram) sendRequest(o *reqOptions) (body []byte, h http.Header, err error) {
	if insta == nil {
		return nil, nil, fmt.Errorf("Error while calling %s: %s", o.Endpoint, ErrInstaNotDefined)
	}

	// Check if a challenge is in progress, if so wait for it to complete (with timeout)
	if insta.privacyRequested.Get() && !insta.privacyCalled.Get() {
		if !insta.checkPrivacy() {
			return nil, nil, errors.New("Privacy check timedout")
		}
	}

	insta.checkXmidExpiry()

	method := "GET"
	if o.IsPost {
		method = "POST"
	}
	if o.Connection == "" {
		o.Connection = "close"
	}

	if o.Timestamp == "" {
		o.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	}

	var nu string
	if o.Useb {
		nu = instaAPIUrlb
	} else {
		nu = instaAPIUrl
	}
	if o.UseV2 && !o.Useb {
		nu = instaAPIUrlv2
	} else if o.UseV2 && o.Useb {
		nu = instaAPIUrlv2b
	}
	if o.OmitAPI {
		nu = baseUrl
		o.IgnoreHeaders = append(o.IgnoreHeaders, omitAPIHeadersExclude...)
	}

	nu = nu + o.Endpoint
	u, err := url.Parse(nu)
	if err != nil {
		return nil, nil, err
	}

	vs := url.Values{}
	bf := bytes.NewBuffer([]byte{})
	reqData := bytes.NewBuffer([]byte{})

	for k, v := range o.Query {
		vs.Add(k, v)
	}

	// If DataBytes has been passed, use that as data, else use Query
	if o.DataBytes != nil {
		reqData = o.DataBytes
	} else {
		reqData.WriteString(vs.Encode())
	}

	var contentEncoding string
	if o.IsPost && o.Gzip {
		// If gzip encoding needs to be applied
		zw := gzip.NewWriter(bf)
		defer zw.Close()
		if _, err := zw.Write(reqData.Bytes()); err != nil {
			return nil, nil, err
		}
		if err := zw.Close(); err != nil {
			return nil, nil, err
		}
		contentEncoding = "gzip"
	} else if o.IsPost {
		// use post form if POST request
		bf = reqData
	} else {
		// append query to url if GET request
		for k, v := range u.Query() {
			vs.Add(k, strings.Join(v, " "))
		}

		u.RawQuery = vs.Encode()
	}

	var req *http.Request
	req, err = http.NewRequest(method, u.String(), bf)
	if err != nil {
		return
	}

	ignoreHeader := func(h string) bool {
		for _, k := range o.IgnoreHeaders {
			if k == h {
				return true
			}
		}
		return false
	}

	setHeaders := func(h map[string]string) {
		for k, v := range h {
			if v != "" && !ignoreHeader(k) {
				req.Header.Set(k, v)
			}
		}
	}
	setHeadersAsync := func(key, value interface{}) bool {
		k, v := key.(string), value.(string)
		if v != "" && !ignoreHeader(k) {
			req.Header.Set(k, v)
		}
		return true
	}

	headers := map[string]string{
		"Accept-Language":             locale,
		"Accept-Encoding":             "gzip,deflate",
		"Connection":                  o.Connection,
		"Content-Type":                "application/x-www-form-urlencoded; charset=UTF-8",
		"User-Agent":                  insta.userAgent,
		"X-Ig-App-Locale":             locale,
		"X-Ig-Device-Locale":          locale,
		"X-Ig-Mapped-Locale":          locale,
		"X-Ig-App-Id":                 fbAnalytics,
		"X-Ig-Device-Id":              insta.uuid,
		"X-Ig-Family-Device-Id":       insta.fID,
		"X-Ig-Android-Id":             insta.dID,
		"X-Ig-Timezone-Offset":        timeOffset,
		"X-Ig-Capabilities":           igCapabilities,
		"X-Ig-Connection-Type":        connType,
		"X-Pigeon-Session-Id":         insta.psID,
		"X-Pigeon-Rawclienttime":      fmt.Sprintf("%s.%d", o.Timestamp, random(100, 900)),
		"X-Ig-Bandwidth-Speed-KBPS":   fmt.Sprintf("%d.000", random(1000, 9000)),
		"X-Ig-Bandwidth-TotalBytes-B": strconv.FormatInt(random(1000000, 5000000), 10),
		"X-Ig-Bandwidth-Totaltime-Ms": strconv.FormatInt(random(200, 800), 10),
		"X-Ig-App-Startup-Country":    "unkown",
		"X-Bloks-Version-Id":          bloksVerID,
		"X-Bloks-Is-Layout-Rtl":       "false",
		"X-Bloks-Is-Panorama-Enabled": "true",
		"X-Fb-Http-Engine":            "Liger",
		"X-Fb-Client-Ip":              "True",
		"X-Fb-Server-Cluster":         "True",
	}
	if insta.Account != nil {
		headers["Ig-Intended-User-Id"] = strconv.FormatInt(insta.Account.ID, 10)
	} else {
		headers["Ig-Intended-User-Id"] = "0"
	}
	if contentEncoding != "" {
		headers["Content-Encoding"] = contentEncoding
	}

	setHeaders(headers)
	setHeaders(o.ExtraHeaders)
	insta.headerOptions.Range(setHeadersAsync)

	resp, err := insta.c.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	insta.extractHeaders(resp.Header)
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// Extract error from request body, if present
	err = insta.isError(resp.StatusCode, body, resp.Status, o.Endpoint)

	// Decode gzip encoded responses
	encoding := resp.Header.Get("Content-Encoding")
	if encoding != "" && encoding == "gzip" {
		buf := bytes.NewBuffer(body)
		zr, err := gzip.NewReader(buf)
		if err != nil {
			return nil, nil, err
		}
		body, err = ioutil.ReadAll(zr)
		if err != nil {
			return nil, nil, err
		}
		if err := zr.Close(); err != nil {
			return nil, nil, err
		}
	}

	// Log complete response body
	if insta.Debug {
		r := map[string]interface{}{
			"status":   resp.StatusCode,
			"endpoint": o.Endpoint,
			"body":     string(body),
		}

		b, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		insta.debugHandler(string(b))
	}

	// Call Request Wrapper
	hCopy := resp.Header.Clone()
	if insta.reqWrapper != nil {
		o.WrapperCount += 1
		body, hCopy, err = insta.reqWrapper.GoInstaWrapper(
			&ReqWrapperArgs{
				insta:      insta,
				reqOptions: o,
				Body:       body,
				Headers:    hCopy,
				Error:      err,
			})
	}

	return body, hCopy, err
}

func (insta *Instagram) checkXmidExpiry() {
	t := time.Now().Unix()
	if insta.xmidExpiry != -1 && t > insta.xmidExpiry-10 {
		insta.xmidExpiry = -1
		insta.zrToken()
	}
}

func (insta *Instagram) extractHeaders(h http.Header) {
	extract := func(in string, out string) {
		x := h[in]
		if len(x) > 0 && x[0] != "" {
			// prevent from auth being set without token post login
			if in == "Ig-Set-Authorization" {
				old, ok := insta.headerOptions.Load(out)
				if ok && len(old.(string)) != 0 {
					current := strings.Split(old.(string), ":")
					newHeader := strings.Split(x[0], ":")
					if len(current[2]) > len(newHeader[2]) {
						return
					}

				}
			}
			insta.headerOptions.Store(out, x[0])
		}
	}

	extract("Ig-Set-Authorization", "Authorization")
	extract("Ig-Set-X-Mid", "X-Mid")
	extract("X-Ig-Set-Www-Claim", "X-Ig-Www-Claim")
	extract("Ig-Set-Ig-U-Ig-Direct-Region-Hint", "Ig-U-Ig-Direct-Region-Hint")
	extract("Ig-Set-Ig-U-Shbid", "Ig-U-Shbid")
	extract("Ig-Set-Ig-U-Shbts", "Ig-U-Shbts")
	extract("Ig-Set-Ig-U-Rur", "Ig-U-Rur")
	extract("Ig-Set-Ig-U-Ds-User-Id", "Ig-U-Ds-User-Id")
}

func (insta *Instagram) checkPrivacy() bool {
	d := time.Now().Add(5 * time.Minute)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	for {
		select {
		case <-time.After(3 * time.Second):
			if insta.privacyRequested.Get() {
				return true
			}
		case <-ctx.Done():
			return false
		}
	}
}

func (insta *Instagram) isError(code int, body []byte, status, endpoint string) (err error) {
	switch code {
	case 200:
	case 202:
	case 400:
		ierr := Error400{Endpoint: endpoint}

		// Ignore error, doesn't matter if types don't always match up
		json.Unmarshal(body, &ierr)

		switch ierr.GetMessage() {
		case "login_required":
			if ierr.ErrorTitle == "You've Been Logged Out" {
				return ErrLoggedOut
			}
			return ErrLoginRequired
		case "bad_password":
			return ErrBadPassword

		case "Sorry, this media has been deleted":
			return ErrMediaDeleted

		case "checkpoint_required":
			// Usually a request to accept cookies
			insta.warnHandler(ierr)
			insta.Checkpoint = &ierr.Checkpoint
			insta.Checkpoint.insta = insta
			return ErrCheckpointRequired

		case "checkpoint_challenge_required":
			fallthrough
		case "challenge_required":
			insta.warnHandler(ierr)
			insta.Challenge = ierr.Challenge
			insta.Challenge.insta = insta
			return ErrChallengeRequired

		case "two_factor_required":
			insta.TwoFactorInfo = ierr.TwoFactorInfo
			insta.TwoFactorInfo.insta = insta
			if insta.Account == nil {
				insta.Account = &Account{
					ID:       insta.TwoFactorInfo.ID,
					Username: insta.TwoFactorInfo.Username,
				}
			}
			return Err2FARequired
		case "Please check the code we sent you and try again.":
			return ErrInvalidCode

		default:
			return ierr
		}

	case 403:
		ierr := Error400{
			Code:     403,
			Endpoint: endpoint,
		}
		err = json.Unmarshal(body, &ierr)
		if err == nil {
			switch ierr.Message {
			case "login_required":
				if ierr.ErrorTitle == "You've Been Logged Out" {
					return ErrLoggedOut
				}
				return ErrLoginRequired
			}
			return ierr
		}
		return err
	case 429:
		return ErrTooManyRequests
	case 500:
		ierr := ErrorN{
			Endpoint:  endpoint,
			Status:    "500",
			Message:   string(body),
			ErrorType: status,
		}
		return ierr
	case 503:
		return Error503{
			Message: "Instagram API error. Try it later.",
		}
	default:
		ierr := ErrorN{
			Endpoint:  endpoint,
			Status:    strconv.Itoa(code),
			Message:   string(body),
			ErrorType: status,
		}
		err = json.Unmarshal(body, &ierr)
		if ierr.Message == "Transcode not finished yet." {
			return nil
		}
		return ierr
	}
	return nil
}

func (insta *Instagram) prepareData(other ...map[string]interface{}) (string, error) {
	data := map[string]interface{}{
		"_uuid": insta.uuid,
	}
	if insta.Account != nil && insta.Account.ID != 0 {
		data["_uid"] = strconv.FormatInt(insta.Account.ID, 10)
	}

	for i := range other {
		for key, value := range other[i] {
			data[key] = value
		}
	}
	b, err := json.Marshal(data)
	if err == nil {
		return string(b), err
	}
	return "", err
}

func (insta *Instagram) prepareDataQuery(other ...map[string]interface{}) map[string]string {
	data := map[string]string{
		"_uuid":     insta.uuid,
		"device_id": insta.dID,
	}
	for i := range other {
		for key, value := range other[i] {
			data[key] = toString(value)
		}
	}
	return data
}

func random(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}
