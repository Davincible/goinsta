package goinsta

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
)

// Broadcast struct represents live video streams.
type Broadcast struct {
	insta *Instagram
	// Mutex to keep track of BroadcastStatus
	mu *sync.RWMutex

	LastLikeTs         int64
	LastCommentTs      int64
	LastCommentFetchTs int64
	LastCommentTotal   int

	ID         int64  `json:"id"`
	MediaID    string `json:"media_id"`
	LivePostID int64  `json:"live_post_id"`

	// BroadcastStatus is either "active", "interrupted", "stopped"
	BroadcastStatus           string  `json:"broadcast_status"`
	DashPlaybackURL           string  `json:"dash_playback_url"`
	DashAbrPlaybackURL        string  `json:"dash_abr_playback_url"`
	DashManifest              string  `json:"dash_manifest"`
	ExpireAt                  int64   `json:"expire_at"`
	EncodingTag               string  `json:"encoding_tag"`
	InternalOnly              bool    `json:"internal_only"`
	NumberOfQualities         int     `json:"number_of_qualities"`
	CoverFrameURL             string  `json:"cover_frame_url"`
	User                      User    `json:"broadcast_owner"`
	Cobroadcasters            []*User `json:"cobroadcasters"`
	PublishedTime             int64   `json:"published_time"`
	Message                   string  `json:"broadcast_message"`
	OrganicTrackingToken      string  `json:"organic_tracking_token"`
	IsPlayerLiveTrace         int     `json:"is_player_live_trace_enabled"`
	IsGamingContent           bool    `json:"is_gaming_content"`
	IsViewerCommentAllowed    bool    `json:"is_viewer_comment_allowed"`
	IsPolicyViolation         bool    `json:"is_policy_violation"`
	PolicyViolationReason     string  `json:"policy_violation_reason"`
	LiveCommentMentionEnabled bool    `json:"is_live_comment_mention_enabled"`
	LiveCommentRepliesEnabled bool    `json:"is_live_comment_replies_enabled"`
	HideFromFeedUnit          bool    `json:"hide_from_feed_unit"`
	VideoDuration             float64 `json:"video_duration"`
	Visibility                int     `json:"visibility"`
	ViewerCount               float64 `json:"viewer_count"`
	ResponseTs                int64   `json:"response_timestamp"`
	Status                    string  `json:"status"`
	Dimensions                struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"dimensions"`
	Experiments     map[string]interface{} `json:"broadcast_experiments"`
	PayViewerConfig struct {
		PayConfig struct {
			ConsumptionSheetConfig struct {
				Description               string `json:"description"`
				PrivacyDisclaimer         string `json:"privacy_disclaimer"`
				PrivacyDisclaimerLink     string `json:"privacy_disclaimer_link"`
				PrivacyDisclaimerLinkText string `json:"privacy_disclaimer_link_text"`
			} `json:"consumption_sheet_config"`
			DigitalNonConsumableProductID int64 `json:"digital_non_consumable_product_id"`
			DigitalProductID              int64 `json:"digital_product_id"`
			PayeeID                       int64 `json:"payee_id"`
			PinnedRowConfig               struct {
				ButtonTitle string `json:"button_title"`
				Description string `json:"description"`
			} `json:"pinned_row_config"`
			TierInfos []struct {
				DigitalProductID int64  `json:"digital_product_id"`
				Sku              string `json:"sku"`
				SupportTier      string `json:"support_tier"`
			} `json:"tier_infos"`
		} `json:"pay_config"`
	} `json:"user_pay_viewer_config"`
}

type BroadcastComments struct {
	CommentLikesEnabled        bool       `json:"comment_likes_enabled"`
	Comments                   []*Comment `json:"comments"`
	PinnedComment              *Comment   `json:"pinned_comment"`
	CommentCount               int        `json:"comment_count"`
	Caption                    *Caption   `json:"caption"`
	CaptionIsEdited            bool       `json:"caption_is_edited"`
	HasMoreComments            bool       `json:"has_more_comments"`
	HasMoreHeadloadComments    bool       `json:"has_more_headload_comments"`
	MediaHeaderDisplay         string     `json:"media_header_display"`
	CanViewMorePreviewComments bool       `json:"can_view_more_preview_comments"`
	LiveSecondsPerComment      int        `json:"live_seconds_per_comment"`
	IsFirstFetch               string     `json:"is_first_fetch"`
	SystemComments             []*Comment `json:"system_comments"`
	CommentMuted               int        `json:"comment_muted"`
	IsViewerCommentAllowed     bool       `json:"is_viewer_comment_allowed"`
	UserPaySupportersInfo      struct {
		SupportersInComments   map[string]interface{} `json:"supporters_in_comments"`
		SupportersInCommentsV2 map[string]interface{} `json:"supporters_in_comments_v2"`
		// SupportersInCommentsV2 map[string]struct {
		// 	SupportTier string `json:"support_tier"`
		// 	BadgesCount int    `json:"badges_count"`
		// } `json:"supporters_in_comments_v2"`
		NewSupportersNextMinID int64          `json:"new_supporters_next_min_id"`
		NewSupporters          []NewSupporter `json:"new_supporters"`
	} `json:"user_pay_supporter_info"`
	Status string `json:"status"`
}

type NewSupporter struct {
	RepeatedSupporter bool    `json:"is_repeat_supporter"`
	SupportTier       string  `json:"support_tier"`
	Timestamp         float64 `json:"ts_secs"`
	User              struct {
		ID         int64  `json:"pk"`
		Username   string `json:"username"`
		FullName   string `json:"full_name"`
		IsPrivate  bool   `json:"is_private"`
		IsVerified bool   `json:"is_verified"`
	}
}

type BroadcastLikes struct {
	Likes      int `json:"likes"`
	BurstLikes int `json:"burst_likes"`
	Likers     []struct {
		UserID        int64  `json:"user_id"`
		ProfilePicUrl string `json:"profile_pic_url"`
		Count         string `json:"count"`
	} `json:"likers"`
	LikeTs           int64  `json:"like_ts"`
	Status           string `json:"status"`
	PaySupporterInfo struct {
		LikeCountByTier []struct {
			BurstLikes  int           `json:"burst_likes"`
			Likers      []interface{} `json:"likers"`
			Likes       int           `json:"likes"`
			SupportTier string        `json:"support_tier"`
		} `json:"like_count_by_support_tier"`
		BurstLikes int `json:"supporter_tier_burst_likes"`
		Likes      int `json:"supporter_tier_likes"`
	} `json:"user_pay_supporter_info"`
}

type BroadcastHeartbeat struct {
	ViewerCount             float64  `json:"viewer_count"`
	BroadcastStatus         string   `json:"broadcast_status"`
	CobroadcasterIds        []string `json:"cobroadcaster_ids"`
	OffsetVideoStart        float64  `json:"offset_to_video_start"`
	RequestToJoinEnabled    int      `json:"request_to_join_enabled"`
	UserPayMaxAmountReached bool     `json:"user_pay_max_amount_reached"`
	Status                  string   `json:"status"`
}

// Discover wraps Instagram.IGTV.Live
func (br *Broadcast) Discover() (*IGTVChannel, error) {
	return br.insta.IGTV.Live()
}

// NewUser returns prepared user to be used with his functions.
func (insta *Instagram) NewBroadcast(id int64) *Broadcast {
	return &Broadcast{insta: insta, ID: id, mu: &sync.RWMutex{}}
}

// GetInfo will fetch the information about a broadcast
func (br *Broadcast) GetInfo() error {
	body, _, err := br.insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlLiveInfo, br.ID),
			Query: map[string]string{
				"view_expired_broadcast": "false",
			},
		})
	if err != nil {
		return err
	}

	br.mu.RLock()
	err = json.Unmarshal(body, br)
	br.mu.RUnlock()
	return err
}

// Call every 2 seconds
func (br *Broadcast) GetComments() (*BroadcastComments, error) {
	br.mu.RLock()
	if br.BroadcastStatus == "stopped" {
		return nil, ErrMediaDeleted
	}
	br.mu.RUnlock()

	body, _, err := br.insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlLiveComments, br.ID),
			Query: map[string]string{
				"last_comment_ts":               strconv.FormatInt(br.LastCommentTs, 10),
				"join_request_last_seen_ts":     "0",
				"join_request_last_fetch_ts":    "0",
				"join_request_last_total_count": "0",
			},
		})
	if err != nil {
		return nil, err
	}

	c := &BroadcastComments{}
	err = json.Unmarshal(body, c)
	if err != nil {
		return nil, err
	}

	if c.CommentCount > 0 {
		br.LastCommentTs = c.Comments[0].CreatedAt
	}
	return c, nil
}

// Call every 6 seconds
func (br *Broadcast) GetLikes() (*BroadcastLikes, error) {
	br.mu.RLock()
	if br.BroadcastStatus == "stopped" {
		return nil, ErrMediaDeleted
	}
	br.mu.RUnlock()

	body, _, err := br.insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlLiveLikeCount, br.ID),
			Query: map[string]string{
				"like_ts": strconv.FormatInt(br.LastLikeTs, 10),
			},
		})
	if err != nil {
		return nil, err
	}

	c := &BroadcastLikes{}
	err = json.Unmarshal(body, c)
	if err != nil {
		return nil, err
	}
	br.LastLikeTs = c.LikeTs
	return c, nil
}

// Call every 3 seconds
func (br *Broadcast) GetHeartbeat() (*BroadcastHeartbeat, error) {
	br.mu.RLock()
	if br.BroadcastStatus == "stopped" {
		return nil, ErrMediaDeleted
	}
	br.mu.RUnlock()

	body, _, err := br.insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlLiveHeartbeat, br.ID),
			IsPost:   true,
			Query: map[string]string{
				"_uuid":                 br.insta.uuid,
				"live_with_eligibility": "2", // What is this?
			},
		})
	if err != nil {
		return nil, err
	}

	c := &BroadcastHeartbeat{}
	err = json.Unmarshal(body, c)
	if err != nil {
		return nil, err
	}

	br.mu.RLock()
	br.BroadcastStatus = c.BroadcastStatus
	br.mu.RUnlock()

	return c, nil
}

// GetLiveChaining traditionally gets called after the live stream has ended, and provides
//   recommendations of other current live streams, as well as past live streams.
func (br *Broadcast) GetLiveChaining() ([]*Broadcast, error) {
	insta := br.insta
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlLiveChaining,
			Query: map[string]string{
				"include_post_lives": "true",
			},
		},
	)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Broadcasts []*Broadcast `json:"broadcasts"`
		Status     string       `json:"string"`
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	for _, br := range resp.Broadcasts {
		br.setValues(insta)
	}
	return resp.Broadcasts, nil
}

func (br *Broadcast) DownloadCoverFrame() ([]byte, error) {
	if br.CoverFrameURL == "" {
		return nil, ErrNoMedia
	}

	b, err := br.insta.download(br.CoverFrameURL)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (br *Broadcast) setValues(insta *Instagram) {
	br.insta = insta
	br.User.insta = insta
	br.mu = &sync.RWMutex{}
	for _, cb := range br.Cobroadcasters {
		cb.insta = insta
	}
}
