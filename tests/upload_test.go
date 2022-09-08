package tests

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"testing"

	"github.com/Davincible/goinsta/v3"
)

func TestUploadPhoto(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	insta.SetWarnHandler(t.Log)

	// Get random photo
	photo, err := getPhoto(1400, 1400)
	if err != nil {
		log.Fatal(err)
	}

	results, err := insta.Searchbar.SearchLocation("New York")
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Places) == 0 {
		t.Fatal(errors.New("No search result found"))
	}
	location := results.Places[0].Location
	if err := results.RegisterLocationClick(location); err != nil {
		t.Fatal(err)
	}

	item, err := insta.Upload(
		&goinsta.UploadOptions{
			File:     photo,
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

func TestUploadThumbVideo(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	insta.SetWarnHandler(t.Log)

	// Get random video
	video, err := getVideo()
	if err != nil {
		t.Fatal(err)
	}
	size := float64(len(video.Content)) / 1000000.0
	t.Logf("Video size: %.2f Mb", size)

	photo, err := getPhoto(1920, 1080)
	if err != nil {
		log.Fatal(err)
	}

	// Find location
	results, err := insta.Searchbar.SearchLocation("Chicago")
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Places) == 0 {
		t.Fatal(errors.New("No search result found"))
	}

	location := results.Places[1].Location
	if err := results.RegisterLocationClick(location); err != nil {
		t.Fatal(err)
	}

	item, err := insta.Upload(
		&goinsta.UploadOptions{
			File:      bytes.NewReader(video.Content),
			Thumbnail: photo,
			Caption:   "What a terrific video! #art",
			Location:  location.NewPostTag(),
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
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	insta.SetWarnHandler(t.Log)

	// Get random video
	video, err := getVideo()
	if err != nil {
		t.Fatal(err)
	}
	size := float64(len(video.Content)) / 1000000.0
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
	if err := results.RegisterLocationClick(location); err != nil {
		t.Fatal(err)
	}

	item, err := insta.Upload(
		&goinsta.UploadOptions{
			File:     bytes.NewReader(video.Content),
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

func TestUploadStoryPhoto(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	insta.SetWarnHandler(t.Log)

	// Get random photo
	photo, err := getPhoto(1400, 1400)
	if err != nil {
		log.Fatal(err)
	}

	item, err := insta.Upload(
		&goinsta.UploadOptions{
			File:    photo,
			IsStory: true,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("The ID of the new upload is %s", item.ID)
}

func TestUploadStoryVideo(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	insta.SetWarnHandler(t.Log)

	// Get random video
	video, err := getVideo(map[string]interface{}{"max_length": 20})
	if err != nil {
		t.Fatal(err)
	}
	size := float64(len(video.Content)) / 1000000.0
	t.Logf("Video size: %.2f Mb", size)

	item, err := insta.Upload(
		&goinsta.UploadOptions{
			File:    bytes.NewReader(video.Content),
			IsStory: true,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("The ID of the new upload is %s", item.ID)
}

func TestUploadStoryMultiVideo(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	insta.SetWarnHandler(t.Log)

	// Get random videos
	album := []io.Reader{}
	for i := 0; i < 5; i++ {
		video, err := getVideo(map[string]interface{}{"max_length": 20})
		if err != nil {
			t.Fatal(err)
		}
		size := float64(len(video.Content)) / 1000000.0
		t.Logf("Video size: %.2f Mb", size)
		album = append(album, bytes.NewReader(video.Content))
	}

	item, err := insta.Upload(
		&goinsta.UploadOptions{
			Album:   album,
			IsStory: true,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("The ID of the last story is %s", item.ID)
}

func TestUploadCarousel(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	insta.SetWarnHandler(t.Log)

	// Get random photos
	album := []io.Reader{}
	for i := 0; i < 5; i++ {
		photo, err := getPhoto(1920, 1080)
		if err != nil {
			log.Fatal(err)
		}

		album = append(album, photo)
	}

	// Add video to album
	video, err := getVideo()
	if err != nil {
		t.Fatal(err)
	}
	size := float64(len(video.Content)) / 1000000.0
	t.Logf("Video size: %.2f Mb", size)
	album = append(album, bytes.NewReader(video.Content))

	results, err := insta.Searchbar.SearchLocation("New York")
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Places) == 0 {
		t.Fatal(errors.New("No search result found"))
	}
	location := results.Places[1].Location
	if err := results.RegisterLocationClick(location); err != nil {
		t.Fatal(err)
	}

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

func TestUploadProfilePicture(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	insta.SetWarnHandler(t.Log)

	file := "./downloads/1645304867/testy_1.jpg"
	b, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewBuffer(b)
	if err := insta.Account.ChangeProfilePic(buf); err != nil {
		t.Fatal(err)
	}
	t.Log("Changed profile picture!")
}
