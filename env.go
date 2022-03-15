package goinsta

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

// EnvPlainAcc represents the plain account details stored in the env variable:
//
//   INSTAGRAM_ACT_<name>="username:password"
type EnvPlainAcc struct {
	Name     string
	Username string
	Password string
}

// EnvEncAcc represents the encoded account details stored in the env variable:
//
//   INSTAGRAM_BASE64_<name>="<base64 encoded config>"
type EnvEncAcc struct {
	Name     string
	Username string
	Base64   string
}

// EnvAcc represents the pair of plain and base64 encoded account pairs as
//   stored in EnvPlainAcc and EnvEncAcc, with env variables:
//
//   INSTAGRAM_ACT_<name>="username:password"
//   INSTAGRAM_BASE64_<name>="<base64 encoded config>"
type EnvAcc struct {
	Plain *EnvPlainAcc
	Enc   *EnvEncAcc
}

var (
	errNoAcc = errors.New("No Account Found")
)

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
	allAccs, err := EnvReadAccs(path...)
	if err != nil {
		return "", "", err
	}

	// extract plain accounts from all accounts
	accs := []*EnvPlainAcc{}
	for _, acc := range allAccs {
		if acc.Plain != nil {
			accs = append(accs, acc.Plain)
		}
	}

	// find valid one, until list exhausted
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

// EnvProvision will check the environment variables for INSTAGRAM_ACT_ and
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
func EnvProvision(path string, refresh ...bool) error {
	// By default, skip exisitng accounts
	refreshFlag := false
	if len(refresh) == 0 || (len(refresh) > 0 && !refresh[0]) {
		refreshFlag = true
	}
	fmt.Printf("Force refresh is set to %v\n", refreshFlag)

	accs, other, err := envLoadAccs(path)
	if err != nil {
		return err
	}

	fmt.Printf("Found %d accounts\n", len(accs))

	for _, acc := range accs {
		if acc.Enc != nil && !refreshFlag {
			fmt.Printf("Skipping account %s\n", acc.Plain.Name)
			continue
		}
		username := acc.Plain.Username
		password := acc.Plain.Password
		fmt.Println("Processing", username)
		insta := New(username, password)
		err := insta.Login()
		if err != nil {
			return err
		}
		// Export Config
		enc, err := insta.ExportAsBase64String()
		if err != nil {
			return err
		}
		acc.Enc.Base64 = enc
		fmt.Println("Sleeping...")
		time.Sleep(20 * time.Second)
	}
	err = accsToFile(path, accs, other)
	if err != nil {
		return err
	}

	return nil
}

// EnvUpdateAccs will update the plain and encoded account variables stored in
//  the .env file:
//
//   INSTAGRAM_ACT_<name>="username:password"
//   INSTAGRAM_BASE64_<name>="<base64 encoded config>"
//
// :param: string:path -- file path of the .env file, typically ".env"
// :param: []*EncAcc:newAccs -- list of updated versions of the accounts
func EnvUpdateAccs(path string, newAccs []*EnvAcc) error {
	return envUpdateAccs(path, newAccs)
}

// EnvUpdateEnc will update the encoded account variables stored in
//  the .env file:
//
//   INSTAGRAM_BASE64_<name>="<base64 encoded config>"
//
// :param: string:path -- file path of the .env file, typically ".env"
// :param: []*EnvEncAcc:newAccs -- list of updated encoded accounts
func EnvUpdateEnc(path string, newAccs []*EnvEncAcc) error {
	return envUpdateAccs(path, newAccs)
}

// EnvPlainAccs will update the plain account variables stored in
//  the .env file:
//
//   INSTAGRAM_ACT_<name>="username:password"
//
// :param: string:path -- file path of the .env file, typically ".env"
// :param: []*EnvPlainAcc:newAccs -- list of updated plain accounts
func EnvUpdatePlain(path string, newAccs []*EnvPlainAcc) error {
	return envUpdateAccs(path, newAccs)
}

func envUpdateAccs(path string, newAccs interface{}) error {
	accs, other, err := dotenv(path)
	if err != nil {
		return err
	}
	switch n := newAccs.(type) {
	case []*EnvAcc:

	case []*EnvPlainAcc:
		for _, acc := range n {
			accs = addOrUpdateAcc(accs, acc)
		}
	case []*EnvEncAcc:
		for _, acc := range n {
			accs = addOrUpdateAcc(accs, acc)
		}
	}

	return accsToFile(path, accs, other)
}

// checkEnv will check the env variables for accounts that do have a login,
//    but no config (INSTAGRAM_BASE64_<...>). If one is found, call ProvisionEnv
func checkEnv(path ...string) error {
	accs, err := EnvReadAccs(path...)
	if err != nil {
		return err
	}

	// Check for every plain acc, if there is an encoded equivalent
	for _, acc := range accs {
		if acc.Plain != nil && acc.Enc == nil {
			fmt.Println("Unable to find:", acc.Plain.Name, acc.Plain.Username)
			p := ".env"
			if len(path) > 0 {
				p = path[0]
			}
			EnvProvision(p, true)
		}
	}
	return nil
}

// getRandAcc returns a random insta instance from env
func getRandAcc(path ...string) (*Instagram, error) {
	allAccs, err := EnvReadAccs(path...)
	if err != nil {
		return nil, err
	}

	// extract encoded accounts from all accounts
	accounts := []*EnvEncAcc{}
	for _, acc := range allAccs {
		if acc.Enc != nil {
			accounts = append(accounts, acc.Enc)
		}
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
	// insta.SetProxy("http://localhost:9090", false, true)
	return insta, err
}

// EnvLoadPlain will load all plain accounts stored in the env variables:
//
//   INSTAGRAM_ACT_<name>="username:password"
//
// :param: path (OPTIONAL) -- .env file to load, default to ".env"
func EnvLoadPlain(path ...string) ([]*EnvPlainAcc, error) {
	allAccs, err := EnvReadAccs(path...)
	if err != nil {
		return nil, err
	}

	// extract plain accounts from all accounts
	accs := []*EnvPlainAcc{}
	for _, acc := range allAccs {
		if acc.Plain != nil {
			accs = append(accs, acc.Plain)
		}
	}
	return accs, nil
}

// EnvLoadAccs loads all the environment variables.
//
// By default, the OS environment variables as well as .env are loaded
// To load a custom file, instead of .env, pass the filepath as an argument.
//
// Don't Sync param is set to true to prevent any http calls on import by default
func EnvLoadAccs(p ...string) ([]*Instagram, error) {
	instas := []*Instagram{}
	accs, _, err := envLoadAccs(p...)

	for _, acc := range accs {
		insta, err := ImportFromBase64String(acc.Enc.Base64, true)
		if err != nil {
			return nil, err
		}
		instas = append(instas, insta)
	}

	return instas, err
}

// EnvReadAccs loads both all plain and base64 encoded accounts
//
// Set in a .env file or export to your environment variables:
//   INSTAGRAM_ACT_<name>="user:pass"
//   INSTAGRAM_BASE64_<name>="..."
//
// :param: p (OPTIONAL) -- env file path, ".env" by default
func EnvReadAccs(p ...string) ([]*EnvAcc, error) {
	accs, _, err := envLoadAccs(p...)
	return accs, err
}

func envLoadAccs(p ...string) ([]*EnvAcc, []string, error) {
	path := ".env"
	if len(p) > 0 {
		path = p[0]
	}
	environ := os.Environ()
	accsOne, other, err := dotenv(path)
	if err != nil {
		return nil, nil, err
	}

	accsTwo, _, err := parseAccs(environ)
	if err != nil {
		return nil, nil, err
	}
	accs := append(accsOne, accsTwo...)
	if len(accs) == 0 {
		return nil, nil, ErrNoValidLogin
	}

	return accs, other, nil
}

// dotenv loads the file in the path parameter (commonly a .env file) and
//   returns its contents.
func dotenv(path string) ([]*EnvAcc, []string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, nil, err
	}
	lines := strings.Split(string(buf.Bytes()), "\n")
	accs, other, err := parseAccs(lines)
	if err != nil {
		return nil, nil, err
	}
	return accs, other, nil
}

func parseAccs(lines []string) ([]*EnvAcc, []string, error) {
	accs := []*EnvAcc{}
	other := []string{}
	for _, line := range lines {
		plain, errPlain := parsePlain(line)
		if errPlain == nil {
			accs = addOrUpdateAcc(accs, plain)
		}

		enc, errEnc := parseEnc(line)
		if errEnc == nil {
			accs = addOrUpdateAcc(accs, enc)
		} else if errEnc != errNoAcc {
			return nil, nil, errEnc
		}

		// Keep track of other lines in the file, not related to goinsta
		if !strings.HasPrefix(line, "INSTAGRAM") {
			other = append(other, line)
		}
	}
	return accs, other, nil
}

func parsePlain(line string) (*EnvPlainAcc, error) {
	if strings.HasPrefix(line, "INSTAGRAM_ACT_") {
		envVar := strings.Split(line, "=")
		name := strings.Split(envVar[0], "_")[2]
		acc := envVar[1]
		if acc[0] == '"' {
			acc = acc[1 : len(acc)-1]
		}
		creds := strings.Split(acc, ":")

		return &EnvPlainAcc{
			Name:     name,
			Username: strings.ToLower(creds[0]),
			Password: creds[1],
		}, nil
	}
	return nil, errNoAcc
}

func parseEnc(line string) (*EnvEncAcc, error) {
	if strings.HasPrefix(line, "INSTAGRAM_BASE64_") {
		envVar := strings.SplitN(line, "=", 2)
		encodedString := strings.TrimSpace(envVar[1])
		name := strings.Split(envVar[0], "_")[2]
		if encodedString[0] == '"' {
			encodedString = encodedString[1 : len(encodedString)-1]
		}
		insta, err := ImportFromBase64String(encodedString, true)
		if err != nil {
			return nil, err
		}

		return &EnvEncAcc{
			Name:     name,
			Username: insta.Account.Username,
			Base64:   encodedString,
		}, nil
	}
	return nil, errNoAcc
}

func addOrUpdateAcc(accs []*EnvAcc, toAdd interface{}) []*EnvAcc {
	switch newAcc := toAdd.(type) {
	case *EnvAcc:
		for _, acc := range accs {
			if (acc.Enc != nil && newAcc.Plain.Username == acc.Enc.Username) ||
				(acc.Plain != nil && newAcc.Enc.Username == acc.Plain.Username) {
				if newAcc.Plain != nil {
					if newAcc.Plain.Name == "" {
						newAcc.Plain.Name = acc.Plain.Name
					}

					acc.Plain = newAcc.Plain
				}
				if newAcc.Enc != nil {
					if newAcc.Enc.Name == "" {
						newAcc.Enc.Name = acc.Enc.Name
					}
					acc.Enc = newAcc.Enc
				}
				return accs
			}
		}
	case *EnvPlainAcc:
		for _, acc := range accs {
			if (acc.Enc != nil && newAcc.Username == acc.Enc.Username) ||
				(acc.Plain != nil && newAcc.Username == acc.Plain.Username) {
				if newAcc.Name == "" {
					newAcc.Name = acc.Plain.Name
				}
				acc.Plain = newAcc
				return accs
			}
		}
		accs = append(accs, &EnvAcc{Plain: newAcc})
	case *EnvEncAcc:
		for _, acc := range accs {
			if (acc.Plain != nil && newAcc.Username == acc.Plain.Username) ||
				(acc.Enc != nil && newAcc.Username == acc.Enc.Username) {
				if newAcc.Name == "" {
					newAcc.Name = acc.Enc.Name
				}
				acc.Enc = newAcc
				return accs
			}
		}
		accs = append(accs, &EnvAcc{Enc: newAcc})
	}
	return accs
}

func accsToFile(path string, accs []*EnvAcc, other []string) error {
	newBuf := new(bytes.Buffer)
	for _, acc := range accs {
		if acc.Plain != nil {
			line := fmt.Sprintf("INSTAGRAM_ACT_%s=\"%s:%s\"\n", acc.Plain.Name, acc.Plain.Username, acc.Plain.Password)
			_, err := newBuf.WriteString(line)
			if err != nil {
				return err
			}

		}
		if acc.Enc != nil {
			encLine := fmt.Sprintf("INSTAGRAM_BASE64_%s=\"%s\"\n\n", acc.Enc.Name, acc.Enc.Base64)
			_, err := newBuf.WriteString(encLine)
			if err != nil {
				return err
			}
		}
	}

	err := ioutil.WriteFile(path, newBuf.Bytes(), 0o644)
	if err != nil {
		return err
	}
	return nil
}

func (acc *Account) GetEnvEncAcc() (*EnvEncAcc, error) {
	b, err := acc.insta.ExportAsBase64String()
	return &EnvEncAcc{
		Username: acc.Username,
		Base64:   b,
	}, err
}
