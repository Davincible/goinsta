package tests

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/Davincible/goinsta"
)

func TestTimeline(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	tl := insta.Timeline
	next := tl.NextID
outside:
	for i := 0; tl.Next(); i++ {
		if tl.NextID == next {
			t.Fatal("Next ID must be different after each request")
		}
		next = tl.NextID
		if i == 5 {
			break outside
		}
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}

	t.Logf("Gathered %d posts, %f on last request\n", len(tl.Items), tl.NumResults)
}

func TestDownload(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	if !insta.Timeline.Next() {
		t.Fatal(insta.Timeline.Error())
	}
	posts := insta.Timeline.Items
	if len(posts) == 0 {
		t.Fatal("No posts found")
	}

	rand.Seed(time.Now().UnixNano())
	randN := rand.Intn(len(posts))
	post := posts[randN]

	folder := "downloads/" + strconv.FormatInt(time.Now().Unix(), 10)
	err = post.Download(folder, "")
	if err != nil {
		t.Fatal(err)
	}

	randN = rand.Intn(len(posts))
	post = posts[randN]
	err = post.Download(folder, "testy")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Downloaded posts to %s", folder)
}
