package tests

import (
	"testing"
)

func TestSearchUser(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := insta.Searchbar.SearchUser("a")
	if err != nil {
		t.Fatal(err)
		return
	}
	if result.Status != "ok" {
		t.Fatal(result.Status)
		return
	}
	t.Logf("result length is %d", len(result.Users))
	for _, user := range result.Users {
		t.Logf("user %s with id %d\n", user.Username, user.ID)
	}
}
