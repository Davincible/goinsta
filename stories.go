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
	insta    *Instagram
	endpoint string
	uid      int64

	err error

	Pk                     interface{} `json:"id"`
	MediaCount             int64       `json:"media_count"`
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
	User                   User        `json:"user"`
	Items                  []Item      `json:"items"`
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
	StoryRankingToken    string      `json:"story_ranking_token"`
	Broadcasts           []Broadcast `json:"broadcasts"`
	FaceFilterNuxVersion int         `json:"face_filter_nux_version"`
	HasNewNuxStory       bool        `json:"has_new_nux_story"`
	Status               string      `json:"status"`
}

// Deletes ALL user's instagram stories.
// If you want to remove a single story, pick one from StoryMedia.Items, and
//   call Item.Delete()
//
// See example: examples/media/deleteStories.go
func (media *StoryMedia) Delete() error {
	for _, item := range media.Items {
		err := item.Delete()
		if err != nil {
			return err
		}
		time.Sleep(200 * time.Millisecond)

	}
	return nil
}

// ID returns Story id
func (media *StoryMedia) GetNextID() string {
	return formatID(media.Pk)
}

func (media *StoryMedia) setValues() {
	for i := range media.Items {
		media.Items[i].insta = media.insta
		setToItem(&media.Items[i], media)
	}
}

// Error returns error happened any error
func (media *StoryMedia) Error() error {
	return media.err
}

func (media *StoryMedia) getInsta() *Instagram {
	return media.insta
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

// Sync function is used when Highlight must be sync.
// Highlight must be sync when User.Highlights does not return any object inside StoryMedia slice.
//
// This function does NOT update Stories items.
//
// This function updates (fetches) StoryMedia.Items
func (media *StoryMedia) Sync() error {
	insta := media.insta
	supCap, err := getSupCap()
	if err != nil {
		return err
	}

	id := media.Pk.(string)
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
	if err == nil {
		resp := trayResp{}
		err = json.Unmarshal(body, &resp)
		if err == nil {
			m, ok := resp.Reels[id]
			if ok {
				media.Items = m.Items
				media.setValues()
				return nil
			}
			err = fmt.Errorf("cannot find %s structure in response", id)
		}
	}
	return err
}

// Next allows pagination after calling:
// User.Stories
//
//
// returns false when list reach the end
// if StoryMedia.Error() is ErrNoMore no problem have been occurred.
func (media *StoryMedia) Next(params ...interface{}) bool {
	if media.err != nil {
		return false
	}

	insta := media.insta
	endpoint := media.endpoint
	if media.uid != 0 {
		endpoint = fmt.Sprintf(endpoint, media.uid)
	}

	body, err := insta.sendSimpleRequest(endpoint)
	if err == nil {
		m := StoryMedia{}
		err = json.Unmarshal(body, &m)
		if err == nil {
			// TODO check NextID media
			*media = m
			media.insta = insta
			media.endpoint = endpoint
			media.err = ErrNoMore // TODO: See if stories has pagination
			media.setValues()
			return true
		}
	}
	media.err = err
	return false
}
