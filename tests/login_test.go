package tests

import (
	"testing"

	"github.com/Davincible/goinsta"
)

func TestImportAccount(t *testing.T) {
	// Test Login
	user, pass, err := getLogin()
	if err != nil {
		t.Fatal(err)
		return
	}

	insta := goinsta.New(user, pass)
	err = insta.Login()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Logged in successfully as %s\n", user)
	logPosts(t, insta)

	// Test Import
	insta, err = getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}
	insta.OpenApp()
	t.Logf("logged into Instagram as user '%s'", insta.Account.Username)
	logPosts(t, insta)
}

func logPosts(t *testing.T, insta *goinsta.Instagram) {
	t.Logf("Gathered %d Timeline posts, %d Stories, %d Discover items, %d Notifications",
		len(insta.Timeline.Items),
		len(insta.Timeline.Tray.Stories),
		len(insta.Discover.AllItems),
		len(insta.Activity.NewStories)+len(insta.Activity.OldStories),
	)
}
