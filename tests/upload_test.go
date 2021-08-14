package tests

import (
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

	item, err := insta.Upload(
		&goinsta.UploadOptions{
			File:    resp.Body,
			Caption: "awesome! :) #41",
		},
	)
	if err != nil {
		Error(t, err)
	}
	t.Logf("The ID of the new upload is %s", item.ID)
}
