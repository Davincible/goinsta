package tests

import (
	"testing"

	"github.com/Davincible/goinsta/v3"
)

func TestProfileVisit(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	profile, err := insta.VisitProfile("miakhalifa")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Gatherd %d posts, %d story items, %d highlights",
		len(profile.Feed.Items),
		len(profile.Stories.Reel.Items),
		len(profile.Highlights),
	)
}

func TestProfilesByName(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	_, err = insta.Profiles.ByName("binance")
	if err != nil {
		t.Fatal(err)
	}
}

func TestProfilesByID(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	_, err = insta.Profiles.ByID("28527810")
	if err != nil {
		t.Fatal(err)
	}
}

func TestProfilesBlocked(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	blocked, err := insta.Profiles.Blocked()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Foud %d blocked users", len(blocked))
}
