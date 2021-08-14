package goinsta

import "errors"

const (
	// urls
	baseUrl        = "https://i.instagram.com/"
	instaAPIUrl    = "https://i.instagram.com/api/v1/"
	instaAPIUrlb   = "https://b.i.instagram.com/api/v1/"
	instaAPIUrlv2  = "https://i.instagram.com/api/v2/"
	instaAPIUrlv2b = "https://b.i.instagram.com/api/v2/"

	// header values
	bloksVerID         = "927f06374b80864ae6a0b04757048065714dc50ff15d2b8b3de8d0b6de961649"
	fbAnalytics        = "567067343352427"
	igCapabilities     = "3brTvx0="
	connType           = "WIFI"
	instaSigKeyVersion = "4"
	locale             = "en_US"
	appVersion         = "195.0.0.31.123"
	appVersionCode     = "302733750"

	// Used for supported_capabilities value used in some requests, e.g. tray requests
	supportedSdkVersions = "100.0,101.0,102.0,103.0,104.0,105.0,106.0,107.0,108.0,109.0,110.0,111.0,112.0,113.0,114.0,115.0,116.0,117.0"
	facetrackerVersion   = "14"
	segmentation         = "segmentation_enabled"
	compression          = "ETC2_COMPRESSION"
	worldTracker         = "world_tracker_enabled"
	gyroscope            = "gyroscope_enabled"

	// Other
	software = "Android RP1A.200720.012.G975FXXSBFUF3"
	hmacKey  = "iN4$aGr0m"
)

var (
	defaultHeaderOptions = map[string]string{
		"X-Ig-Www-Claim": "0",
	}
	// Default Device
	GalaxyS10 = Device{
		Manufacturer:     "samsung",
		Model:            "SM-G975F",
		CodeName:         "beyond2",
		AndroidVersion:   30,
		AndroidRelease:   11,
		ScreenDpi:        "560dpi",
		ScreenResolution: "1440x2898",
		Chipset:          "exynos9820",
	}
	G6 = Device{
		Manufacturer:     "LGE/lge",
		Model:            "LG-H870DS",
		CodeName:         "lucye",
		AndroidVersion:   28,
		AndroidRelease:   9,
		ScreenDpi:        "560dpi",
		ScreenResolution: "1440x2698",
		Chipset:          "lucye",
	}
	timeOffset = getTimeOffset()
)

type muteOption string

const (
	MuteAll   muteOption = "all"
	MuteStory muteOption = "reel"
	MutePosts muteOption = "post"
)

// Endpoints (with format vars)
const (
	// Login
	urlMsisdnHeader               = "accounts/read_msisdn_header/"
	urlGetPrefill                 = "accounts/get_prefill_candidates/"
	urlContactPrefill             = "accounts/contact_point_prefill/"
	urlGetAccFamily               = "multiple_accounts/get_account_family/"
	urlZrToken                    = "zr/token/result/"
	urlLogin                      = "accounts/login/"
	urlLogout                     = "accounts/logout/"
	urlAutoComplete               = "friendships/autocomplete_user_list/"
	urlQeSync                     = "qe/sync/"
	urlSync                       = "launcher/sync/"
	urlLogAttribution             = "attribution/log_attribution/"
	urlMegaphoneLog               = "megaphone/log/"
	urlExpose                     = "qe/expose/"
	urlGetNdxSteps                = "devices/ndx/api/async_get_ndx_ig_steps/"
	urlBanyan                     = "banyan/banyan/"
	urlCooldowns                  = "qp/get_cooldowns/"
	urlFetchConfig                = "loom/fetch_config/"
	urlBootstrapUserScores        = "scores/bootstrap/users/"
	urlStoreClientPushPermissions = "notifications/store_client_push_permissions/"
	urlProcessContactPointSignals = "accounts/process_contact_point_signals/"

	// Account
	urlCurrentUser      = "accounts/current_user/"
	urlChangePass       = "accounts/change_password/"
	urlSetPrivate       = "accounts/set_private/"
	urlSetPublic        = "accounts/set_public/"
	urlRemoveProfPic    = "accounts/remove_profile_picture/"
	urlChangeProfPic    = "accounts/change_profile_picture/"
	urlFeedSaved        = "feed/saved/all/"
	urlFeedSavedPosts   = "feed/saved/posts/"
	urlFeedSavedIGTV    = "feed/saved/igtv/"
	urlEditProfile      = "accounts/edit_profile/"
	urlFeedLiked        = "feed/liked/"
	urlConsent          = "consent/existing_user_flow/"
	urlNotifBadge       = "notifications/badge/"
	urlFeaturedAccounts = "multiple_accounts/get_featured_accounts/"

	// Collections
	urlCollectionsList     = "collections/list/"
	urlCollectionsCreate   = "collections/create/"
	urlCollectionEdit      = "collections/%s/edit/"
	urlCollectionDelete    = "collections/%s/delete/"
	urlCollectionFeedAll   = "feed/collection/%s/all/"
	urlCollectionFeedPosts = "feed/collection/%s/posts/"

	// Account and profile
	urlFollowers = "friendships/%d/followers/"
	urlFollowing = "friendships/%d/following/"

	// Users
	urlUserArchived      = "feed/only_me_feed/"
	urlUserByName        = "users/%s/usernameinfo/"
	urlUserByID          = "users/%s/info/"
	urlUserBlock         = "friendships/block/%d/"
	urlUserUnblock       = "friendships/unblock/%d/"
	urlUserMute          = "friendships/mute_posts_or_story_from_follow/"
	urlUserUnmute        = "friendships/unmute_posts_or_story_from_follow/"
	urlUserFollow        = "friendships/create/%d/"
	urlUserUnfollow      = "friendships/destroy/%d/"
	urlUserFeed          = "feed/user/%d/"
	urlFriendship        = "friendships/show/%d/"
	urlFriendshipPending = "friendships/pending/"
	urlUserStories       = "feed/user/%d/story/"
	urlUserTags          = "usertags/%d/feed/"
	urlBlockedList       = "users/blocked_list/"
	urlUserInfo          = "users/%d/info/"
	urlUserHighlights    = "highlights/%d/highlights_tray/"

	// Timeline
	urlTimeline  = "feed/timeline/"
	urlStories   = "feed/reels_tray/"
	urlReelMedia = "feed/reels_media/"

	// Search
	urlSearchTop           = "fbsearch/topsearch_flat/"
	urlSearchUser          = "users/search/"
	urlSearchTag           = "tags/search/"
	urlSearchLocation      = "fbsearch/places/"
	urlSearchRecent        = "fbsearch/recent_searches/"
	urlSearchNullState     = "fbsearch/nullstate_dynamic_sections/"
	urlSearchRegisterClick = "fbsearch/register_recent_search_click/"

	// Feeds
	urlFeedLocationID = "feed/location/%d/"
	urlFeedLocations  = "locations/%d/sections/"
	urlFeedTag        = "feed/tag/%s/"

	// Media
	urlMediaInfo    = "media/%s/info/"
	urlMediaDelete  = "media/%s/delete/"
	urlMediaLike    = "media/%s/like/"
	urlMediaUnlike  = "media/%s/unlike/"
	urlMediaSave    = "media/%s/save/"
	urlMediaUnsave  = "media/%s/unsave/"
	urlMediaSeen    = "media/seen/"
	urlMediaLikers  = "media/%s/likers/"
	urlMediaBlocked = "media/blocked/"

	// Broadcasts
	urlLiveInfo      = "live/%s/info/"
	urlLiveComments  = "live/%s/get_comment/"
	urlLiveLikeCount = "live/%s/get_like_count/"
	urlLiveHeartbeat = "live/%s/heartbeat_and_get_viewer_count/"

	// IGTV
	urlIGTVChannel = "igtv/channel/"
	urlIGTVSeen    = "igtv/write_seen_state/"

	// Discover
	urlDiscoverExplore = "discover/topical_explore/"

	// Comments
	urlCommentAdd        = "media/%d/comment/"
	urlCommentDelete     = "media/%s/comment/%s/delete/"
	urlCommentBulkDelete = "media/%s/comment/bulk_delete/"
	urlCommentSync       = "media/%s/comments/"
	urlCommentDisable    = "media/%s/disable_comments/"
	urlCommentEnable     = "media/%s/enable_comments/"
	urlCommentLike       = "media/%s/comment_like/"
	urlCommentUnlike     = "media/%s/comment_unlike/"
	urlCommentOffensive  = "media/comment/check_offensive_comment/"

	// Activity
	urlActivityFollowing = "news/"
	urlActivityRecent    = "news/inbox/"

	// Inbox
	urlInbox         = "direct_v2/inbox/"
	urlInboxPending  = "direct_v2/pending_inbox/"
	urlInboxSend     = "direct_v2/threads/broadcast/text/"
	urlInboxSendLike = "direct_v2/threads/broadcast/like/"
	urlReplyStory    = "direct_v2/threads/broadcast/reel_share/"
	urlInboxThread   = "direct_v2/threads/%s/"
	urlInboxMute     = "direct_v2/threads/%s/mute/"
	urlInboxUnmute   = "direct_v2/threads/%s/unmute/"

	// Tags
	urlTagInfo    = "tags/%s/info/"
	urlTagStories = "tags/%s/story/"
	urlTagContent = "tags/%s/ranked_sections/"

	// Upload
	urlUploadPhoto      = "rupload_igphoto/%s"
	urlUploadVideo      = "rupload_igvideo/%s"
	urlConfigure        = "media/configure/"
	urlConfigureSidecar = "media/configure_sidecar/"
	urlConfigureIGTV    = "media/configure_to_igtv/?video=1"
	urlConfigureStory   = "media/configure_to_story/"

	// 2FA
	url2FACheckTrusted = "two_factor/check_trusted_notification_status/"
	url2FALogin        = "accounts/two_factor_login/"
)

// Errors
var (
	RespErr2FA = "two_factor_required"

	// Account & Login Errors
	ErrBadPassword     = errors.New("Password is incorrect")
	ErrTooManyRequests = errors.New("Too many requests, please wait a few minutes before you try again")

	// Upload Errors
	ErrInvalidFormat      = errors.New("Invalid file type, please use one of jpeg, jpg, mp4")
	ErrCarouselType       = errors.New("Invalid file type, please use a jpeg or jpg image")
	ErrCarouselMediaLimit = errors.New("Carousel media limit of 10 exceeded")
	ErrStoryBadMediaType  = errors.New("When uploading multiple items to your story at once, all have to be mp4")
	ErrStoryMediaTooLong  = errors.New("Story media must not exceed 15 seconds per item")

	// Search Errors
	ErrSearchUserNotFound = errors.New("User not found in search result")

	// Feed Errors
	ErrInvalidTab   = errors.New("Invalid tab, please select top or recent")
	ErrNoMore       = errors.New("No more posts availible, page end has been reached")
	ErrNotHighlight = errors.New("Unable to sync, Reel is not of type highlight")

	// Misc
	ErrByteIndexNotFound = errors.New("Failed to index byte slice, delim not found")
)
