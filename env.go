package goinsta

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

type plainAcc struct {
	Name     string
	Username string
	Password string
}

type encAcc struct {
	Username string
	Line     string
	Base64   string
}

// EnvRandAcc will check the environment variables, and the .env file in
//   the current working directory (unless another path has been provided),
//   for either a base64 encoded goinsta config, or plain credentials.
//
// To use this function, add one or multiple of the following:
//   INSTAGRAM_ACT_<name>="username:password"
//   INSTAGRAM_BASE64_<name>="<base64 encoded config>"
//
// INSTAGRAM_ACT_ variables will automatiaclly be converted to INSTAGRAM_BASE64_
//
func EnvRandAcc(path ...string) (*Instagram, error) {
	err := checkEnv(path...)
	if err != nil {
		return nil, err
	}
	return getRandAcc(path...)
}

// EnvRandLogin fetches a random login from the env.
// :param: path (OPTIONAL) - path to a file, by default .env
//
// Looks for INSTAGRAM_ACT_<name>="username:password" in env
//
func EnvRandLogin(path ...string) (string, string, error) {
	accs, err := avPlainAcc(path...)
	if err != nil {
		return "", "", err
	}

	for {
		i := rand.Intn(len(accs))
		r := accs[i]
		if r.Username != "" && r.Password != "" {
			return r.Username, r.Password, nil
		}
		accs = append(accs[:i], accs[i+1:]...)
		if len(accs) == 0 {
			return "", "", ErrNoValidLogin
		}
	}
}

// ProvisionEnv will check the environment variables for INSTAGRAM_ACT_ and
//   create a base64 encoded config for the account, and write it to path.
//
// :param: path               - path a file to use as env, commonly a .env, but not required.
// :param: refresh (OPTIONAL) - refresh all plaintext credentials, don't skip already converted accounts
//
// This function has been created the use of a .env file in mind.
//
// .env contents:
//   INSTAGRAM_ACT_<name>="user:pass"
//
// This function will add to the .env:
//   INSTAGRAM_BASE64_<name>="..."
//
func ProvisionEnv(path string, refresh ...bool) error {
	skipList := []encAcc{}

	// By default, skip exisitng accounts
	if len(refresh) == 0 || (len(refresh) > 0 && !refresh[0]) {
		var err error
		skipList, err = avEncAcc()
		if err != nil {
			return err
		}
	}

	// Open File
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read file
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return err
	}

	newBuf := new(bytes.Buffer)

	// Read Lines
	for buf.Len() != 0 {
		l, err := buf.ReadBytes(byte('\n'))
		if err != nil {
			return err
		}

		// Process Lines
		line := string(l)
		if strings.HasPrefix(line, "INSTAGRAM_ACT_") {
			// Extract Creds
			split := strings.Split(line, "=")
			name := split[0][14:]
			c := split[1][1:]
			creds := strings.Split(strings.Split(c, "\"")[0], ":")

			// Create new config if not already exists
			var encLine string
			for _, existing := range skipList {
				found := false
				if existing.Username == strings.ToLower(creds[0]) {
					found = true
				}

				if found {
					// Use old config
					encLine = existing.Line
				} else {
					// Login
					fmt.Println("Processing", creds[0])
					insta := New(creds[0], creds[1])
					err := insta.Login()
					if err != nil {
						return err
					}
					// Export Config
					enc, err := insta.ExportAsBase64String()
					if err != nil {
						return err
					}
					encLine = fmt.Sprintf("INSTAGRAM_BASE64_%s=\"%s\"\n\n", name, enc)
				}
			}

			// Write Config
			_, err = newBuf.WriteString(line)
			if err != nil {
				return err
			}
			_, err = newBuf.WriteString(encLine)
			if err != nil {
				return err
			}
		} else if !strings.HasPrefix(line, "INSTAGRAM") {
			_, err = newBuf.WriteString(line)
			if err != nil {
				return err
			}
		}
	}

	// Write File
	err = ioutil.WriteFile(path, newBuf.Bytes(), 0o644)
	if err != nil {
		return err
	}
	return nil
}

// checkEnv will check the env variables for accounts that do have a login,
//    but no config (INSTAGRAM_BASE64_<...>). If one is found, call ProvisionEnv
func checkEnv(path ...string) error {
	avEnc, err := avEncAcc(path...)
	if err != nil {
		return err
	}
	avPlain, err := avPlainAcc(path...)
	if err != nil {
		return err
	}

	// Check for every plain acc, if there is an encoded equivalent
	for _, plain := range avPlain {
		found := false
		for _, enc := range avEnc {
			insta, err := ImportFromBase64String(enc.Base64, true)
			if err != nil {
				return err
			}
			if insta.Account.Username == plain.Username {
				found = true
			}
		}
		// If even one account has no config, provision
		if !found {
			fmt.Println("Unable to find:", plain.Name, plain.Username)
			p := ".env"
			if len(path) > 0 {
				p = path[0]
			}
			ProvisionEnv(p, true)
			break
		}
	}
	return nil
}

// getRandAcc returns a random insta instance from env
func getRandAcc(path ...string) (*Instagram, error) {
	accounts, err := avEncAcc(path...)
	if err != nil {
		return nil, err
	}

	// validate result
	if len(accounts) == 0 {
		return nil, ErrNoValidLogin
	}

	// select rand account
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(len(accounts))

	// load account config
	insta, err := ImportFromBase64String(accounts[r].Base64, true)
	if err != nil {
		return nil, err
	}
	return insta, err
}

// avEncAcc gets all available encoded accounts
func avEncAcc(path ...string) ([]encAcc, error) {
	output := []encAcc{}

	environ, err := loadEnv(path...)
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

			insta, err := ImportFromBase64String(encodedString, true)
			if err != nil {
				return nil, err
			}

			x := encAcc{
				Username: insta.Account.Username,
				Line:     env,
				Base64:   encodedString,
			}
			output = append(output, x)
		}
	}

	return output, nil
}

// avPlainAcc will check the environment for INSTAGRAM_ACT_, and return a list
//   of all found results.
func avPlainAcc(path ...string) ([]plainAcc, error) {
	var accs []plainAcc

	environ, err := loadEnv(path...)
	if err != nil {
		return nil, err
	}

	for _, env := range environ {
		if strings.HasPrefix(env, "INSTAGRAM_ACT_") {
			envVar := strings.Split(env, "=")
			name := strings.Split(envVar[0], "_")[2]
			acc := envVar[1]
			if acc[0] == '"' {
				acc = acc[1 : len(acc)-1]
			}
			creds := strings.Split(acc, ":")

			x := plainAcc{
				Name:     name,
				Username: strings.ToLower(creds[0]),
				Password: creds[1],
			}
			accs = append(accs, x)
		}
	}

	if len(accs) == 0 {
		return nil, ErrNoValidLogin
	}
	return accs, nil
}

// loadEnv loads all the environment variables
func loadEnv(p ...string) ([]string, error) {
	path := ".env"
	if len(p) > 0 {
		path = p[0]
	}
	environ := os.Environ()
	env, err := dotenv(path)
	if err != nil {
		return nil, err
	}
	environ = append(environ, env...)
	return environ, nil
}

// dotenv loads the file in the path parameter (commonly a .env file) and
//   returns its contents.
func dotenv(path string) ([]string, error) {
	file, err := os.Open(path)
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
