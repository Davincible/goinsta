package tests

import (
	"testing"
	"time"

	"github.com/Davincible/goinsta/v3"
)

func TestIGTVChannel(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
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
	insta, err := goinsta.EnvRandAcc()
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
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	igtv, err := insta.IGTV.Live()
	if err != nil {
		t.Fatal(err)
	}

	broadcasts := igtv.Broadcasts
	t.Logf("Found %d broadcasts. More Available = %v\n", len(broadcasts), igtv.MoreAvailable)
	if len(broadcasts) == 0 {
		t.Fatal("No broadcasts found")
	}

	if igtv.MoreAvailable {
		time.Sleep(3 * time.Second)
		igtv, err := igtv.Live()
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("Found %d broadcasts. More Available = %v\n", len(igtv.Broadcasts), igtv.MoreAvailable)
		}
	}

	for _, br := range broadcasts {
		t.Logf("Broadcast by %s has %d viewers\n", br.User.Username, int(br.ViewerCount))
	}

	br := broadcasts[0]
	if err := br.GetInfo(); err != nil {
		t.Fatal(err)
	}

	comments, err := br.GetComments()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Comment Count: %d\n", comments.CommentCount)

	likes, err := br.GetLikes()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Like Count: %d\n", likes.Likes)

	heartbeat, err := br.GetHeartbeat()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Viewer Count: %d\n", int(heartbeat.ViewerCount))
}

func TestIGTVDiscover(t *testing.T) {
	t.Skip("Skipping IGTV Discover, depricated")

	insta, err := goinsta.EnvRandAcc()
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
