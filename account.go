package goinsta

import (
	"encoding/json"
	"fmt"
	"io"
)

type accountResp struct {
	Status  string  `json:"status"`
	Account Account `json:"logged_in_user"`
}

// Account is personal account object
//
// See examples: examples/account/*
type Account struct {
	insta *Instagram

	ID                         int64        `json:"pk"`
	Username                   string       `json:"username"`
	FullName                   string       `json:"full_name"`
	Biography                  string       `json:"biography"`
	ProfilePicURL              string       `json:"profile_pic_url"`
	Email                      string       `json:"email"`
	PhoneNumber                string       `json:"phone_number"`
	IsBusiness                 bool         `json:"is_business"`
	Gender                     int          `json:"gender"`
	ProfilePicID               string       `json:"profile_pic_id"`
	CanSeeOrganicInsights      bool         `json:"can_see_organic_insights"`
	ShowInsightsTerms          bool         `json:"show_insights_terms"`
	Nametag                    Nametag      `json:"nametag"`
	HasAnonymousProfilePicture bool         `json:"has_anonymous_profile_picture"`
	IsPrivate                  bool         `json:"is_private"`
	IsUnpublished              bool         `json:"is_unpublished"`
	AllowedCommenterType       string       `json:"allowed_commenter_type"`
	IsVerified                 bool         `json:"is_verified"`
	MediaCount                 int          `json:"media_count"`
	FollowerCount              int          `json:"follower_count"`
	FollowingCount             int          `json:"following_count"`
	GeoMediaCount              int          `json:"geo_media_count"`
	ExternalURL                string       `json:"external_url"`
	HasBiographyTranslation    bool         `json:"has_biography_translation"`
	ExternalLynxURL            string       `json:"external_lynx_url"`
	HdProfilePicURLInfo        PicURLInfo   `json:"hd_profile_pic_url_info"`
	HdProfilePicVersions       []PicURLInfo `json:"hd_profile_pic_versions"`
	UsertagsCount              int          `json:"usertags_count"`
	HasChaining                bool         `json:"has_chaining"`
	ReelAutoArchive            string       `json:"reel_auto_archive"`
	PublicEmail                string       `json:"public_email"`
	PublicPhoneNumber          string       `json:"public_phone_number"`
	PublicPhoneCountryCode     string       `json:"public_phone_country_code"`
	ContactPhoneNumber         string       `json:"contact_phone_number"`
	Byline                     string       `json:"byline"`
	SocialContext              string       `json:"social_context,omitempty"`
	SearchSocialContext        string       `json:"search_social_context,omitempty"`
	MutualFollowersCount       float64      `json:"mutual_followers_count"`
	LatestReelMedia            int64        `json:"latest_reel_media,omitempty"`
	CityID                     int64        `json:"city_id"`
	CityName                   string       `json:"city_name"`
	AddressStreet              string       `json:"address_street"`
	DirectMessaging            string       `json:"direct_messaging"`
	Latitude                   float64      `json:"latitude"`
	Longitude                  float64      `json:"longitude"`
	Category                   string       `json:"category"`
	BusinessContactMethod      string       `json:"business_contact_method"`
	IsCallToActionEnabled      bool         `json:"is_call_to_action_enabled"`
	FbPageCallToActionID       string       `json:"fb_page_call_to_action_id"`
	Zip                        string       `json:"zip"`
	AllowContactsSync          bool         `json:"allow_contacts_sync"`
	CanBoostPost               bool         `json:"can_boost_post"`
}

// Sync updates account information
func (account *Account) Sync() error {
	insta := account.insta
	body, _, err := insta.sendRequest(&reqOptions{
		Endpoint: urlCurrentUser,
		Query: map[string]string{
			"edit": "true",
		},
	})
	if err == nil {
		resp := profResp{}
		err = json.Unmarshal(body, &resp)
		if err == nil {
			*account = resp.Account
			account.insta = insta
		}
	}
	return err
}

// ChangePassword changes current password.
//
// GoInsta does not store current instagram password (for security reasons)
// If you want to change your password you must parse old and new password.
//
// See example: examples/account/changePass.go
func (account *Account) ChangePassword(old, new string) error {
	insta := account.insta
	data, err := insta.prepareData(
		map[string]interface{}{
			"old_password":  old,
			"new_password1": new,
			"new_password2": new,
		},
	)
	if err == nil {
		_, _, err = insta.sendRequest(
			&reqOptions{
				Endpoint: urlChangePass,
				Query:    generateSignature(data),
				IsPost:   true,
			},
		)
	}
	return err
}

type profResp struct {
	Status  string  `json:"status"`
	Account Account `json:"user"`
}

// RemoveProfilePic removes current profile picture
//
// This function updates current Account information.
//
// See example: examples/account/removeProfilePic.go
func (account *Account) RemoveProfilePic() error {
	insta := account.insta
	data, err := insta.prepareData()
	if err != nil {
		return err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlRemoveProfPic,
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err == nil {
		resp := profResp{}
		err = json.Unmarshal(body, &resp)
		if err == nil {
			*account = resp.Account
			account.insta = insta
		}
	}
	return err
}

// ChangeProfilePic Update profile picture
//
// See example: examples/account/change-profile-pic/main.go
func (account *Account) ChangeProfilePic(photo io.Reader) error {
	insta := account.insta
	config, err := insta.postPhoto(photo, "", 1, 1, false)
	if err != nil {
		return err
	}
	data, err := insta.prepareData(config)
	if err != nil {
		return err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlChangeProfPic,
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err == nil {
		resp := profResp{}
		err = json.Unmarshal(body, &resp)
		if err == nil {
			*account = resp.Account
			account.insta = insta
		}
	}
	return err
}

// SetPrivate sets account to private mode.
//
// This function updates current Account information.
//
// See example: examples/account/setPrivate.go
func (account *Account) SetPrivate() error {
	insta := account.insta
	data, err := insta.prepareData()
	if err != nil {
		return err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlSetPrivate,
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err == nil {
		resp := profResp{}
		err = json.Unmarshal(body, &resp)
		if err == nil {
			*account = resp.Account
			account.insta = insta
		}
	}
	return err
}

// SetPublic sets account to public mode.
//
// This function updates current Account information.
//
// See example: examples/account/setPublic.go
func (account *Account) SetPublic() error {
	insta := account.insta
	data, err := insta.prepareData()
	if err != nil {
		return err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlSetPublic,
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err == nil {
		resp := profResp{}
		err = json.Unmarshal(body, &resp)
		if err == nil {
			*account = resp.Account
			account.insta = insta
		}
	}
	return err
}

// Followers returns a list of user followers.
//
// Users.Next can be used to paginate
//
// See example: examples/account/followers.go
func (account *Account) Followers() *Users {
	endpoint := fmt.Sprintf(urlFollowers, account.ID)
	users := &Users{}
	users.insta = account.insta
	users.endpoint = endpoint
	return users
}

// Following returns a list of user following.
//
// Users.Next can be used to paginate
//
// See example: examples/account/following.go
func (account *Account) Following() *Users {
	endpoint := fmt.Sprintf(urlFollowing, account.ID)
	users := &Users{}
	users.insta = account.insta
	users.endpoint = endpoint
	return users
}

// Feed returns current account feed
//
// 	params can be:
// 		string: timestamp of the minimum media timestamp.
//
// minTime is the minimum timestamp of media.
//
// For pagination use FeedMedia.Next()
func (account *Account) Feed(params ...interface{}) *FeedMedia {
	insta := account.insta

	media := &FeedMedia{}
	media.insta = insta
	media.endpoint = urlUserFeed
	media.uid = account.ID

	for _, param := range params {
		switch s := param.(type) {
		case string:
			media.timestamp = s
		}
	}

	return media
}

// Stories returns account stories.
//
// Use StoryMedia.Next for pagination.
//
// See example: examples/account/stories.go
func (account *Account) Stories() *StoryMedia {
	media := &StoryMedia{}
	media.uid = account.ID
	media.insta = account.insta
	media.endpoint = urlUserStories
	return media
}

// Tags returns media where account is tagged in
//
// For pagination use FeedMedia.Next()
func (account *Account) Tags(minTimestamp []byte) (*FeedMedia, error) {
	timestamp := b2s(minTimestamp)
	body, _, err := account.insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlUserTags, account.ID),
			Query: map[string]string{
				"max_id":         "",
				"rank_token":     account.insta.rankToken,
				"min_timestamp":  timestamp,
				"ranked_content": "true",
			},
		},
	)
	if err != nil {
		return nil, err
	}

	media := &FeedMedia{}
	err = json.Unmarshal(body, media)
	media.insta = account.insta
	media.endpoint = urlUserTags
	media.uid = account.ID
	return media, err
}

// Saved returns saved media.
// To get all the media you have to
// use the Next() method.
func (account *Account) Saved() *SavedMedia {
	return &SavedMedia{
		insta:    account.insta,
		endpoint: urlFeedSaved,
		err:      nil,
	}
}

type editResp struct {
	Status  string  `json:"status"`
	Account Account `json:"user"`
}

func (account *Account) edit() {
	insta := account.insta
	acResp := editResp{}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlCurrentUser,
			Query: map[string]string{
				"edit": "true",
			},
		},
	)
	if err == nil {
		err = json.Unmarshal(body, &acResp)
		if err == nil {
			acResp.Account.insta = insta
			*account = acResp.Account
		}
	}
}

// UpdateProfile This function updates current Account information.
func (account *Account) UpdateProfile(
	editProfileForm map[string]interface{},
) error {
	insta := account.insta
	data, err := insta.prepareData(editProfileForm)
	if err != nil {
		return err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint:   urlEditProfile,
			IsPost:     true,
			Connection: "keep-alive",
			Query:      generateSignature(data),
		},
	)
	if err != nil {
		return err
	}
	var respEdit struct {
		Status string `json:"status"`
	}
	err = json.Unmarshal(body, &respEdit)
	if err != nil {
		return err
	}
	if respEdit.Status == "success" {
		return fmt.Errorf("Can't update profile")
	}
	return nil
}

// SetBiography changes your Instagram's biography.
//
// This function updates current Account information.
func (account *Account) SetBiography(bio string) error {
	account.edit() // preparing to edit
	insta := account.insta
	data, err := insta.prepareData(
		map[string]interface{}{
			"raw_text": bio,
		},
	)
	if err != nil {
		return err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlSetBiography,
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err == nil {
		var resp struct {
			User struct {
				Pk        int64  `json:"pk"`
				Biography string `json:"biography"`
			} `json:"user"`
			Status string `json:"status"`
		}
		err = json.Unmarshal(body, &resp)
		if err == nil {
			account.Biography = resp.User.Biography
		}
	}
	return err
}

// Liked are liked publications
func (account *Account) Liked() *FeedMedia {
	insta := account.insta

	media := &FeedMedia{}
	media.insta = insta
	media.endpoint = urlFeedLiked
	return media
}

// PendingFollowRequests returns pending follow requests.
func (account *Account) PendingFollowRequests() ([]User, error) {
	insta := account.insta
	resp, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlFriendshipPending,
		},
	)
	if err != nil {
		return nil, err
	}

	var result struct {
		Users []User `json:"users"`
		// TODO: pagination
		// TODO: SuggestedUsers
		Status string `json:"status"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}
	if result.Status != "ok" {
		return nil, fmt.Errorf("bad status: %s", result.Status)
	}

	return result.Users, nil
}

// Archived returns current account archive feed
//
// For pagination use FeedMedia.Next()
func (account *Account) Archived(params ...interface{}) *FeedMedia {
	insta := account.insta

	media := &FeedMedia{}
	media.insta = insta
	media.endpoint = urlUserArchived

	for _, param := range params {
		switch s := param.(type) {
		case string:
			media.timestamp = s
		}
	}

	return media
}
