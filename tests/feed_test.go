package tests

import (
	"math/rand"
	"testing"
	"time"
)

func TestFeedUser(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	sr, err := insta.Searchbar.SearchUser("miakhalifa")
	if err != nil {
		t.Fatal(err)
		return
	}
	user := sr.Users[0]
	feed := user.Feed()

	next := feed.NextID
outside:
	for i := 0; feed.Next(); i++ {
		if feed.NextID == next {
			t.Fatal("Next ID must be different after each request")
			return
		}
		if feed.Status != "ok" {
			t.Fatalf("Status not ok: %s\n", feed.Status)
			return
		}
		next = feed.NextID
		if i == 5 {
			break outside
		}
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}

	t.Logf("Gathered %d posts, %d on last request\n", len(feed.Items), feed.NumResults)
}

func TestFeedDiscover(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	feed := insta.Discover
	next := feed.NextID

outside:
	for i := 0; feed.Next(); i++ {
		if feed.NextID == next {
			t.Fatal("Next ID must be different after each request")
			return
		}
		if feed.Status != "ok" {
			t.Fatalf("Status not ok: %s\n", feed.Status)
			return
		}
		next = feed.NextID
		if i == 5 {
			break outside
		}
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}

	t.Logf("Gathered %d posts, %d on last request\n", len(feed.AllItems), feed.NumResults)
}

func TestFeedTagLike(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	hashtag := insta.NewHashtag("golang")
	err = hashtag.Info()
	if err != nil {
		t.Fatal(err)
		return
	}

	// First round
	if s := hashtag.Next(); !s {
		t.Fatal(hashtag.Error())
		return
	}

	if len(hashtag.Items) == 0 {
		t.Logf("%+v", hashtag.Sections)
		t.Fatalf("Items length is 0, section length is %d\n", len(hashtag.Sections))
		return
	}

	for i, item := range hashtag.Items {
		err = item.Like()
		if err != nil {
			t.Fatal(err)
			return
		}
		t.Logf("media %s liked by goinsta", item.ID)
		if i == 5 {
			return
		}
		time.Sleep(3 * time.Second)
	}
}

func TestFeedTagNextOld(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	feedTag, err := insta.Feed.Tags("golang")
	if err != nil {
		t.Fatal(err)
		return
	}

	initNextID := feedTag.NextID
	success := feedTag.Next()
	if !success {
		t.Fatal("Failed to fetch next page")
		return
	}
	gotStatus := feedTag.Status

	if gotStatus != "ok" {
		t.Errorf("Status = %s; want ok", gotStatus)
	}

	gotNextID := feedTag.NextID
	if gotNextID == initNextID {
		t.Errorf("NextID must differ after FeedTag.Next() call")
	}
	t.Logf("Fetched %d posts", len(feedTag.Items))
}

func TestFeedTagNext(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	hashtag := insta.NewHashtag("golang")
	err = hashtag.Info()
	if err != nil {
		t.Fatal(err)
		return
	}

	// First round
	if s := hashtag.Next(); !s {
		t.Fatal(hashtag.Error())
		return
	}

	initNextID := hashtag.NextID

	// Second round
	if s := hashtag.Next(); !s {
		t.Fatal(hashtag.Error())
		return
	}

	if hashtag.Status != "ok" {
		t.Errorf("Status = %s; want ok", hashtag.Status)
	}

	gotNextID := hashtag.NextID
	if gotNextID == initNextID {
		t.Errorf("NextID must differ after FeedTag.Next() call")
	}
	t.Logf("Fetched %d posts", len(hashtag.Items))
}

func TestFeedTagNextRecent(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)
	hashtag := insta.NewHashtag("golang")
	err = hashtag.Info()
	if err != nil {
		t.Fatal(err)
		return
	}

	// First round
	if s := hashtag.NextRecent(); !s {
		t.Fatal(hashtag.Error())
		return
	}

	initNextID := hashtag.NextID

	// Second round
	if s := hashtag.NextRecent(); !s {
		t.Fatal(hashtag.Error())
		return
	}

	if hashtag.Status != "ok" {
		t.Errorf("Status = %s; want ok", hashtag.Status)
	}

	gotNextID := hashtag.NextID
	if gotNextID == initNextID {
		t.Errorf("NextID must differ after FeedTag.Next() call")
	}
	t.Logf("Fetched %d posts", len(hashtag.Recent))
}
