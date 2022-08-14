package utilities

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"strconv"
	"strings"
	"time"
)

// prefix will append extra 0s if the length of otp is less than 6.
func prefix(otp string) string {
	if len(otp) == 6 {
		return otp
	}
	for i := (6 - len(otp)); i > 0; i-- {
		otp = "0" + otp
	}
	return otp
}

// genHOTP will generate a Hmac One Time Password
func genHOTP(secret string, interval int64) (string, error) {
	// Decode base32 secret
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", err
	}

	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(interval))

	// Signing the value using HMAC-SHA1 Algorithm
	hash := hmac.New(sha1.New, key)
	hash.Write(bs)
	h := hash.Sum(nil)

	// Grab offset
	o := (h[19] & 15)

	// Get 32 bit chunk from hash starting at the o
	var header uint32
	r := bytes.NewReader(h[o : o+4])
	err = binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return "", err
	}

	// Ignore most significant bits as per RFC 4226, and crop to 6 digits.
	h12 := (int(header) & 0x7fffffff) % 1000000

	otp := strconv.Itoa(h12)
	otp = prefix(otp)

	return otp, nil
}

// GenTOTP will generate a one time pass based on the secret.
func GenTOTP(secret string) (string, error) {
	//The TOTP token is just a HOTP token seeded with every 30 seconds.
	interval := time.Now().Unix() / 30
	otp, err := genHOTP(secret, interval)
	return otp, err
}
