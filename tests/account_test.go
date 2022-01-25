package tests

import (
	"testing"

	"github.com/Davincible/goinsta"
)

func TestPendingFriendships(t *testing.T) {
	insta, err := goinsta.EnvRandAcc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logged in as %s\n", insta.Account.Username)

	count, err := insta.Account.PendingRequestCount()
	if err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("No pending friend requests found")
	}
	t.Logf("Found %d pending frienships\n", count)

	result, err := insta.Account.PendingFollowRequests()
	if err != nil {
		t.Fatal(err)
	}
	pending := result.Users

	err = pending[0].ApprovePending()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Approved request for %s\n", pending[0].Username)

	if len(pending) >= 2 {
		err = pending[1].IgnorePending()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Ignored request for %s\n", pending[1].Username)
	}
	count, err = insta.Account.PendingRequestCount()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("After approving there are %d pending friendships remaining\n", count)
}
