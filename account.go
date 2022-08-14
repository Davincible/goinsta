package goinsta

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Davincible/goinsta/v3/utilities"
)

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
	if err != nil {
		return err
	}

	resp := profResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}

	*account = resp.Account
	account.insta = insta
	return nil
}

// ChangePassword changes current password.
//
// GoInsta does not store current instagram password (for security reasons)
// If you want to change your password you must parse old and new password.
//
// See example: examples/account/changePass.go
func (account *Account) ChangePassword(old, new_ string) error {
	insta := account.insta

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	old, err := utilities.EncryptPassword(old, insta.pubKey, insta.pubKeyID, timestamp)
	if err != nil {
		return err
	}
	new_, err = utilities.EncryptPassword(new_, insta.pubKey, insta.pubKeyID, timestamp)
	if err != nil {
		return err
	}

	data, err := json.Marshal(
		map[string]string{
			"_uid":              toString(insta.Account.ID),
			"_uuid":             insta.uuid,
			"enc_old_password":  old,
			"enc_new_password1": new_,
			"enc_new_password2": new_,
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
	data, err := json.Marshal(
		map[string]string{
			"_uid":  toString(insta.Account.ID),
			"_uuid": insta.uuid,
		},
	)
	if err != nil {
		return err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlRemoveProfPic,
			IsPost:   true,
			Query:    generateSignature(data),
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

	buf, err := readFile(photo)
	if err != nil {
		return fmt.Errorf("ChangeProfilePic readFile: %w", err)
	}

	o := UploadOptions{
		insta: insta,
		File:  photo,
		buf:   buf,
	}

	t := http.DetectContentType(buf.Bytes())
	if t != "image/jpeg" {
		return ErrInvalidImage
	}

	if err = o.uploadPhoto(); err != nil {
		return fmt.Errorf("ChangeProfilePic uploadPhoto: %w", err)
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlChangeProfPic,
			Query: map[string]string{
				"_uuid":          insta.uuid,
				"use_fbuploader": "true",
				"upload_id":      o.uploadID,
			},
			IsPost: true,
		},
	)
	if err != nil {
		return fmt.Errorf("ChangeProfilePic unmarshal response json: %w", err)
	}

	resp := profResp{}
	if err = json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("Failed to unmarshal account from json resposne: %w", err)
	}
	*account = resp.Account
	account.insta = insta
	return nil
}

// SetPrivate sets account to private mode.
//
// This function updates current Account information.
//
// See example: examples/account/setPrivate.go
func (account *Account) SetPrivate() error {
	return account.changePublic(urlSetPrivate)
}

// SetPublic sets account to public mode.
//
// This function updates current Account information.
//
// See example: examples/account/setPublic.go
func (account *Account) SetPublic() error {
	return account.changePublic(urlSetPublic)
}

func (account *Account) changePublic(endpoint string) error {
	insta := account.insta
	data, err := json.Marshal(
		map[string]string{
			"_uid":  toString(insta.Account.ID),
			"_uuid": insta.uuid,
		})
	if err != nil {
		return err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: endpoint,
			IsPost:   true,
			Query:    generateSignature(data),
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

	media := &FeedMedia{
		insta:    insta,
		endpoint: urlUserFeed,
		uid:      account.ID,
	}

	for _, param := range params {
		switch s := param.(type) {
		case string:
			media.timestamp = s
		}
	}

	return media
}

// Stories returns account stories.
func (account *Account) Stories() (*StoryMedia, error) {
	return account.insta.fetchStories(account.ID)
}

// Tags returns media where account is tagged in
//
// For pagination use FeedMedia.Next()
func (account *Account) Tags(minTimestamp []byte) (*FeedMedia, error) {
	insta := account.insta

	timestamp := string(minTimestamp)
	body, _, err := insta.sendRequest(
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

	media := &FeedMedia{
		insta:    insta,
		endpoint: urlUserTags,
		uid:      account.ID,
	}
	err = json.Unmarshal(body, media)
	return media, err
}

// Saved returns saved media.
// To get all the media you have to
// use the Next() method.
func (account *Account) Saved() *SavedMedia {
	return &SavedMedia{
		insta:    account.insta,
		endpoint: urlFeedSavedPosts,
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

// UpdateProfile method allows you to update your current account information.
// :param: form takes a map[string]string, the common values are:
//
// external_url
// phone_number
// username
// first_name -- is actually your full name
// biography
// email
//
func (account *Account) UpdateProfile(form map[string]string) error {
	insta := account.insta
	query := map[string]string{
		"external_url": "",
		"phone_number": "",
		"username":     insta.Account.Username,
		"first_name":   insta.Account.FullName,
		"_uid":         toString(insta.Account.ID),
		"device_id":    insta.dID,
		"biography":    insta.Account.Biography,
		"_uuid":        insta.uuid,
		"email":        insta.Account.Email,
	}

	for k, v := range form {
		query[k] = v
	}
	data, err := json.Marshal(query)
	if err != nil {
		return err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlEditProfile,
			IsPost:   true,
			Query:    generateSignature(data),
		},
	)
	if err != nil {
		return err
	}
	resp := struct {
		Status string   `json:"status"`
		User   *Account `json:"user"`
	}{
		User: insta.Account,
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}
	if resp.Status != "ok" {
		return fmt.Errorf("Can't update profile")
	}
	insta.Account = resp.User
	return nil
}

// EditBiography changes your Instagram's biography.
//
// This function updates current Account information.
func (account *Account) EditBiography(bio string) error {
	return account.UpdateProfile(map[string]string{"biography": bio})
}

// EditName changes your Instagram account name.
//
// This function updates current Account information.
func (account *Account) EditName(name string) error {
	return account.UpdateProfile(map[string]string{"first_name": name})
}

// EditUrl changes your Instagram account url.
//
// This function updates current Account information.
func (account *Account) EditUrl(url string) error {
	return account.UpdateProfile(map[string]string{"external_url": url})
}

// Liked are liked publications
func (account *Account) Liked() *FeedMedia {
	insta := account.insta

	media := &FeedMedia{
		insta:    insta,
		endpoint: urlFeedLiked,
	}
	return media
}

// PendingFollowRequests returns pending follow requests.
func (account *Account) PendingFollowRequests() (*PendingRequests, error) {
	insta := account.insta
	resp, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlFriendshipPending,
		},
	)
	if err != nil {
		return nil, err
	}

	var result PendingRequests
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}
	if result.Status != "ok" {
		return nil, fmt.Errorf("bad status: %s", result.Status)
	}

	var users []string
	for _, u := range result.Users {
		u.insta = insta
		users = append(users, toString(u.ID))
	}
	for _, u := range result.SuggestedUsers.Suggestions {
		u.User.insta = insta
	}

	friendships, err := account.FriendhipsShowMany(users)
	if err != nil {
		return nil, err
	}

	for _, u := range result.Users {
		if val, ok := friendships[toString(u.ID)]; ok {
			u.Friendship = val
		}
	}

	return &result, nil
}

// PendingRequestCount returns the number of open pending friendships as int
func (account *Account) PendingRequestCount() (int, error) {
	insta := account.insta
	resp, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlFriendshipPendingCount,
		},
	)
	if err != nil {
		return 0, err
	}

	var result struct {
		Count  int    `json:"count"`
		Status string `json:"status"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return 0, err
	}
	if result.Status != "ok" {
		return 0, fmt.Errorf("bad status: %s", result.Status)
	}
	return result.Count, nil
}

func (account *Account) FriendhipsShowMany(userIds []string) (map[string]Friendship, error) {
	insta := account.insta
	resp, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlFriendshipShowMany,
			IsPost:   true,
			Query: map[string]string{
				"user_ids": strings.Join(userIds, ","),
				"_uuid":    insta.uuid,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	var result struct {
		Friendships map[string]Friendship `json:"friendship_statuses"`
		Status      string                `json:"status"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}
	if result.Status != "ok" {
		return nil, fmt.Errorf("bad status: %s", result.Status)
	}
	return result.Friendships, nil
}

// Archived returns current account archive feed
//
// For pagination use FeedMedia.Next()
func (account *Account) Archived(params ...interface{}) *FeedMedia {
	insta := account.insta

	media := &FeedMedia{
		insta:    insta,
		endpoint: urlUserArchived,
	}

	for _, param := range params {
		switch s := param.(type) {
		case string:
			media.timestamp = s
		}
	}

	return media
}
