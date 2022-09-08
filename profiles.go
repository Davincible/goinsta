package goinsta

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

// Profiles allows user function interactions
type Profiles struct {
	insta *Instagram
}

// Profile represents an instagram user with their various properties, such as
//   their account info, stored in Profile.User (a *User struct), feed, stories,
//   Highlights, IGTV posts, and friendship status.
//
type Profile struct {
	User       *User
	Friendship *Friendship

	Feed       *FeedMedia
	Stories    *StoryMedia
	Highlights []*Reel
	IGTV       *IGTVChannel
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
		if r.User != nil && r.User.Username == handle {
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
	p := &Profile{User: user}

	wg := &sync.WaitGroup{}
	info := &sync.WaitGroup{}
	errChan := make(chan error, 10)

	// Fetch Profile Info
	wg.Add(1)
	info.Add(1)
	go func(wg, info *sync.WaitGroup) {
		defer wg.Done()
		defer info.Done()
		if err := user.Info("entry_point", "profile", "from_module", "blended_search"); err != nil {
			errChan <- err
		}
	}(wg, info)

	// Fetch Friendship
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		fr, err := user.GetFriendship()
		if err != nil {
			errChan <- err
		}
		p.Friendship = fr
	}(wg)

	// Fetch Feed
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		feed := user.Feed()
		p.Feed = feed
		if !feed.Next() && feed.Error() != ErrNoMore {
			errChan <- feed.Error()
		}
	}(wg)

	// Fetch Stories
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		stories, err := user.Stories()
		p.Stories = stories
		if err != nil {
			errChan <- err
		}
	}(wg)

	// Fetch Highlights
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		h, err := user.Highlights()
		p.Highlights = h
		if err != nil {
			errChan <- err
		}
	}(wg)

	// Fetch Featured Accounts
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		_, err := user.GetFeaturedAccounts()
		if err != nil {
			user.insta.warnHandler(err)
		}
	}(wg)

	// Fetch IGTV
	wg.Add(1)
	go func(wg, info *sync.WaitGroup) {
		defer wg.Done()
		info.Wait()
		if user.IGTVCount > 0 {
			igtv, err := user.IGTV()
			if err != nil {
				errChan <- err
			}
			p.IGTV = igtv
		}
	}(wg, info)

	wg.Wait()
	select {
	case err := <-errChan:
		return p, err
	default:
		return p, nil
	}
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
		return nil, errors.New("invalid id, please provide a string or int(64)")
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
