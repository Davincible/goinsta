package tests

import (
	"math/rand"
	"testing"
)

func TestSearchUser(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	// Search for users
	result, err := insta.Searchbar.SearchUser("a")
	if err != nil {
		t.Fatal(err)
		return
	}
	if result.Status != "ok" {
		t.Fatal(result.Status)
		return
	}
	t.Logf("Result length is %d", len(result.Users))

	// Select a random user
	user := result.Users[rand.Intn(len(result.Users))]
	err = result.RegisterUserClick(user)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Get user info
	err = user.Info()
	if err != nil {
		t.Fatal(err)
		return
	}

	// Get user feed
	feed := user.Feed()
	s := feed.Next()
	if !s {
		t.Fatalf("Failed to get feed: %s", feed.Error())
		return
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
			return
		}
	}
}

func TestSearchHashtag(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	// Search for hashtags
	result, err := insta.Searchbar.SearchHashtag("photography")
	if err != nil {
		t.Fatal(err)
		return
	}
	if result.Status != "ok" {
		t.Fatal(result.Status)
		return
	}
	t.Logf("Result length is %d", len(result.Users))

	// Select a random user
	// tag := result.Ta
	user := result.Users[rand.Intn(len(result.Users))]
	err = result.RegisterUserClick(user)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Get user info
	err = user.Info()
	if err != nil {
		t.Fatal(err)
		return
	}

	// Get user feed
	feed := user.Feed()
	s := feed.Next()
	if !s {
		t.Fatalf("Failed to get feed: %s", feed.Error())
		return
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
			return
		}
	}
}

func TestSearchLocation(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	// Search for hashtags
	result, err := insta.Searchbar.SearchLocation("New York")
	if err != nil {
		t.Fatal(err)
		return
	}
	if result.Status != "ok" {
		t.Fatal(result.Status)
		return
	}
	t.Logf("Result length is %d", len(result.Places))
}
