package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Davincible/goinsta/v3"
)

func TestStories(t *testing.T) {
	loginUser, pass, err := goinsta.EnvRandLogin()
	if err != nil {
		t.Fatal(err)
	}

	insta := goinsta.New(loginUser, pass)

	err = insta.Login()
	require.NoError(t, err, "failed to create new instance")

	userName := os.Getenv("IG_USERNAME")
	user, err := insta.UserInfo(userName)
	require.NoError(t, err, "failed to get user info")

	err = user.Info()
	require.NoError(t, err, "failed to get user info")

	stories, err := user.Stories()
	require.NoError(t, err, "failed to get user stories")

	t.Logf("Stories: %d", len(stories.Reel.MediaIDs))

	for _, story := range stories.Reel.Items {
		b, err := story.Download()
		require.NoError(t, err, "failed to download story")

		err = os.WriteFile(fmt.Sprintf("./downloads/story-%s-%s.jpg", userName, story.ID), b, 0644)
		require.NoError(t, err, "failed to write story")
	}
}
