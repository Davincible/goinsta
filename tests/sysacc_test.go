package tests

import (
	"bytes"
	"encoding/base64"
	"errors"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/Davincible/goinsta"
)

var errNoValidLogin = errors.New("No valid login found")

func readFromBase64(base64EncodedString string) (*goinsta.Instagram, error) {
	base64Bytes, err := base64.StdEncoding.DecodeString(base64EncodedString)
	if err != nil {
		return nil, err
	}
	return goinsta.ImportReader(bytes.NewReader(base64Bytes))
}

func availableEncodedAccounts() ([]string, error) {
	output := make([]string, 0)

	environ, err := loadEnv()
	if err != nil {
		return nil, err
	}

	for _, env := range environ {
		if strings.HasPrefix(env, "INSTAGRAM_BASE64_") {
			index := strings.Index(env, "=")
			encodedString := env[index+1:]
			if encodedString[0] == '"' {
				encodedString = encodedString[1 : len(encodedString)-1]
			}
			output = append(output, encodedString)
		}
	}

	return output, nil
}

func TestGetRandomAccount(t *testing.T) {
	for i := 0; i < 50; i++ {
		insta, err := getRandomAccount(t)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(insta.Account.Username)
	}
}

func TestGetRandomLogin(t *testing.T) {
	for i := 0; i < 50; i++ {
		uname, pw, err := getLogin()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(uname, pw)
	}
}

func getRandomAccount(ts ...*testing.T) (*goinsta.Instagram, error) {
	accounts, err := availableEncodedAccounts()
	if err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, errNoValidLogin
	}

	r := rand.Intn(len(accounts))
	encodedAccount := accounts[r]
	insta, err := readFromBase64(encodedAccount)
	if err != nil {
		return nil, err
	}

	// if t != nil {
	// 	t.Logf("Found %d accounts, rand n = %d aka '%s' %p, %s, %+v\n\"%s\"\n",
	// 		len(accounts),
	// 		r,
	// 		insta.Account.Username,
	// 		insta.Account,
	// 		insta.Account.FullName,
	// 		insta.Account,
	// 		encodedAccount,
	// 	)
	// }
	return insta, err
}

func getLogin() (string, string, error) {
	var accs [][]string

	environ, err := loadEnv()
	if err != nil {
		return "", "", nil
	}

	for _, env := range environ {
		if strings.HasPrefix(env, "INSTAGRAM_ACT_") {
			acc := strings.Split(env, "=")[1]
			if acc[0] == '"' {
				acc = acc[1 : len(acc)-1]
			}

			accs = append(accs, strings.Split(acc, ":"))
		}
	}

	if len(accs) == 0 {
		return "", "", errNoValidLogin
	}

	for {
		i := rand.Intn(len(accs))
		r := accs[i]
		if len(r) == 2 {
			return r[0], r[1], nil
		}
		accs = append(accs[:i], accs[i+1:]...)
		if len(accs) == 0 {
			return "", "", errNoValidLogin
		}
	}
}

func loadEnv() ([]string, error) {
	environ := os.Environ()
	env, err := dotenv()
	if err != nil {
		return nil, err
	}
	environ = append(environ, env...)
	return environ, nil
}

func dotenv() ([]string, error) {
	file, err := os.Open(".env")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, err
	}
	env := strings.Split(string(buf.Bytes()), "\n")
	return env, nil
}
