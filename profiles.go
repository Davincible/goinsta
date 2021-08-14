package goinsta

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Profiles allows user function interactions
type Profiles struct {
	insta *Instagram
}

// Profile represnts an instagram user with their various properties, such as
//   their account info, stored in Profile.User (a *User struct), feed, stories,
//   Highlights, IGTV posts, and friendship status.
//
type Profile struct {
	insta *Instagram

	User       *User
	Friendship *Friendship

	Feed       FeedMedia
	Stories    StoryMedia
	Highlights []Reel
	IGTV       IGTVChannel
}

// VisitProfile will perform the same request sequence as if you visited a profile
//   in the app. It will first call Instagram.Search(user), then register the click,
//   and lastly visit the profile with User.VisitProfile() and gather (some) posts
//   from the user feed, stories, grab the friendship status, and if available IGTV posts.
//
// You can access the profile info from the profile struct by calling Profile.Feed,
//   Profile.Stories, Profile.User etc. See the Profile struct for all properties.
//
func (insta *Instagram) VisitProfile(handle string) (*Profile, error) {
	sr, err := insta.Search(handle)
	if err != nil {
		return nil, err
	}
	for _, r := range sr.Results {
		if r.User.Username == handle {
			err = r.RegisterClick()
			if err != nil {
				return nil, err
			}
			return r.User.VisitProfile()
		}
	}
	return nil, errors.New("Profile not found")
}

// VisitProfile will perform the same request sequence as if you visited a profile
//   in the app. Thus it will gather (some) posts from the user feed, stories,
//   grab the friendship status, and if available IGTV posts.
//
// You can access the profile info from the profile struct by calling Profile.Feed,
//   Profile.Stories, Profile.User etc. See the Profile struct for all properties.
//
// This method will visit a profile directly from an already existing User instance.
// To visit a profile without the User struct, you can use Insta.VisitProfile(user),
//   which will perform a search, register the click, and call this method.
//
func (user *User) VisitProfile() (*Profile, error) {
	p := Profile{User: user}

	// Fetch Friendship
	fr, err := user.GetFriendship()
	if err != nil {
		return nil, err
	}
	p.Friendship = fr

	// Fetch Feed
	// usually gets called 3 times on profile visit, if enough media available
	feed := user.Feed()
	for i := 0; i < 3; i++ {
		s := feed.Next()
		err := feed.Error()
		if !s && err != ErrNoMore {
			return nil, err
		}
		time.Sleep(200 * time.Millisecond)
	}
	p.Feed = *feed

	// Fetch Stories
	stories, err := user.Stories()
	if err != nil {
		return nil, err
	}
	p.Stories = *stories

	// Fetch Highlights
	h, err := user.Highlights()
	if err != nil {
		return nil, err
	}
	p.Highlights = h

	// Fetch Profile Info
	err = user.Info("entry_point", "profile", "from_module", "blended_search")
	if err != nil {
		return nil, err
	}

	// Fetch Featured Accounts
	_, err = user.GetFeaturedAccounts()
	if err != nil {
		user.insta.ErrHandler(err)
	}

	// Fetch IGTV
	if user.HasIGTVSeries {
		igtv, err := user.IGTV()
		if err != nil {
			return nil, err
		}
		p.IGTV = *igtv
	}
	return &p, nil
}

func newProfiles(insta *Instagram) *Profiles {
	profiles := &Profiles{
		insta: insta,
	}
	return profiles
}

// ByName return a *User structure parsed by username.
// This is not the preffered method to fetch a profile, as the app will
//   not simply call this endpoint. It is better to use insta.Search(user),
//   or insta.Searchbar.SearchUser(user).
//
func (prof *Profiles) ByName(name string) (*User, error) {
	body, err := prof.insta.sendSimpleRequest(urlUserByName, name)
	if err == nil {
		resp := userResp{}
		err = json.Unmarshal(body, &resp)
		if err == nil {
			user := &resp.User
			user.insta = prof.insta
			return user, err
		}
	}
	return nil, err
}

// GetUserByID is a wrapper for Profiles.ByID(id)
func (insta *Instagram) GetUserByID(id interface{}) (*User, error) {
	return insta.Profiles.ByID(id)
}

// ByID returns a *User structure parsed by user id.
func (prof *Profiles) ByID(id_ interface{}) (*User, error) {
	var id string
	switch x := id_.(type) {
	case int64:
		id = fmt.Sprintf("%d", x)
	case int:
		id = fmt.Sprintf("%d", x)
	case string:
		id = x
	default:
		return nil, errors.New("Invalid id, please provide a string or int(64)")
	}

	body, _, err := prof.insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlUserByID, id),
		},
	)
	if err == nil {
		resp := userResp{}
		err = json.Unmarshal(body, &resp)
		if err == nil {
			user := &resp.User
			user.insta = prof.insta
			return user, err
		}
	}
	return nil, err
}

// Blocked returns a list of users you have blocked.
func (prof *Profiles) Blocked() ([]BlockedUser, error) {
	body, err := prof.insta.sendSimpleRequest(urlBlockedList)
	if err == nil {
		resp := blockedListResp{}
		err = json.Unmarshal(body, &resp)
		return resp.BlockedList, err
	}
	return nil, err
}
