package goinsta

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

const (
	volatileSeed = "12345"
)

func generateMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func generateHMAC(text, key string) string {
	hasher := hmac.New(sha256.New, []byte(key))
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func generateDeviceID(seed string) string {
	hash := generateMD5Hash(seed + volatileSeed)
	return "android-" + hash[:16]
}

func generateUserBreadcrumb(text string) string {
	ts := time.Now().Unix()
	d := fmt.Sprintf("%d %d %d %d%d",
		len(text), 0, random(3000, 10000), ts, random(100, 999))
	hmac := base64.StdEncoding.EncodeToString([]byte(generateHMAC(d, hmacKey)))
	enc := base64.StdEncoding.EncodeToString([]byte(d))
	return hmac + "\n" + enc + "\n"
}

// generateSignature takes a string of byte slice as argument, and prepents the signature
func generateSignature(d interface{}, extra ...map[string]string) map[string]string {
	var data string
	switch x := d.(type) {
	case []byte:
		data = string(x)
	case string:
		data = x
	}
	r := map[string]string{
		"signed_body": "SIGNATURE." + data,
	}
	for _, e := range extra {
		for k, v := range e {
			r[k] = v
		}
	}

	return r
}

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func generateUUID() string {
	uuid, err := newUUID()
	if err != nil {
		return "cb479ee7-a50d-49e7-8b7b-60cc1a105e22" // default value when error occurred
	}
	return uuid
}
