package goinsta

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Davincible/goinsta/v3/utilities"
)

type TwoFactorInfo struct {
	insta *Instagram

	ID       int64  `json:"pk"`
	Username string `json:"username"`

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

	PhoneVerificationSettings phoneVerificationSettings `json:"phone_verification_settings"`
}

type phoneVerificationSettings struct {
	MaxSMSCount          int  `json:"max_sms_count"`
	ResendSMSDelaySec    int  `json:"resend_sms_delay_sec"`
	RobocallAfterMaxSms  bool `json:"robocall_after_max_sms"`
	RobocallCountDownSec int  `json:"robocall_count_down_time_sec"`
}

// Login2FA allows for a login through 2FA
// You can either provide a code directly by passing it as a parameter, or
//  goinsta can generate one for you as long as the TOTP seed is set.
func (info *TwoFactorInfo) Login2FA(in ...string) error {
	insta := info.insta

	var code string
	if len(in) > 0 {
		code = in[0]
	} else if info.insta.totp == nil || info.insta.totp.Seed == "" {
		return Err2FANoCode
	} else {
		otp, err := utilities.GenTOTP(insta.totp.Seed)
		if err != nil {
			return fmt.Errorf("Failed to generate 2FA OTP code: %w", err)
		}
		code = otp
	}

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
			"verification_method":   "3",
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
			IgnoreHeaders: []string{
				"Ig-U-Shbts",
				"Ig-U-Shbid",
				"Ig-U-Rur",
				"Authorization",
			},
			ExtraHeaders: map[string]string{
				"X-Ig-Www-Claim":      "0",
				"Ig-Intended-User-Id": "0",
			},
		},
	)
	if err != nil {
		return err
	}

	err = insta.parseLogin(body)
	if err != nil {
		return err
	}

	err = insta.OpenApp()
	return err
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
