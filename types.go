package goinsta

import (
	"encoding/json"
	"fmt"
)

// ConfigFile is a structure to store the session information so that can be exported or imported.
type ConfigFile struct {
	ID            int64             `json:"id"`
	User          string            `json:"username"`
	DeviceID      string            `json:"device_id"`
	FamilyID      string            `json:"family_id"`
	UUID          string            `json:"uuid"`
	RankToken     string            `json:"rank_token"`
	Token         string            `json:"token"`
	PhoneID       string            `json:"phone_id"`
	XmidExpiry    int64             `json:"xmid_expiry"`
	HeaderOptions map[string]string `json:"header_options"`
	Account       *Account          `json:"account"`
	Device        Device            `json:"device"`
	TOTP          *TOTP             `json:"totp"`
	SessionNonce  string            `json:"session"`
}

type Device struct {
	Manufacturer     string `json:"manufacturer"`
	Model            string `json:"model"`
	CodeName         string `json:"code_name"`
	AndroidVersion   int    `json:"android_version"`
	AndroidRelease   int    `json:"android_release"`
	ScreenDpi        string `json:"screen_dpi"`
	ScreenResolution string `json:"screen_resolution"`
	Chipset          string `json:"chipset"`
}

// School is void structure (yet). Whats this even for lol
type School struct{}

// PicURLInfo repre
type PicURLInfo struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

// ErrorN is general instagram error
type ErrorN struct {
	Message   string `json:"message"`
	Endpoint  string `json:"endpoint"`
	Status    string `json:"status"`
	ErrorType string `json:"error_type"`
}

// Error503 is instagram API error
type Error503 struct {
	Message string
}

func (e Error503) Error() string {
	return e.Message
}

func (e ErrorN) Error() string {
	return fmt.Sprintf(
		"Error while calling %s, status code %s: %s (%s)",
		e.Endpoint, e.Status, e.Message, e.ErrorType,
	)
}

// Error400 is error returned by HTTP 400 status code.
type Error400 struct {
	Checkpoint

	Status  string `json:"status"`
	Message string `json:"message"`

	ErrorType  string `json:"error_type"`
	ErrorTitle string `json:"error_title"`
	ErrorBody  string `json:"error_body"`

	// Status code
	Code int

	// The endpoint that returned the 400 status code
	Endpoint string `json:"endpoint"`

	Challenge         *Challenge     `json:"challenge"`
	TwoFactorRequired bool           `json:"two_factor_required"`
	TwoFactorInfo     *TwoFactorInfo `json:"two_factor_info"`

	// This is double, as also present inside TwoFactorInfo
	PhoneVerificationSettings phoneVerificationSettings `json:"phone_verification_settings"`

	Payload struct {
		ClientContext string `json:"client_context"`
		Message       string `json:"message"`
	} `json:"payload"`

	DebugInfo struct {
		Message   string `json:"string"`
		Retriable bool   `json:"retriable"`
		Type      string `json:"type"`
	} `json:"debug_info"`
}

func (e Error400) Error() string {
	var msg string
	if e.Payload.Message != "" {
		msg = e.Payload.Message
	}
	if e.DebugInfo.Message != "" {
		msg = e.DebugInfo.Message
	}
	if e.Message != "" {
		if msg != "" {
			msg += "; " + e.Message
		} else {
			msg = e.Message
		}
	}

	if e.Code == 0 {
		e.Code = 400
	}
	return fmt.Sprintf("Request Status Code %d: %s, %s", e.Code, e.Status, msg)
}

func (e *Error400) GetMessage() string {
	if e.ErrorType != "" {
		return e.ErrorType
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Challenge != nil && len(e.Challenge.Errors) > 0 {
		return e.Challenge.Errors[0]
	}
	return ""
}

// ChallengeError is error returned by HTTP 400 status code.
type ChallengeError struct {
	Challenge struct {
		URL               string `json:"url"`
		APIPath           string `json:"api_path"`
		HideWebviewHeader bool   `json:"hide_webview_header"`
		Lock              bool   `json:"lock"`
		Logout            bool   `json:"logout"`
		NativeFlow        bool   `json:"native_flow"`
	} `json:"challenge"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	ErrorType string `json:"error_type"`
}

func (e ChallengeError) Error() string {
	return fmt.Sprintf("%s: %s, %s", ErrChallengeRequired.Error(), e.Status, e.Message)
}

// Nametag is part of the account information.
type Nametag struct {
	Mode          int64       `json:"mode"`
	Gradient      json.Number `json:"gradient,Number"`
	Emoji         string      `json:"emoji"`
	SelfieSticker json.Number `json:"selfie_sticker,Number"`
}

type friendResp struct {
	Status     string     `json:"status"`
	Friendship Friendship `json:"friendship_status"`
}

// Location stores media location information.
type Location struct {
	insta *Instagram

	ID               int64   `json:"pk"`
	Name             string  `json:"name"`
	Address          string  `json:"address"`
	City             string  `json:"city"`
	ShortName        string  `json:"short_name"`
	Lng              float64 `json:"lng"`
	Lat              float64 `json:"lat"`
	ExternalSource   string  `json:"external_source"`
	FacebookPlacesID int64   `json:"facebook_places_id"`
}

// SuggestedUsers stores the information about user suggestions.
type SuggestedUsers struct {
	Type        int `json:"type"`
	Suggestions []struct {
		User struct {
			ID                         interface{}   `json:"pk"`
			Username                   string        `json:"username"`
			FullName                   string        `json:"full_name"`
			IsVerified                 bool          `json:"is_verified"`
			IsPrivate                  bool          `json:"is_private"`
			HasHighlightReels          bool          `json:"has_highlight_reels"`
			HasAnonymousProfilePicture bool          `json:"has_anonymous_profile_picture"`
			ProfilePicID               string        `json:"profile_pic_id"`
			ProfilePicURL              string        `json:"profile_pic_url"`
			AccountBadges              []interface{} `json:"account_badges"`
		} `json:"user"`
		Algorithm       string        `json:"algorithm"`
		SocialContext   string        `json:"social_context"`
		Icon            string        `json:"icon"`
		Caption         string        `json:"caption"`
		MediaIds        []interface{} `json:"media_ids"`
		ThumbnailUrls   []string      `json:"thumbnail_urls"`
		LargeUrls       []string      `json:"large_urls"`
		MediaInfos      []interface{} `json:"media_infos"`
		Value           float64       `json:"value"`
		IsNewSuggestion bool          `json:"is_new_suggestion"`
	} `json:"suggestions"`
	LandingSiteType  string `json:"landing_site_type"`
	Title            string `json:"title"`
	ViewAllText      string `json:"view_all_text"`
	LandingSiteTitle string `json:"landing_site_title"`
	NetegoType       string `json:"netego_type"`
	UpsellFbPos      string `json:"upsell_fb_pos"`
	AutoDvance       string `json:"auto_dvance"`
	ID               string `json:"id"`
	TrackingToken    string `json:"tracking_token"`
}
type PendingRequests struct {
	Users []*User `json:"users"`
	// TODO: pagination
	BigList                      bool           `json:"big_list"`
	GlobalBlacklistSample        interface{}    `json:"global_blacklist_sample"`
	NextMaxID                    string         `json:"next_max_id"`
	PageSize                     int            `json:"page_size"`
	TruncateFollowRequestAtIndex int            `json:"truncate_follow_requests_at_index"`
	Sections                     interface{}    `json:"sections"`
	SuggestedUsers               SuggestedUsers `json:"suggested_users"`
	Status                       string         `json:"status"`
}

// Friendship stores the details of the relationship between two users.
type Friendship struct {
	Following       bool `json:"following"`
	FollowedBy      bool `json:"followed_by"`
	IncomingRequest bool `json:"incoming_request"`
	OutgoingRequest bool `json:"outgoing_request"`
	Muting          bool `json:"muting"`
	Blocking        bool `json:"blocking"`
	IsBestie        bool `json:"is_bestie"`
	IsBlockingReel  bool `json:"is_blocking_reel"`
	IsMutingReel    bool `json:"is_muting_reel"`
	IsPrivate       bool `json:"is_private"`
	IsRestricted    bool `jsoN:"is_restricted"`
}

// Images are different quality images
type Images struct {
	Versions []Candidate `json:"candidates"`
}

// GetBest returns the URL of the image with the best quality.
func (img Images) GetBest() string {
	best := ""
	var mh, mw int
	for _, v := range img.Versions {
		if v.Width > mw || v.Height > mh {
			best = v.URL
			mh, mw = v.Height, v.Width
		}
	}
	return best
}

// Candidate is something that I really have no idea what it is.
type Candidate struct {
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	URL          string `json:"url"`
	ScansProfile string `json:"scans_profile"`
}

// Tag is the information of an user being tagged on any media.
type Tag struct {
	In []struct {
		User                  User        `json:"user"`
		Position              []float64   `json:"position"`
		StartTimeInVideoInSec interface{} `json:"start_time_in_video_in_sec"`
		DurationInVideoInSec  interface{} `json:"duration_in_video_in_sec"`
	} `json:"in"`
}

// Caption is media caption
type Caption struct {
	// can be both string or int
	ID              interface{} `json:"pk"`
	UserID          int64       `json:"user_id"`
	Text            string      `json:"text"`
	Type            int         `json:"type"`
	CreatedAt       int64       `json:"created_at"`
	CreatedAtUtc    int64       `json:"created_at_utc"`
	ContentType     string      `json:"content_type"`
	Status          string      `json:"status"`
	BitFlags        int         `json:"bit_flags"`
	User            User        `json:"user"`
	DidReportAsSpam bool        `json:"did_report_as_spam"`
	MediaID         int64       `json:"media_id"`
	HasTranslation  bool        `json:"has_translation"`
}

// Mentions is a user being mentioned on media.
type Mentions struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        int64   `json:"z"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Rotation float64 `json:"rotation"`
	IsPinned int     `json:"is_pinned"`
	User     User    `json:"user"`
}

// Video are different quality videos
type Video struct {
	Type   int    `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
	ID     string `json:"id"`
}

type timeStoryResp struct {
	Status string        `json:"status"`
	Media  []*StoryMedia `json:"tray"`
}

type trayResp struct {
	Reels  map[string]Reel `json:"reels"`
	Status string          `json:"status"`
}

// Tray is a set of story media received from timeline calls.
type Tray struct {
	Stories []*Reel `json:"tray"`
	// think this is depricated, and only broadcasts are used
	Lives struct {
		LiveItems []*LiveItems `json:"post_live_items"`
	} `json:"post_live"`
	StoryRankingToken    string       `json:"story_ranking_token"`
	Broadcasts           []*Broadcast `json:"broadcasts"`
	FaceFilterNuxVersion int          `json:"face_filter_nux_version"`
	HasNewNuxStory       bool         `json:"has_new_nux_story"`
	NuxElegible          bool         `json:"stories_viewer_gestures_nux_eligible"`
	StickerVersion       float64      `json:"sticker_version"`
	ReponseTS            float64      `json:"response_timestamp"`
	Status               string       `json:"status"`
	EmojiReactionsConfig struct {
		UfiType                        float64 `json:"ufi_type"`
		DeliveryType                   float64 `json:"delivery_type"`
		OverlaySkinTonePickerEnabled   bool    `json:"overlay_skin_tone_picker_enabled"`
		SwipeUpToShowReactions         bool    `json:"swipe_up_to_show_reactions"`
		ComposerNuxType                float64 `json:"composer_nux_type"`
		HideStoryViewCount             bool    `json:"hide_story_view_count"`
		ReactionTrayInteractivePanning bool    `json:"reaction_tray_interactive_panning_enabled"`
		PersistentSelfStoryBadge       bool    `json:"persistent_self_story_badge_enabled"`
		SelfstoryBadging               bool    `json:"self_story_badging_enabled"`
		ExitTestNux                    bool    `json:"exit_test_nux_enabled"`
	} `json:"emoji_reactions_config"`
}

func (tray *Tray) set(insta *Instagram) {
	for i := range tray.Stories {
		tray.Stories[i].setValues(insta)
	}
	for i := range tray.Lives.LiveItems {
		tray.Lives.LiveItems[i].User.insta = insta
		for j := range tray.Lives.LiveItems[i].Broadcasts {
			tray.Lives.LiveItems[i].Broadcasts[j].User.insta = insta
		}
	}
	for i := range tray.Broadcasts {
		tray.Broadcasts[i].setValues(insta)
	}
}

// LiveItems are Live media items
type LiveItems struct {
	ID                  string       `json:"pk"`
	User                User         `json:"user"`
	Broadcasts          []*Broadcast `json:"broadcasts"`
	LastSeenBroadcastTs float64      `json:"last_seen_broadcast_ts"`
	RankedPosition      int64        `json:"ranked_position"`
	SeenRankedPosition  int64        `json:"seen_ranked_position"`
	Muted               bool         `json:"muted"`
	CanReply            bool         `json:"can_reply"`
	CanReshare          bool         `json:"can_reshare"`
}

// BlockedUser stores information about a used that has been blocked before.
type BlockedUser struct {
	// TODO: Convert to user
	UserID        int64  `json:"user_id"`
	Username      string `json:"username"`
	FullName      string `json:"full_name"`
	ProfilePicURL string `json:"profile_pic_url"`
	BlockAt       int64  `json:"block_at"`
}

// Unblock unblocks blocked user.
func (b *BlockedUser) Unblock() error {
	u := User{ID: b.UserID}
	return u.Unblock()
}

type blockedListResp struct {
	BlockedList []BlockedUser `json:"blocked_list"`
	PageSize    int           `json:"page_size"`
	Status      string        `json:"status"`
}

// InboxItemMedia is inbox media item
type InboxItemMedia struct {
	ClientContext              string `json:"client_context"`
	ExpiringMediaActionSummary struct {
		Count     int    `json:"count"`
		Timestamp int64  `json:"timestamp"`
		Type      string `json:"type"`
	} `json:"expiring_media_action_summary"`
	ItemID     string `json:"item_id"`
	ItemType   string `json:"item_type"`
	RavenMedia struct {
		MediaType int64 `json:"media_type"`
	} `json:"raven_media"`
	ReplyChainCount int           `json:"reply_chain_count"`
	SeenUserIds     []interface{} `json:"seen_user_ids"`
	Timestamp       int64         `json:"timestamp"`
	UserID          int64         `json:"user_id"`
	ViewMode        string        `json:"view_mode"`
}

// InboxItemLike is the heart sent during a conversation.
type InboxItemLike struct {
	ItemID    string `json:"item_id"`
	ItemType  string `json:"item_type"`
	Timestamp int64  `json:"timestamp"`
	UserID    int64  `json:"user_id"`
}

type respLikers struct {
	Users     []*User `json:"users"`
	UserCount int64   `json:"user_count"`
	Status    string  `json:"status"`
}

type ErrChallengeProcess struct {
	StepName string
}

func (ec ErrChallengeProcess) Error() string {
	return ec.StepName
}

type Cooldowns struct {
	Default int    `json:"default"`
	Global  int    `json:"global"`
	Status  string `json:"status"`
	TTL     int    `json:"ttl"`
	Slots   []struct {
		Cooldown int    `json:"cooldown"`
		Slot     string `json:"slot"`
	} `json:"slots"`
	Surfaces []struct {
		Cooldown int    `json:"cooldown"`
		Slot     string `json:"slot"`
	} `json:"surfaces"`
}

type ScoresBootstrapUsers struct {
	Status   string `json:"status"`
	Surfaces []struct {
		Name      string         `json:"name"`
		RankToken string         `json:"rank_token"`
		Scores    map[string]int `json:"scores"`
		TTLSecs   int            `json:"ttl_secs"`
	} `json:"surfaces"`
	Users []*User `json:"users"`
}

type CommentOffensive struct {
	BullyClassifier  float64 `json:"bully_classifier"`
	SexualClassifier float64 `json:"sexual_classifier"`
	HateClassifier   float64 `json:"hate_classifier"`
	IsOffensive      bool    `json:"is_offensive"`
	Status           string  `json:"status"`
	TextLanguage     string  `json:"text_language"`
}

// Two factor authentication seed, used to generte the one time passwords
type TOTP struct {
	ID   int64  `json:"totp_seed_id"`
	Seed string `json:"totp_seed"`
}

type CommentInfo struct {
	LikesEnabled                   bool          `json:"comment_likes_enabled"`
	ThreadingEnabled               bool          `json:"comment_threading_enabled"`
	HasMore                        bool          `json:"has_more_comments"`
	MaxNumVisiblePreview           int           `json:"max_num_visible_preview_comments"`
	PreviewComments                []interface{} `json:"preview_comments"`
	CanViewMorePreview             bool          `json:"can_view_more_preview_comments"`
	CommentCount                   int           `json:"comment_count"`
	HideViewAllCommentEntrypoint   bool          `json:"hide_view_all_comment_entrypoint"`
	InlineComposerDisplayCondition string        `json:"inline_composer_display_condition"`
	InlineComposerImpTriggerTime   int           `json:"inline_composer_imp_trigger_time"`
}
