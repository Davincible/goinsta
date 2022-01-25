package goinsta

import (
	"encoding/json"
	"fmt"
	"time"
)

// StoryReelMention represent story reel mention
type StoryReelMention struct {
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Z           int     `json:"z"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	Rotation    float64 `json:"rotation"`
	IsPinned    int     `json:"is_pinned"`
	IsHidden    int     `json:"is_hidden"`
	IsSticker   int     `json:"is_sticker"`
	IsFBSticker int     `json:"is_fb_sticker"`
	User        User
	DisplayType string `json:"display_type"`
}

// StoryCTA represent story cta
type StoryCTA struct {
	Links []struct {
		LinkType                                int         `json:"linkType"`
		WebURI                                  string      `json:"webUri"`
		AndroidClass                            string      `json:"androidClass"`
		Package                                 string      `json:"package"`
		DeeplinkURI                             string      `json:"deeplinkUri"`
		CallToActionTitle                       string      `json:"callToActionTitle"`
		RedirectURI                             interface{} `json:"redirectUri"`
		LeadGenFormID                           string      `json:"leadGenFormId"`
		IgUserID                                string      `json:"igUserId"`
		AppInstallObjectiveInvalidationBehavior interface{} `json:"appInstallObjectiveInvalidationBehavior"`
	} `json:"links"`
}

// StoryMedia is the struct that handles the information from the methods to get info about Stories.
type StoryMedia struct {
	Reel       Reel         `json:"reel"`
	Broadcast  *Broadcast   `json:"broadcast"`
	Broadcasts []*Broadcast `json:"broadcasts"`
	Status     string       `json:"status"`
}

// Reel represents a single user's story collection.
// Every user has one reel, and one reel can contain many story items
type Reel struct {
	insta *Instagram

	ID                     interface{} `json:"id"`
	Items                  []*Item     `json:"items"`
	MediaCount             int         `json:"media_count"`
	MediaIDs               []int64     `json:"media_ids"`
	Muted                  bool        `json:"muted"`
	LatestReelMedia        int64       `json:"latest_reel_media"`
	LatestBestiesReelMedia float64     `json:"latest_besties_reel_media"`
	ExpiringAt             float64     `json:"expiring_at"`
	Seen                   float64     `json:"seen"`
	SeenRankedPosition     int         `json:"seen_ranked_position"`
	CanReply               bool        `json:"can_reply"`
	CanGifQuickReply       bool        `json:"can_gif_quick_reply"`
	ClientPrefetchScore    float64     `json:"client_prefetch_score"`
	Title                  string      `json:"title"`
	CanReshare             bool        `json:"can_reshare"`
	ReelType               string      `json:"reel_type"`
	ReelMentions           []string    `json:"reel_mentions"`
	PrefetchCount          int         `json:"prefetch_count"`
	// this field can be int or bool
	HasBestiesMedia       interface{} `json:"has_besties_media"`
	HasPrideMedia         bool        `json:"has_pride_media"`
	HasVideo              bool        `json:"has_video"`
	IsCacheable           bool        `json:"is_cacheable"`
	IsSensitiveVerticalAd bool        `json:"is_sensitive_vertical_ad"`
	RankedPosition        int         `json:"ranked_position"`
	RankerScores          struct {
		Fp   float64 `json:"fp"`
		Ptap float64 `json:"ptap"`
		Vm   float64 `json:"vm"`
	} `json:"ranker_scores"`
	StoryRankingToken    string `json:"story_ranking_token"`
	FaceFilterNuxVersion int    `json:"face_filter_nux_version"`
	HasNewNuxStory       bool   `json:"has_new_nux_story"`
	User                 User   `json:"user"`
}

// Stories will fetch a user's stories.
func (user *User) Stories() (*StoryMedia, error) {
	return user.insta.fetchStories(user.ID)
}

// Highlights will fetch a user's highlights.
func (user *User) Highlights() ([]*Reel, error) {
	data, err := getSupCap()
	if err != nil {
		return nil, err
	}

	body, _, err := user.insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlUserHighlights, user.ID),
			Query:    map[string]string{"supported_capabilities_new": data},
		},
	)
	if err != nil {
		return nil, err
	}

	tray := &Tray{}
	err = json.Unmarshal(body, &tray)
	if err != nil {
		return nil, err
	}

	tray.set(user.insta)
	return tray.Stories, nil
}

// Deletes ALL user's instagram stories.
// If you want to remove a single story, pick one from StoryMedia.Items, and
//   call Item.Delete()
//
// See example: examples/media/deleteStories.go
func (media *Reel) Delete() error {
	// TODO: update to reel
	for _, item := range media.Items {
		err := item.Delete()
		if err != nil {
			return err
		}
		time.Sleep(200 * time.Millisecond)

	}
	return nil
}

func (media *StoryMedia) setValues(insta *Instagram) {
	media.Reel.setValues(insta)
	if media.Broadcast != nil {
		media.Broadcast.setValues(insta)
	}
	for _, br := range media.Broadcasts {
		br.setValues(insta)
	}
}

func (media *Reel) setValues(insta *Instagram) {
	media.insta = insta
	media.User.insta = insta
	for _, i := range media.Items {
		i.insta = insta
		i.User.insta = insta
	}
}

// Seen marks story as seen.
/*
func (media *StoryMedia) Seen() error {
	insta := media.inst
	data, err := insta.prepareData(
		map[string]interface{}{
			"container_module":   "feed_timeline",
			"live_vods_skipped":  "",
			"nuxes_skipped":      "",
			"nuxes":              "",
			"reels":              "", // TODO xd
			"live_vods":          "",
			"reel_media_skipped": "",
		},
	)
	if err == nil {
		_, _, err = insta.sendRequest(
			&reqOptions{
				Endpoint: urlMediaSeen, // reel=1&live_vod=0
				Query:    generateSignature(data),
				IsPost:   true,
				UseV2:    true,
			},
		)
	}
	return err
}
*/

type trayRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (insta *Instagram) fetchStories(id int64) (*StoryMedia, error) {
	supCap, err := getSupCap()
	if err != nil {
		return nil, err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlUserStories, id),
			Query:    map[string]string{"supported_capabilities_new": supCap},
		},
	)
	if err != nil {
		return nil, err
	}

	m := &StoryMedia{}
	err = json.Unmarshal(body, m)
	if err != nil {
		return nil, err
	}

	m.setValues(insta)
	return m, nil
}

// Sync function is used when Highlights must be sync.
// Highlight must be sync when User.Highlights does not return any object inside StoryMedia slice.
//
// This function does NOT update Stories items.
//
// This function updates (fetches) StoryMedia.Items
func (media *Reel) Sync() error {
	if media.ReelType != "highlight_reel" {
		return ErrNotHighlight
	}

	insta := media.insta
	supCap, err := getSupCap()
	if err != nil {
		return err
	}

	id := media.ID.(string)
	data, err := json.Marshal(
		map[string]interface{}{
			"exclude_media_ids":          "[]",
			"supported_capabilities_new": supCap,
			"source":                     "reel_feed_timeline",
			"_uid":                       toString(insta.Account.ID),
			"_uuid":                      insta.uuid,
			"user_ids":                   []string{id},
		},
	)
	if err != nil {
		return err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlReelMedia,
			IsPost:   true,
			Query:    generateSignature(data),
		},
	)
	if err != nil {
		return err
	}

	resp := trayResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}

	m, ok := resp.Reels[id]
	if ok {
		*media = m
		media.setValues(insta)
		return nil
	}
	return fmt.Errorf("cannot find %s structure in response", id)
}
