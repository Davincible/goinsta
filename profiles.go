package goinsta

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Profiles allows user function interactions
type Profiles struct {
	insta *Instagram
}

func (insta *Instagram) VisitProfile(handle string) {
}

func newProfiles(insta *Instagram) *Profiles {
	profiles := &Profiles{
		insta: insta,
	}
	return profiles
}

// ByName return a *User structure parsed by username
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

func (insta *Instagram) GetUserByID(id interface{}) (*User, error) {
	return insta.Profiles.ByID(id)
}

// ByID returns a *User structure parsed by user id
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

// Blocked returns a list of blocked profiles.
func (prof *Profiles) Blocked() ([]BlockedUser, error) {
	body, err := prof.insta.sendSimpleRequest(urlBlockedList)
	if err == nil {
		resp := blockedListResp{}
		err = json.Unmarshal(body, &resp)
		return resp.BlockedList, err
	}
	return nil, err
}
