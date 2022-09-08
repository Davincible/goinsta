package tests

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/Davincible/goinsta/v3"
)

func TestSearchUser(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	// Search for users
	result, err := insta.Searchbar.SearchUser("nicky")
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "ok" {
		t.Fatal(result.Status)
	}
	t.Logf("Result length is %d", len(result.Users))

	// Select a random user
	if len(result.Users) == 0 {
		t.Fatal("No search results found! Change search query or fix api")
	}
	user := result.Users[rand.Intn(len(result.Users))]
	err = result.RegisterUserClick(user)
	if err != nil {
		t.Fatal(err)
	}

	// Get user info
	err = user.Info()
	if err != nil {
		t.Fatal(err)
	}

	// Get user feed
	feed := user.Feed()
	if !feed.Next() {
		t.Fatalf("Failed to get feed: %s", feed.Error())
	}
	t.Logf("Found %d posts", len(feed.Items))

	// Err if no posts are found
	if len(feed.Items) == 0 && user.MediaCount != 0 {
		t.Fatal("Failed to fetch any posts while the user does have posts")
	} else if len(feed.Items) != 0 {
		// Like a random post to make sure the insta pointer is set and working
		post := feed.Items[rand.Intn(len(feed.Items))]
		err := post.Like()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestSearchHashtag(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	// Search for hashtags
	query := "photography"
	result, err := insta.Searchbar.SearchHashtag(query)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "ok" {
		t.Fatal(result.Status)
	}
	if len(result.Tags) == 0 {
		t.Fatal("No results found")
	}
	t.Logf("Result length is %d", len(result.Tags))

	var hashtag *goinsta.Hashtag
	for _, tag := range result.Tags {
		if tag.Name == query {
			result.RegisterHashtagClick(tag)
			hashtag = tag
			break
		}
	}

	for i := 0; i < 5; i++ {
		if !hashtag.Next() {
			t.Fatal(hashtag.Error())
		}
		t.Logf("Fetched %d posts", len(hashtag.Items))
	}

	for i := 0; i < 5; i++ {
		if !hashtag.NextRecent() {
			t.Log(hashtag.Error())
		}
		t.Logf("Fetched %d recent posts", len(hashtag.ItemsRecent))
	}
}

func TestSearchLocation(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	// Search for hashtags
	result, err := insta.Searchbar.SearchLocation("New York")
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "ok" {
		t.Fatal(errors.New(result.Status))
	}
	t.Logf("Result length is %d", len(result.Places))

	location := result.Places[0].Location
	feed, err := location.Feed()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Found %d sections", len(feed.Sections))
}
