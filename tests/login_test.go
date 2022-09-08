package tests

import (
	"testing"

	"github.com/Davincible/goinsta/v3"
)

func TestImportAccount(t *testing.T) {
	// Test Import
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	if err := insta.OpenApp(); err != nil {
		t.Fatal(err)
	}

	t.Logf("logged into Instagram as user '%s'", insta.Account.Username)
	logPosts(t, insta)
}

func TestLogin(t *testing.T) {
	// Test Login
	user, pass, err := goinsta.EnvRandLogin()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Attempting to login as %s\n", user)
	if user == "codebrewernl" {
		t.Skip()
	}

	insta := goinsta.New(user, pass)
	if err = insta.Login(); err != nil {
		t.Fatal(err)
	}
	t.Log("Logged in successfully")
	logPosts(t, insta)

	// Test Logout
	if err := insta.Logout(); err != nil {
		t.Fatal(err)
	}
	t.Log("Logged out successfully")
}

func logPosts(t *testing.T, insta *goinsta.Instagram) {
	t.Logf("Gathered %d Timeline posts, %d Stories, %d Discover items, %d Notifications",
		len(insta.Timeline.Items),
		len(insta.Timeline.Tray.Stories),
		len(insta.Discover.Items),
		len(insta.Activity.NewStories)+len(insta.Activity.OldStories),
	)
}
