package tests

import (
	"bytes"
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

	"github.com/Davincible/goinsta/v3"
)

var errNoAPIKEY = errors.New("No Pixabay API Key has been found. Please add one to .env as PIXABAY_API_KEY")

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
	URL     string `json:"url"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Size    int    `json:"size"`
	Content []byte
}

func TestEnvLoadAccs(t *testing.T) {
	accs, err := goinsta.EnvLoadAccs()
	if err != nil {
		t.Fatal(err)
	}

	if len(accs) == 0 {
		t.Fatalf("No accounts found")
	}
	t.Logf("Found %d accounts", len(accs))
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

func getPhoto(width, height int, i ...int) (io.Reader, error) {
	url := fmt.Sprintf("https://picsum.photos/%d/%d", width, height)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// Retry on failure
		if len(i) == 0 || i[0] < 5 {
			c := 0
			if len(i) > 0 {
				c = i[0] + 1
			}
			fmt.Println("Failed to get photo, retrying...")
			time.Sleep(5 * time.Second)
			return getPhoto(width, height, c)
		}
		return nil, fmt.Errorf("Get image status code %d", resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)

	return buf, err
}

func getVideo(o ...map[string]interface{}) (*video, error) {
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
	vid.Content = video
	return &vid, nil
}
