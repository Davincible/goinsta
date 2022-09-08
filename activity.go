package goinsta

import (
	"encoding/json"
	"fmt"
)

// Activity is the recent activity menu.
//
// See example: examples/activity/recent.go
type Activity struct {
	insta *Instagram
	err   error

	// Ad is every column of Activity section
	Ad struct {
		Items []struct {
			// User            User          `json:"user"`
			Algorithm       string        `json:"algorithm"`
			SocialContext   string        `json:"social_context"`
			Icon            string        `json:"icon"`
			Caption         string        `json:"caption"`
			MediaIds        []interface{} `json:"media_ids"`
			ThumbnailUrls   []interface{} `json:"thumbnail_urls"`
			LargeUrls       []interface{} `json:"large_urls"`
			MediaInfos      []interface{} `json:"media_infos"`
			Value           float64       `json:"value"`
			IsNewSuggestion bool          `json:"is_new_suggestion"`
		} `json:"items"`
		MoreAvailable bool `json:"more_available"`
	} `json:"aymf"`
	Counts struct {
		Campaign      int `json:"campaign_notification"`
		CommentLikes  int `json:"comment_likes"`
		Comments      int `json:"comments"`
		Fundraiser    int `json:"fundraiser"`
		Likes         int `json:"likes"`
		NewPosts      int `json:"new_posts"`
		PhotosOfYou   int `json:"photos_of_you"`
		Relationships int `json:"relationships"`
		Requests      int `json:"requests"`
		Shopping      int `json:"shopping_notification"`
		UserTags      int `json:"usertags"`
	} `json:"counts"`
	FriendRequestStories []interface{} `json:"friend_request_stories"`
	NewStories           []RecentItems `json:"new_stories"`
	OldStories           []RecentItems `json:"old_stories"`
	ContinuationToken    int64         `json:"continuation_token"`
	Subscription         interface{}   `json:"subscription"`
	NextID               string        `json:"next_max_id"`
	LastChecked          float64       `json:"last_checked"`
	FirstRecTs           float64       `json:"pagination_first_record_timestamp"`

	Status string `json:"status"`
}

type RecentItems struct {
	Type      int `json:"type"`
	StoryType int `json:"story_type"`
	Args      struct {
		Text     string `json:"text"`
		RichText string `json:"rich_text"`
		IconUrl  string `json:"icon_url"`
		Links    []struct {
			Start int         `json:"start"`
			End   int         `json:"end"`
			Type  string      `json:"type"`
			ID    interface{} `json:"id"`
		} `json:"links"`
		InlineFollow struct {
			UserInfo        User `json:"user_info"`
			Following       bool `json:"following"`
			OutgoingRequest bool `json:"outgoing_request"`
		} `json:"inline_follow"`
		Actions         []string    `json:"actions"`
		AfCandidateId   int         `json:"af_candidate_id"`
		ProfileID       int64       `json:"profile_id"`
		ProfileImage    string      `json:"profile_image"`
		Timestamp       float64     `json:"timestamp"`
		Tuuid           string      `json:"tuuid"`
		Clicked         bool        `json:"clicked"`
		ProfileName     string      `json:"profile_name"`
		LatestReelMedia int64       `json:"latest_reel_media"`
		Destination     string      `json:"destination"`
		Extra           interface{} `json:"extra"`
	} `json:"args"`
	Counts struct{} `json:"counts"`
	Pk     string   `json:"pk"`
}

func (act *Activity) Error() error {
	return act.err
}

// Next function allows pagination over notifications.
//
// See example: examples/activity/recent.go
func (act *Activity) Next() bool {
	if act.err != nil {
		return false
	}
	var first bool
	if act.Status == "" {
		first = true
	}

	query := map[string]string{
		"mark_as_seen":    "false",
		"timezone_offset": timeOffset,
	}
	if act.NextID != "" {
		query["max_id"] = act.NextID
		query["last_checked"] = fmt.Sprintf("%f", act.LastChecked)
		query["pagination_first_record_timestamp"] = fmt.Sprintf("%f", act.FirstRecTs)
	}

	insta := act.insta
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlActivityRecent,
			Query:    query,
		},
	)
	if err != nil {
		act.err = err
		return false
	}

	act2 := Activity{}
	err = json.Unmarshal(body, &act2)
	if err == nil {
		*act = act2
		act.insta = insta
		if first {
			if err := act.MarkAsSeen(); err != nil {
				act.err = err
				return false
			}
		}

		if act.NextID == "" {
			act.err = ErrNoMore
			return false
		}
		return true
	}
	act.err = err
	return false
}

// MarkAsSeen will let instagram know you visited the activity page, and mark
//   current items as seen.
func (act *Activity) MarkAsSeen() error {
	insta := act.insta
	_, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlActivitySeen,
			IsPost:   true,
			Query: map[string]string{
				"_uuid": insta.uuid,
			},
		},
	)
	return err
}

func newActivity(insta *Instagram) *Activity {
	act := &Activity{
		insta: insta,
	}
	return act
}
