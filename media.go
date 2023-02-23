package goinsta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	neturl "net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Media interface defines methods for both StoryMedia and FeedMedia.
type Media interface {
	// Next allows pagination
	Next(...interface{}) bool
	// Error returns error (in case it have been occurred)
	Error() error
	// ID returns media id
	GetNextID() string
	// Delete removes media
	Delete() error
	getInsta() *Instagram
}

// Item represents media items
//
// All Item has Images or Videos objects which contains the url(s).
// You can use the Download function to get the best quality Image or Video from Item.
type Item struct {
	insta    *Instagram
	media    Media
	module   string
	Comments *Comments `json:"-"`

	// Post Info
	TakenAt           int64       `json:"taken_at"`
	Pk                int64       `json:"pk"`
	ID                interface{} `json:"id"` // Most of the times a string
	Index             int         // position in feed
	CommentsDisabled  bool        `json:"comments_disabled"`
	DeviceTimestamp   int64       `json:"device_timestamp"`
	FacepileTopLikers []struct {
		FollowFrictionType float64 `json:"follow_friction_type"`
		FullNeme           string  `json:"ful_name"`
		IsPrivate          bool    `json:"is_private"`
		IsVerified         bool    `json:"is_verified"`
		Pk                 float64 `json:"pk"`
		ProfilePicID       string  `json:"profile_pic_id"`
		ProfilePicURL      string  `json:"profile_pic_url"`
		Username           string  `json:"username"`
	} `json:"facepile_top_likers"`
	MediaType             int     `json:"media_type"`
	Code                  string  `json:"code"`
	ClientCacheKey        string  `json:"client_cache_key"`
	FilterType            int     `json:"filter_type"`
	User                  User    `json:"user"`
	CanReply              bool    `json:"can_reply"`
	CanReshare            bool    `json:"can_reshare"` // Used for stories
	CanViewerReshare      bool    `json:"can_viewer_reshare"`
	Caption               Caption `json:"caption"`
	CaptionIsEdited       bool    `json:"caption_is_edited"`
	LikeViewCountDisabled bool    `json:"like_and_view_counts_disabled"`
	FundraiserTag         struct {
		HasStandaloneFundraiser bool `json:"has_standalone_fundraiser"`
	} `json:"fundraiser_tag"`
	IsSeen                       bool   `json:"is_seen"`
	InventorySource              string `json:"inventory_source"`
	ProductType                  string `json:"product_type"`
	Likes                        int    `json:"like_count"`
	HasLiked                     bool   `json:"has_liked"`
	NearlyCompleteCopyRightMatch bool   `json:"nearly_complete_copyright_match"`
	// Toplikers can be `string` or `[]string`.
	// Use TopLikers function instead of getting it directly.
	Toplikers  interface{} `json:"top_likers"`
	Likers     []*User     `json:"likers"`
	PhotoOfYou bool        `json:"photo_of_you"`

	// Comments
	CommentLikesEnabled          bool `json:"comment_likes_enabled"`
	CommentThreadingEnabled      bool `json:"comment_threading_enabled"`
	HasMoreComments              bool `json:"has_more_comments"`
	MaxNumVisiblePreviewComments int  `json:"max_num_visible_preview_comments"`

	// To fetch, call feed.GetCommentInfo(), or item.GetCommentInfo()
	CommentInfo *CommentInfo

	// Will always be zero, call feed.GetCommentInfo()
	CommentCount int `json:"comment_count"`

	// Previewcomments can be `string` or `[]string` or `[]Comment`.
	// Use PreviewComments function instead of getting it directly.
	Previewcomments interface{} `json:"preview_comments,omitempty"`

	// Tags are tagged people in photo
	Tags struct {
		In []Tag `json:"in"`
	} `json:"usertags,omitempty"`
	FbUserTags           Tag    `json:"fb_user_tags"`
	CanViewerSave        bool   `json:"can_viewer_save"`
	OrganicTrackingToken string `json:"organic_tracking_token"`
	// Images contains URL images in different versions.
	// Version = quality.
	Images          Images   `json:"image_versions2,omitempty"`
	OriginalWidth   int      `json:"original_width,omitempty"`
	OriginalHeight  int      `json:"original_height,omitempty"`
	ImportedTakenAt int64    `json:"imported_taken_at,omitempty"`
	Location        Location `json:"location,omitempty"`
	Lat             float64  `json:"lat,omitempty"`
	Lng             float64  `json:"lng,omitempty"`

	// Carousel
	CarouselParentID string `json:"carousel_parent_id"`
	CarouselMedia    []Item `json:"carousel_media,omitempty"`

	// Live
	IsPostLive bool `json:"is_post_live"`

	// Videos
	Videos            []Video `json:"video_versions,omitempty"`
	VideoCodec        string  `json:"video_codec"`
	HasAudio          bool    `json:"has_audio,omitempty"`
	VideoDuration     float64 `json:"video_duration,omitempty"`
	ViewCount         float64 `json:"view_count,omitempty"`
	PlayCount         float64 `json:"play_count,omitempty"`
	IsDashEligible    int     `json:"is_dash_eligible,omitempty"`
	IsUnifiedVideo    bool    `json:"is_unified_video"`
	VideoDashManifest string  `json:"video_dash_manifest,omitempty"`
	NumberOfQualities int     `json:"number_of_qualities,omitempty"`

	// IGTV
	Title                    string `json:"title"`
	IGTVExistsInViewerSeries bool   `json:"igtv_exists_in_viewer_series"`
	IGTVSeriesInfo           struct {
		HasCoverPhoto bool `json:"has_cover_photo"`
		ID            int64
		NumEpisodes   int    `json:"num_episodes"`
		Title         string `json:"title"`
	} `json:"igtv_series_info"`
	IGTVAdsInfo struct {
		AdsToggledOn            bool `json:"ads_toggled_on"`
		ElegibleForInsertingAds bool `json:"is_video_elegible_for_inserting_ads"`
	} `json:"igtv_ads_info"`

	// Ads
	IsCommercial        bool   `json:"is_commercial"`
	IsPaidPartnership   bool   `json:"is_paid_partnership"`
	CommercialityStatus string `json:"commerciality_status"`
	AdLink              string `json:"link"`
	AdLinkText          string `json:"link_text"`
	AdLinkHint          string `json:"link_hint_text"`
	AdTitle             string `json:"overlay_title"`
	AdSubtitle          string `json:"overlay_subtitle"`
	AdText              string `json:"overlay_text"`
	AdAction            string `json:"ad_action"`
	AdHeaderStyle       int    `json:"ad_header_style"`
	AdLinkType          int    `json:"ad_link_type"`
	AdMetadata          []struct {
		Type  int         `json:"type"`
		Value interface{} `json:"value"`
	} `json:"ad_metadata"`
	AndroidLinks []struct {
		AndroidClass      string `json:"androidClass"`
		CallToActionTitle string `json:"callToActionTitle"`
		DeeplinkUri       string `json:"deeplinkUri"`
		LinkType          int    `json:"linkType"`
		Package           string `json:"package"`
		WebUri            string `json:"webUri"`
	} `json:"android_links"`

	// Only for stories
	StoryEvents              []interface{}      `json:"story_events"`
	StoryHashtags            []interface{}      `json:"story_hashtags"`
	StoryPolls               []interface{}      `json:"story_polls"`
	StoryFeedMedia           []interface{}      `json:"story_feed_media"`
	StorySoundOn             []interface{}      `json:"story_sound_on"`
	CreativeConfig           interface{}        `json:"creative_config"`
	StoryLocations           []interface{}      `json:"story_locations"`
	StorySliders             []interface{}      `json:"story_sliders"`
	StoryQuestions           []interface{}      `json:"story_questions"`
	StoryProductItems        []interface{}      `json:"story_product_items"`
	StoryCTA                 []StoryCTA         `json:"story_cta"`
	IntegrityReviewDecision  string             `json:"integrity_review_decision"`
	IsReelMedia              bool               `json:"is_reel_media"`
	ProfileGridControl       bool               `json:"profile_grid_control_enabled"`
	ReelMentions             []StoryReelMention `json:"reel_mentions"`
	ExpiringAt               int64              `json:"expiring_at"`
	CanSendCustomEmojis      bool               `json:"can_send_custom_emojis"`
	SupportsReelReactions    bool               `json:"supports_reel_reactions"`
	ShowOneTapFbShareTooltip bool               `json:"show_one_tap_fb_share_tooltip"`
	HasSharedToFb            int64              `json:"has_shared_to_fb"`
	Mentions                 []Mentions
	Audience                 string `json:"audience,omitempty"`
	StoryMusicStickers       []struct {
		X              float64 `json:"x"`
		Y              float64 `json:"y"`
		Z              int     `json:"z"`
		Width          float64 `json:"width"`
		Height         float64 `json:"height"`
		Rotation       float64 `json:"rotation"`
		IsPinned       int     `json:"is_pinned"`
		IsHidden       int     `json:"is_hidden"`
		IsSticker      int     `json:"is_sticker"`
		MusicAssetInfo struct {
			ID                       string `json:"id"`
			Title                    string `json:"title"`
			Subtitle                 string `json:"subtitle"`
			DisplayArtist            string `json:"display_artist"`
			CoverArtworkURI          string `json:"cover_artwork_uri"`
			CoverArtworkThumbnailURI string `json:"cover_artwork_thumbnail_uri"`
			ProgressiveDownloadURL   string `json:"progressive_download_url"`
			HighlightStartTimesInMs  []int  `json:"highlight_start_times_in_ms"`
			IsExplicit               bool   `json:"is_explicit"`
			DashManifest             string `json:"dash_manifest"`
			HasLyrics                bool   `json:"has_lyrics"`
			AudioAssetID             string `json:"audio_asset_id"`
			IgArtist                 struct {
				Pk            int64  `json:"pk"`
				Username      string `json:"username"`
				FullName      string `json:"full_name"`
				IsPrivate     bool   `json:"is_private"`
				ProfilePicURL string `json:"profile_pic_url"`
				ProfilePicID  string `json:"profile_pic_id"`
				IsVerified    bool   `json:"is_verified"`
			} `json:"ig_artist"`
			PlaceholderProfilePicURL string `json:"placeholder_profile_pic_url"`
			ShouldMuteAudio          bool   `json:"should_mute_audio"`
			ShouldMuteAudioReason    string `json:"should_mute_audio_reason"`
			OverlapDurationInMs      int    `json:"overlap_duration_in_ms"`
			AudioAssetStartTimeInMs  int    `json:"audio_asset_start_time_in_ms"`
			StoryLinkStickers        []struct {
				X           float64 `json:"x"`
				Y           float64 `json:"y"`
				Z           int     `json:"z"`
				Width       float64 `json:"width"`
				Height      float64 `json:"height"`
				Rotation    int     `json:"rotation"`
				IsPinned    int     `json:"is_pinned"`
				IsHidden    int     `json:"is_hidden"`
				IsSticker   int     `json:"is_sticker"`
				IsFbSticker int     `json:"is_fb_sticker"`
				StoryLink   struct {
					LinkType   string `json:"link_type"`
					URL        string `json:"url"`
					LinkTitle  string `json:"link_title"`
					DisplayURL string `json:"display_url"`
				} `json:"story_link"`
			} `json:"story_link_stickers"`
		} `json:"music_asset_info"`
	} `json:"story_music_stickers,omitempty"`
}

// Comment pushes a text comment to media item.
//
// If parent media is a Story this function will send a private message
// replying the Instagram story.
func (item *Item) Comment(text string) error {
	// if comment is called on a story media, use reply
	if item.ProductType == "story" {
		return item.Reply(text)
	}

	o, err := item.CommentCheckOffensive(text)
	if err != nil {
		return err
	}
	if !o.IsOffensive {
		return item.comment(text)
	}
	return errors.New("Failed to post comment, flagged as offensive")
}

func (item *Item) comment(text string) error {
	insta := item.insta
	query := map[string]string{
		// "feed_position":           "",
		"container_module":        "feed_timeline",
		"user_breadcrumb":         generateUserBreadcrumb(text),
		"nav_chain":               "",
		"_uid":                    toString(insta.Account.ID),
		"_uuid":                   insta.uuid,
		"idempotence_token":       generateUUID(),
		"radio_type":              "wifi-none",
		"is_carousel_bumped_post": "false", // not sure when this would be true
		"comment_text":            text,
	}
	if item.module != "" {
		query["container_module"] = item.module
	}
	if item.IsCommercial {
		query["delivery_class"] = "ad"
	} else {
		query["delivery_class"] = "organic"
	}
	if item.InventorySource != "" {
		query["inventory_source"] = item.InventorySource
	}
	if len(item.CarouselMedia) > 0 || item.CarouselParentID != "" {
		query["carousel_index"] = "0"
	}
	b, err := json.Marshal(query)
	if err != nil {
		return err
	}

	// ignoring response
	_, _, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlCommentAdd, item.Pk),
			Query:    map[string]string{"signed_body": "SIGNATURE." + string(b)},
			IsPost:   true,
		},
	)
	return err
}

func (item *Item) CommentCheckOffensive(comment string) (*CommentOffensive, error) {
	insta := item.insta
	data, err := json.Marshal(map[string]string{
		"media_id":           item.GetID(),
		"_uid":               toString(insta.Account.ID),
		"comment_session_id": generateUUID(),
		"_uuid":              insta.uuid,
		"comment_text":       comment,
	})
	if err != nil {
		return nil, err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlCommentOffensive,
			IsPost:   true,
			Query:    generateSignature(data),
		},
	)
	if err != nil {
		return nil, err
	}
	r := &CommentOffensive{}
	err = json.Unmarshal(body, r)
	return r, err
}

func (item *Item) Reply(text string) error {
	if item.ProductType != "story" {
		return item.Comment(text)
	}

	insta := item.insta

	to, err := prepareRecipients(item)
	if err != nil {
		return err
	}

	token := "68" + randNum(17)
	_, _, err = insta.sendRequest(
		&reqOptions{
			Connection: "keep-alive",
			Endpoint:   fmt.Sprintf("%s?media_type=%s", urlReplyStory, item.MediaToString()),
			IsPost:     true,
			Query: map[string]string{
				"recipient_users":      to,
				"action":               "send_item",
				"is_shh_mode":          "0", // not sure when this would be 1/true
				"send_attribution":     "reel",
				"client_context":       token,
				"media_id":             item.GetID(),
				"text":                 text,
				"device_id":            insta.dID,
				"mutation_token":       token,
				"_uuid":                insta.uuid,
				"entry":                "reel",
				"reel_id":              toString(item.User.ID),
				"offline_threading_id": token,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// MediaToString returns Item.MediaType as string.
func (item *Item) MediaToString() string {
	return MediaToString(item.MediaType)
}

func MediaToString(t int) string {
	switch t {
	case 1:
		return "photo"
	case 2:
		return "video"
	case 6:
		return "ad_map"
	case 7:
		return "live"
	case 8:
		return "carousel"
	case 9:
		return "live_replay"
	case 10:
		return "collection"
	case 11:
		return "audio"
	case 12:
		return "showreel_native"
	case 13:
		return "guide_facade"
	}
	return ""
}

// TODO: remove excessive insta = media.insta lines from setValues() funcs
func setToItem(item *Item, media Media) {
	item.media = media
	item.insta = media.getInsta()
	item.User.insta = media.getInsta()
	item.Comments = newComments(item)
	for i := range item.CarouselMedia {
		item.CarouselMedia[i].User = item.User
		setToItem(&item.CarouselMedia[i], media)
	}
}

// setToMediaItem is a utility function that
// mimics the setToItem but for the SavedMedia items
func setToMediaItem(item *MediaItem, media Media) {
	item.Media.media = media
	item.Media.Comments = newComments(&item.Media)

	for i := range item.Media.CarouselMedia {
		item.Media.CarouselMedia[i].User = item.Media.User
		setToItem(&item.Media.CarouselMedia[i], media)
	}
}

func getname(name string) string {
	nname := name
	i := 1
	for {
		ext := path.Ext(name)

		_, err := os.Stat(name)
		if err != nil {
			break
		}
		if ext != "" {
			nname = strings.Replace(nname, ext, "", -1)
		}
		name = fmt.Sprintf("%s.%d%s", nname, i, ext)
		i++
	}
	return name
}

type bestMedia struct {
	w, h int
	url  string
}

// GetBest returns url to best quality image or video.
//
// Arguments can be []Video or []Candidate
func GetBest(obj interface{}) string {
	m := bestMedia{}

	switch t := obj.(type) {
	// getting best video
	case []Video:
		for _, video := range t {
			if m.w < video.Width && video.Height > m.h && video.URL != "" {
				m.w = video.Width
				m.h = video.Height
				m.url = video.URL
			}
		}
		// getting best image
	case []Candidate:
		for _, image := range t {
			if m.w < image.Width && image.Height > m.h && image.URL != "" {
				m.w = image.Width
				m.h = image.Height
				m.url = image.URL
			}
		}
	}
	return m.url
}

var rxpTags = regexp.MustCompile(`#\w+`)

// Hashtags returns caption hashtags.
//
// Item media parent must be FeedMedia.
//
// See example: examples/media/hashtags.go
func (item *Item) Hashtags() []Hashtag {
	tags := rxpTags.FindAllString(item.Caption.Text, -1)

	hsh := make([]Hashtag, len(tags))

	i := 0
	for _, tag := range tags {
		hsh[i].Name = tag[1:]
		i++
	}

	for _, comment := range item.PreviewComments() {
		tags := rxpTags.FindAllString(comment.Text, -1)

		for _, tag := range tags {
			hsh = append(hsh, Hashtag{Name: tag[1:]})
		}
	}

	return hsh
}

// Delete deletes your media item. StoryMedia or FeedMedia
//
// See example: examples/media/mediaDelete.go
func (item *Item) Delete() error {
	return item.insta.delete(toString(item.ID), item.MediaToString(), item)
}

func (insta *Instagram) delete(id, media string, mediaType interface{}) error {
	query := map[string]string{
		"media_id": id,
		"_uid":     toString(insta.Account.ID),
		"_uuid":    insta.uuid,
	}

	switch mediaType.(type) {
	case *StoryMedia:
		query["deep_delete_waterfall"] = generateUUID()
	case *FeedMedia:
		query["igtv_feed_preview"] = "false"
	case *Timeline:
		query["igtv_feed_preview"] = "false"
	}

	data, err := json.Marshal(query)
	if err != nil {
		return err
	}

	mediaParam := "?media_type=" + strings.ToUpper(media)

	_, _, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlMediaDelete, id) + mediaParam,
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// SyncLikers fetch new likers of a media
//
// This function updates Item.Likers value
func (item *Item) SyncLikers() error {
	resp := respLikers{}
	insta := item.insta
	body, err := insta.sendSimpleRequest(urlMediaLikers, item.ID)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &resp)
	if err == nil {
		item.Likers = resp.Users
	}

	for _, u := range item.Likers {
		u.insta = insta
	}
	return err
}

// Like mark media item as liked.
//
// See example: examples/media/like.go
func (item *Item) Like() error {
	if !item.HasLiked {
		return item.changeLike(urlMediaLike)
	}
	return nil
}

// Unlike mark media item as unliked.
//
// See example: examples/media/unlike.go
func (item *Item) Unlike() error {
	if item.HasLiked {
		return item.changeLike(urlMediaUnlike)
	}
	return nil
}

func (item *Item) changeLike(endpoint string) error {
	insta := item.insta
	query := map[string]string{
		// "feed_position":           "",
		"container_module":        "feed_timelin",
		"nav_chain":               "",
		"_uid":                    toString(insta.Account.ID),
		"_uuid":                   insta.uuid,
		"radio_type":              "wifi-none",
		"is_carousel_bumped_post": "false", // not sure when this would be true
	}
	if item.module != "" {
		query["container_module"] = item.module
	}
	if item.IsCommercial {
		query["delivery_class"] = "ad"
	} else {
		query["delivery_class"] = "organic"
	}
	if item.InventorySource != "" {
		query["inventory_source"] = item.InventorySource
	}
	if len(item.CarouselMedia) > 0 || item.CarouselParentID != "" {
		query["carousel_index"] = "0"
	}

	data, err := json.Marshal(query)
	if err != nil {
		return err
	}

	_, _, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(endpoint, item.ID),
			Query: map[string]string{
				"signed_body": "SIGNATURE." + string(data),
				"d":           "0",
			},
			IsPost: true,
		},
	)
	return err
}

// Download downloads media item (video or image) with the best quality.
//
// Input parameter is a path to either a directory or a file. If no file is
//   specified it will try to extract a file name from the image and use that.
// If file extentions will automatically be appended to the file name, and do
//   not need to be set manually.
//
// If file exists it will be saved
// This function makes folder automatically
//
// See example: examples/media/itemDownload.go

// func (item *Item) Download(folder, name string) (m []byte, err error) {
// 	return nil, nil
// }
func (item *Item) DownloadTo(dst string) error {
	insta := item.insta
	folder, file := path.Split(dst)

	if err := os.MkdirAll(folder, 0o777); err != nil {
		return err
	}

	switch item.MediaType {
	case 1:
		return insta.downloadTo(folder, file, item.Images.Versions)
	case 2:
		return insta.downloadTo(folder, file, item.Videos)
	case 8:
		return item.downloadCarousel(folder, file)
	}

	insta.warnHandler(
		fmt.Sprintf(
			"Unable to download %s media (media type %d), this has not been implemented",
			item.MediaToString(),
			item.MediaType,
		),
	)
	return ErrNoMedia

}

// Download will download a media item and directly return it as a byte slice.
// If you wish to download a picture to a folder, use item.DownloadTo(path)
func (item *Item) Download() ([]byte, error) {
	insta := item.insta

	switch item.MediaType {
	case 1:
		url := GetBest(item.Images.Versions)
		return insta.download(url)
	case 2:
		url := GetBest(item.Videos)
		return insta.download(url)
	case 8:
		return nil, fmt.Errorf("Unable to download a carousel with this method, use DownloadTo instead to save it to a file. If this is a feature you wish to use please let me know.")
	}
	return nil, ErrNoMedia
}

func (item *Item) downloadCarousel(folder, fn string) error {
	if fn == "" {
		fn = item.GetID()
	}
	for i, media := range item.CarouselMedia {
		n := fmt.Sprintf("%s_%d", fn, i+1)
		if err := media.DownloadTo(path.Join(folder, n)); err != nil {
			return err
		}
	}
	return nil
}

// downloadTo saves a media item to folder/file
func (insta *Instagram) downloadTo(folder, fn string, media interface{}) error {
	url := GetBest(media)
	fn, err := getDownloadName(url, fn)
	if err != nil {
		return err
	}
	b, err := insta.download(url)
	if err != nil {
		return err
	}
	err = saveToFolder(folder, fn, b)
	return err
}

// download the media from a url and return the bytes
func (insta *Instagram) download(url string) ([]byte, error) {
	resp, err := insta.c.Get(url)
	if err != nil {
		return nil, err
	}
	media, err := io.ReadAll(resp.Body)
	return media, err
}

// saveToFolder writes bytes to a file
func saveToFolder(folder, fn string, media []byte) error {
	dst := path.Join(folder, fn)

	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()

	b := bytes.NewBuffer(media)
	_, err = io.Copy(file, b)
	return err

}

func getDownloadName(url, name string) (string, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}
	ext := path.Ext(u.Path)
	if name == "" {
		name = path.Base(u.Path)
	} else if !strings.HasSuffix(name, ext) {
		name += ext
	}
	name = getname(name)
	return name, nil
}

// TopLikers returns string slice or single string (inside string slice)
// Depending on TopLikers parameter.
func (item *Item) TopLikers() []string {
	switch s := item.Toplikers.(type) {
	case string:
		return []string{s}
	case []string:
		return s
	}
	return nil
}

// PreviewComments returns string slice or single string (inside Comment slice)
// Depending on PreviewComments parameter.
// If PreviewComments are string or []string only the Text field will be filled.
func (item *Item) PreviewComments() []Comment {
	switch s := item.Previewcomments.(type) {
	case []interface{}:
		if len(s) == 0 {
			return nil
		}

		switch s[0].(type) {
		case string:
			comments := make([]Comment, 0)
			for i := range s {
				comments = append(comments, Comment{
					Text: s[i].(string),
				})
			}
			return comments
		case interface{}:
			comments := make([]Comment, 0)
			for i := range s {
				if buf, err := json.Marshal(s[i]); err != nil {
					return nil
				} else {
					comment := &Comment{}

					if err = json.Unmarshal(buf, comment); err != nil {
						return nil
					} else {
						comments = append(comments, *comment)
					}
				}
			}
			return comments
		}
	case string:
		comments := []Comment{
			{
				Text: s,
			},
		}
		return comments
	}
	return nil
}

// StoryIsCloseFriends returns a bool
// If the returned value is true the story was published only for close friends
func (item *Item) StoryIsCloseFriends() bool {
	return item.Audience == "besties"
}

func (item *Item) GetID() string {
	return toString(item.ID)
}

// FeedMedia represent a set of media items
// Mainly used for user profile feeds. To get your main timeline use insta.Timeline
type FeedMedia struct {
	insta *Instagram

	err error

	uid       int64
	endpoint  string
	timestamp string

	Items               []*Item `json:"items"`
	NumResults          int     `json:"num_results"`
	MoreAvailable       bool    `json:"more_available"`
	AutoLoadMoreEnabled bool    `json:"auto_load_more_enabled"`
	Status              string  `json:"status"`
	// Can be int64 and string
	// this is why we recommend Next() usage :')
	NextID interface{} `json:"next_max_id"`
}

// Delete deletes ALL items in media.
// If you want to delete one item, pick one from media.Items and call Item.Delete()
//
// See example: examples/media/mediaDelete.go
func (media *FeedMedia) Delete() error {
	for i := range media.Items {
		if err := media.Items[i].Delete(); err != nil {
			return errors.Wrap(err, "failed to delete item")
		}
	}
	return nil
}

// SetInstagram set instagram
func (media *FeedMedia) SetInstagram(insta *Instagram) {
	media.insta = insta
}

// SetID sets media.GetNextID
// this value can be int64 or string
func (media *FeedMedia) SetID(id interface{}) {
	media.NextID = id
}

// Sync updates media values.
func (media *FeedMedia) Sync() error {
	// haven't actually seen this being called, so not sure if its correct
	id := media.GetNextID()
	insta := media.insta

	data, err := json.Marshal(
		map[string]string{
			"media_id": id,
			"_uid":     toString(insta.Account.ID),
			"_uuid":    insta.uuid,
		},
	)
	if err != nil {
		return err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlMediaInfo, id),
			IsPost:   false,
			Query:    generateSignature(data),
		},
	)
	if err != nil {
		return err
	}

	m := FeedMedia{
		insta:    insta,
		endpoint: urlMediaInfo,
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return err
	}

	*media = m
	media.NextID = id
	media.setValues()
	return nil
}

func (media *FeedMedia) setIndex() {
	for i, v := range media.Items {
		v.Index = i
	}
}

func (media *FeedMedia) setValues() {
	for i := range media.Items {
		media.Items[i].insta = media.insta
		media.Items[i].User.insta = media.insta
		setToItem(media.Items[i], media)
	}
}

func (media *FeedMedia) Error() error {
	return media.err
}

func (media *FeedMedia) getInsta() *Instagram {
	return media.insta
}

// ID returns media id.
func (media *FeedMedia) GetNextID() string {
	return formatID(media.NextID)
}

// Next allows pagination after calling:
// User.Feed
// extra query arguments can be passes one after another as func(key, value).
// Only if an even number of string arguements will be passed, they will be
//   used in the query.
// returns false when list reach the end.
// if FeedMedia.Error() is ErrNoMore no problems have occurred.
func (media *FeedMedia) Next(params ...interface{}) bool {
	if media.err != nil {
		return false
	}

	insta := media.insta
	endpoint := media.endpoint
	if media.uid != 0 {
		endpoint = fmt.Sprintf(endpoint, media.uid)
	}

	query := map[string]string{
		"exclude_comment":                 "true",
		"only_fetch_first_carousel_media": "false",
	}
	if len(params)%2 == 0 {
		for i := 0; i < len(params); i = i + 2 {
			query[params[i].(string)] = params[i+1].(string)
		}
	}

	if next := media.GetNextID(); next != "" {
		query["max_id"] = next
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: endpoint,
			Query:    query,
		},
	)
	if err == nil {
		m := FeedMedia{
			insta: insta,
		}
		d := json.NewDecoder(bytes.NewReader(body))
		d.UseNumber()
		err = d.Decode(&m)
		if err == nil {
			media.NextID = m.NextID
			media.MoreAvailable = m.MoreAvailable
			media.NumResults = m.NumResults
			media.AutoLoadMoreEnabled = m.AutoLoadMoreEnabled
			media.Status = m.Status
			if m.NextID == 0 || !m.MoreAvailable {
				media.err = ErrNoMore
			}
			m.setValues()
			media.Items = append(media.Items, m.Items...)
			media.setIndex()
			return true
		}
	}
	return false
}

// GetCommentInfo will fetch the item.CommentInfo for an item
func (item *Item) GetCommentInfo() error {
	insta := item.insta

	query := map[string]string{
		"media_ids": item.GetID(),
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlMediaCommentInfos,
			Query:    query,
		},
	)
	if err != nil {
		return err
	}

	var rsp struct {
		CommentInfos map[string]*CommentInfo `json:"comment_infos"`
	}
	if err := json.Unmarshal(body, &rsp); err != nil {
		return err
	}
	item.CommentInfo = rsp.CommentInfos[item.GetID()]

	return nil
}

// GetCommentInfo will fetch the item.CommentInfo; e.g. comment counts, and
//  other comment information for the feed.Latest() items
func (media *FeedMedia) GetCommentInfo() error {
	insta := media.insta

	items := media.Latest()
	query := map[string]string{
		"media_ids": "",
	}

	for _, item := range items {
		query["media_ids"] = query["media_ids"] + "," + item.GetID()
	}
	query["media_ids"] = strings.TrimPrefix(query["media_ids"], ",")

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlMediaCommentInfos,
			Query:    query,
		},
	)
	if err != nil {
		return err
	}

	var rsp struct {
		CommentInfos map[string]*CommentInfo `json:"comment_infos"`
	}
	if err := json.Unmarshal(body, &rsp); err != nil {
		return err
	}
	for id, info := range rsp.CommentInfos {
		for _, item := range items {
			if item.GetID() == id {
				item.CommentInfo = info
				break
			}
		}
	}

	return nil
}

// Latest returns a slice of the latest fetched items of the list of all items.
// The Next method keeps adding to the list, with Latest you can retrieve only
// the newest items.
func (media *FeedMedia) Latest() []*Item {
	return media.Items[len(media.Items)-media.NumResults:]
}
