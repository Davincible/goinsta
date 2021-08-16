package tests

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Davincible/goinsta"
)

var (
	errNoValidLogin = errors.New("No valid login found")
	errNoAPIKEY     = errors.New("No API Key has been found. Please add one to .env")
)

type pixaBayRes struct {
	Total     int `json:"total"`
	TotalHits int `json:"totalHits"`
	Hits      []struct {
		ID        int    `json:"id"`
		URL       string `json:"page_url"`
		Type      string `json:"type"`
		Tags      string `json:"tags"`
		Duration  int    `json:"duration"`
		PictureID string `json:"picture_id"`
		Videos    struct {
			Large  video `json:"large"`
			Medium video `json:"medium"`
			Small  video `json:"small"`
			Tiny   video `json:"tiny"`
		} `json:"videos"`
	} `json:"hits"`
}

type video struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Size   int    `json:"size"`
}

func readFromBase64(base64EncodedString string) (*goinsta.Instagram, error) {
	base64Bytes, err := base64.StdEncoding.DecodeString(base64EncodedString)
	if err != nil {
		return nil, err
	}
	return goinsta.ImportReader(bytes.NewReader(base64Bytes), true)
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
		insta, err := getRandomAccount()
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

func getRandomAccount() (*goinsta.Instagram, error) {
	accounts, err := availableEncodedAccounts()
	if err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, errNoValidLogin
	}

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(len(accounts))
	encodedAccount := accounts[r]
	insta, err := readFromBase64(encodedAccount)
	if err != nil {
		return nil, err
	}
	return insta, err
}

func getLogin() (string, string, error) {
	var accs [][]string

	environ, err := loadEnv()
	if err != nil {
		return "", "", err
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

func Error(t *testing.T, err error) {
	t.Log(err.Error())
	t.Fatal(err)
}

func getPixabayAPIKey() (string, error) {
	environ, err := loadEnv()
	if err != nil {
		return "", err
	}

	for _, env := range environ {
		if strings.HasPrefix(env, "PIXABAY_API_KEY") {
			key := strings.Split(env, "=")[1]
			if key[0] == '"' {
				key = key[1 : len(key)-1]
			}
			return key, nil
		}
	}
	return "", errNoAPIKEY
}

func getVideo(o ...map[string]interface{}) ([]byte, error) {
	key, err := getPixabayAPIKey()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://pixabay.com/api/videos/?key=%s&per_page=200", key)

	// Get video list
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	var res pixaBayRes
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	// Select random video
	max_length := 0
	if len(o) > 0 {
		options := o[0]
		max_length, _ = options["max_length"].(int)
	}

	valid := false
	var vid video
	for !valid {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(len(res.Hits))
		vid = res.Hits[r].Videos.Small
		if max_length == 0 || res.Hits[r].Duration < max_length {
			valid = true
		}
	}

	// Download video
	resp, err = http.Get(vid.URL)
	if err != nil {
		return nil, err
	}

	video, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}
	return video, nil
}
