package tests

import (
	"math/rand"
	"testing"
	"time"

	"github.com/Davincible/goinsta/v3"
)

// Random big accounts used for story reply and DM tests
var possibleUsers = []string{
	"kimkardashian",
	"kyliejenner",
	"kendaljenner",
	"addisonraee",
	"iamcardib",
	"snoopdogg",
	"dualipa",
	"stassiebaby",
	"kourtneykardash",
	"f1",
	"madscandids",
	"9gag",
}

func TestStoryReply(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	insta.SetWarnHandler(t.Log)

	for _, acc := range possibleUsers {
		user, err := insta.Profiles.ByName(acc)
		if err != nil {
			t.Fatal(err)
		}
		stories, err := user.Stories()
		if err != nil {
			t.Fatal(err)
		}
		for _, story := range stories.Reel.Items {
			if story.CanReply {
				err := story.Reply("Nice! :)")
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("Replied to a story of %s", acc)
				return
			}
		}
	}
	t.Fatal("Unable to find a story to reply to")
}

func TestInboxSync(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	insta.SetWarnHandler(t.Log)

	if !insta.Inbox.InitialSnapshot() && insta.Inbox.Error() != goinsta.ErrNoMore {
		t.Fatal(insta.Inbox.Error())
	}

	if err := insta.Inbox.Sync(); err != nil {
		t.Fatal(err)
	}
	t.Logf("Fetched %d conversations", len(insta.Inbox.Conversations))
}

func TestInboxNew(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	insta.SetWarnHandler(t.Log)

	rand.Seed(time.Now().UnixNano())
	randUser := possibleUsers[rand.Intn(len(possibleUsers))]
	user, err := insta.Profiles.ByName(randUser)
	if err != nil {
		t.Fatal(err)
	}

	conv, err := insta.Inbox.New(user, "Roses are red, violets are blue, cushions are soft, and so are you")
	if err != nil {
		t.Fatal(err)
	}

	if err := conv.Send("Feeling poetic today uknow"); err != nil {
		t.Fatal(err)
	}
	t.Logf("DM'ed %s", randUser)
}
