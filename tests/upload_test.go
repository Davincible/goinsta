package tests

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/Davincible/goinsta"
)

func TestUploadPhoto(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	insta.SetErrorHandler(t.Log)

	// Get random photo
	resp, err := http.Get("https://picsum.photos/800/800")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	results, err := insta.Searchbar.SearchLocation("New York")
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Places) == 0 {
		t.Fatal(errors.New("No search result found"))
	}
	location := results.Places[0].Location
	results.RegisterLocationClick(location)

	item, err := insta.Upload(
		&goinsta.UploadOptions{
			File:     resp.Body,
			Caption:  "awesome! :) #41",
			Location: location.NewPostTag(),
			UserTags: &[]goinsta.UserTag{
				{
					User: &goinsta.User{
						ID: insta.Account.ID,
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("The ID of the new upload is %s", item.ID)
}

func TestUploadVideo(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	insta.SetErrorHandler(t.Log)

	// Get random video
	video, err := getVideo()
	if err != nil {
		t.Fatal(err)
	}
	size := float64(len(video)) / 1000000.0
	t.Logf("Video size: %.2f Mb", size)

	// Find location
	results, err := insta.Searchbar.SearchLocation("Bali")
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Places) == 0 {
		t.Fatal(errors.New("No search result found"))
	}
	location := results.Places[1].Location
	results.RegisterLocationClick(location)

	item, err := insta.Upload(
		&goinsta.UploadOptions{
			File:     bytes.NewReader(video),
			Caption:  "What a terrific video! #art",
			Location: location.NewPostTag(),
			UserTags: &[]goinsta.UserTag{
				{
					User: &goinsta.User{
						ID: insta.Account.ID,
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("The ID of the new upload is %s", item.ID)
}

func TestUploadCarousel(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	insta.SetErrorHandler(t.Log)

	// Get random photos
	album := []io.Reader{}
	for i := 0; i < 5; i++ {
		resp, err := http.Get("https://picsum.photos/1400/1400")
		if err != nil {
			log.Fatal(err)
		}
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if err := resp.Body.Close(); err != nil {
			t.Fatal(err)
		}
		buf := bytes.NewReader(bodyBytes)
		album = append(album, buf)
	}

	// Add video to album
	video, err := getVideo()
	if err != nil {
		t.Fatal(err)
	}
	size := float64(len(video)) / 1000000.0
	t.Logf("Video size: %.2f Mb", size)
	album = append(album, bytes.NewReader(video))

	results, err := insta.Searchbar.SearchLocation("New York")
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Places) == 0 {
		t.Fatal(errors.New("No search result found"))
	}
	location := results.Places[1].Location
	results.RegisterLocationClick(location)

	// Upload Album
	item, err := insta.Upload(
		&goinsta.UploadOptions{
			Album:    album,
			Location: location.NewPostTag(),
			Caption:  "The best photos I've seen all morning!",
			UserTags: &[]goinsta.UserTag{
				{
					User: &goinsta.User{
						ID: insta.Account.ID,
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("The ID of the new upload is %s", item.ID)
}
