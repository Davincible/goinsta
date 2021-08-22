package goinsta

import (
	"encoding/json"
	"errors"
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

type Challenge struct {
	insta *Instagram

	LoggedInUser *Account `json:"logged_in_user,omitempty"`
	UserID       int64    `json:"user_id"`
	Status       string   `json:"status"`

	ApiPath           string            `json:"api_path"`
	Context           *ChallengeContext `json:"challenge_context"`
	FlowRenderType    int               `json:"flow_render_type"`
	HideWebviewHeader bool              `json:"hide_webview_header"`
	Lock              bool              `json:"lock"`
	Logout            bool              `json:"logout"`
	NativeFlow        bool              `json:"native_flow"`
	URL               string            `json:"url"`

	TwoFactorRequired bool
	TwoFactorInfo     TwoFactorInfo
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

type TwoFactorInfo struct {
	insta *Instagram

	ElegibleForMultipleTotp    bool   `json:"elegible_for_multiple_totp"`
	ObfuscatedPhoneNr          string `json:"obfuscated_phone_number"`
	PendingTrustedNotification bool   `json:"pending_trusted_notification"`
	ShouldOptInTrustedDevice   bool   `json:"should_opt_in_trusted_device_option"`
	ShowMessengerCodeOption    bool   `json:"show_messenger_code_option"`
	ShowTrustedDeviceOption    bool   `json:"show_trusted_device_option"`
	SMSNotAllowedReason        string `json:"sms_not_allowed_reason"`
	SMSTwoFactorOn             bool   `json:"sms_two_factor_on"`
	TotpTwoFactorOn            bool   `json:"totp_two_factor_on"`
	WhatsappTwoFactorOn        bool   `json:"whatsapp_two_factor_on"`
	TwoFactorIdentifier        string `json:"two_factor_identifier"`
	Username                   string `json:"username"`

	PhoneVerificationSettings phoneVerificationSettings `json:"phone_verification_settings"`
}

type phoneVerificationSettings struct {
	MaxSMSCount          int  `json:"max_sms_count"`
	ResendSMSDelaySec    int  `json:"resend_sms_delay_sec"`
	RobocallAfterMaxSms  bool `json:"robocall_after_max_sms"`
	RobocallCountDownSec int  `json:"robocall_count_down_time_sec"`
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

func (c *Challenge) Process(apiURL string) error {
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

// Check2FATrusted checks whether the device has been trusted.
// When you enable 2FA, you can verify, or trust, the device with one of your
//   other devices. This method will check if this device has been trusted.
// if so, it will login, if not, it will return an error.
// The android app calls this method every 3 seconds
func (info *TwoFactorInfo) Check2FATrusted() error {
	insta := info.insta
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: url2FACheckTrusted,
			Query: map[string]string{
				"two_factor_identifier": info.TwoFactorIdentifier,
				"username":              insta.user,
				"device_id":             insta.dID,
			},
		},
	)
	if err != nil {
		return err
	}

	stat := struct {
		ReviewStatus int    `json:"review_status"`
		Status       string `json:"status"`
	}{}
	err = json.Unmarshal(body, &stat)
	if stat.ReviewStatus == 0 {
		return errors.New("Two factor authentication not yet verified")
	}

	err = info.Login2FA("")
	return err
}

// Login2FA allows for a login through 2FA
func (info *TwoFactorInfo) Login2FA(code string) error {
	insta := info.insta
	data, err := json.Marshal(
		map[string]string{
			"verification_code":     code,
			"phone_id":              insta.fID,
			"two_factor_identifier": info.TwoFactorIdentifier,
			"username":              insta.user,
			"trust_this_device":     "1",
			"guid":                  insta.uuid,
			"device_id":             insta.dID,
			"waterfall_id":          generateUUID(),
			"verification_method":   "4",
		},
	)
	if err != nil {
		return err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: url2FALogin,
			IsPost:   true,
			Query:    generateSignature(data),
		},
	)
	if err != nil {
		return err
	}

	err = insta.verifyLogin(body)
	if err != nil {
		return err
	}

	err = insta.OpenApp()
	return err
}

func (c *ChallengeError) Process() ([]byte, error) {
	insta := c.insta

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: c.Challenge.URL,
			IsPost:   true,
			Query: map[string]string{
				"guid":      insta.uuid,
				"device_id": insta.dID,
			},
		},
	)
	return body, err
}
