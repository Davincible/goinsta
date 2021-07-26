package goinsta

import (
	"encoding/json"
	"fmt"
)

// Profiles allows user function interactions
type Profiles struct {
	insta *Instagram
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

// ByID returns a *User structure parsed by user id
func (prof *Profiles) ByID(id int64) (*User, error) {
	data, err := prof.insta.prepareData()
	if err != nil {
		return nil, err
	}

	body, _, err := prof.insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlUserByID, id),
			Query:    generateSignature(data),
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
