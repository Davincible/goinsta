package tests

import (
	"math/rand"
	"testing"
	"time"
)

func TestTimeline(t *testing.T) {
	insta, err := getRandomAccount()
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
