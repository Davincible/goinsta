package tests

import (
	"testing"

	"github.com/Davincible/goinsta"
)

func TestIGTVChannel(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	user, err := insta.Profiles.ByName("f1")
	if err != nil {
		t.Fatal(err)
	}

	channel, err := user.IGTV()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Fetched %d IGTV items", len(channel.Items))

	if !channel.Next() {
		err := channel.Error()
		if err == goinsta.ErrNoMore {
			t.Logf("No more posts available (%v)", channel.MoreAvailable)
			return
		}
		t.Fatal(channel.Error())
	}
	t.Logf("Fetched %d IGTV items", len(channel.Items))
}

func TestIGTVSeries(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	user, err := insta.Profiles.ByName("kimkotter")
	if err != nil {
		t.Fatal(err)
	}

	series, err := user.IGTVSeries()
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, serie := range series {
		total += len(serie.Items)
	}
	t.Logf("Found %d series. Containing %d IGTV posts in total.", len(series), total)
}

func TestIGTVLive(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	broadcasts, err := insta.IGTV.Live()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Found %d broadcasts. More Available = %v", len(broadcasts.Broadcasts), broadcasts.MoreAvailable)

	if broadcasts.MoreAvailable {
		broadcasts, err := broadcasts.Live()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Found %d broadcasts. More Available = %v", len(broadcasts.Broadcasts), broadcasts.MoreAvailable)

	}
}

func TestIGTVDiscover(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	for i := 0; i < 5; i++ {
		if !insta.IGTV.Next() {
			t.Fatal(insta.IGTV.Error())
		}
		t.Logf("Fetched %d posts", len(insta.IGTV.Items))
	}
}
