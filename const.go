package goinsta

const (
	// urls
	baseUrl        = "https://i.instagram.com/"
	instaAPIUrl    = "https://i.instagram.com/api/v1/"
	instaAPIUrlb   = "https://b.i.instagram.com/api/v1/"
	instaAPIUrlv2  = "https://i.instagram.com/api/v2/"
	instaAPIUrlv2b = "https://b.i.instagram.com/api/v2/"

	// header values
	instaUserAgent     = "Instagram 107.0.0.27.121 Android (24/7.0; 380dpi; 1080x1920; OnePlus; ONEPLUS A3010; OnePlus3T; qcom; en_US)"
	instaBloksVerID    = "927f06374b80864ae6a0b04757048065714dc50ff15d2b8b3de8d0b6de961649"
	fbAnalytics        = "567067343352427"
	igCapabilities     = "3brTvx0="
	connType           = "WIFI"
	instaSigKeyVersion = "4"
	locale             = "en_US"
	appVersion         = "195.0.0.31.123"
	appVersionCode     = "302733750"

	// Used for tray requests
	supportedSdkVersions = "100.0,101.0,102.0,103.0,104.0,105.0,106.0,107.0,108.0,109.0,110.0,111.0,112.0,113.0,114.0,115.0,116.0,117.0"
	facetrackerVersion   = "14"
	segmentation         = "segmentation_enabled"
	compression          = "ETC2_COMPRESSION"
	worldTracker         = "world_tracker_enabled"
	gyroscope            = "gyroscope_enabled"
)

var (
	goInstaDeviceSettings = map[string]interface{}{
		"manufacturer":      "samsung",
		"model":             "SM-G975F",
		"code_name":         "beyond2",
		"android_version":   30,
		"android_release":   "11",
		"screen_dpi":        "560dpi",
		"screen_resolution": "1440x2898",
		"chipset":           "exynos9820",
	}
	timeOffset = getTimeOffset()
	userAgent  = createUserAgent()
)

type muteOption string

const (
	MuteAll   muteOption = "all"
	MuteStory muteOption = "story"
	MuteFeed  muteOption = "feed"
)

// Endpoints (with format vars)
const (
	// login
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

	// account
	urlCurrentUser      = "accounts/current_user/"
	urlChangePass       = "accounts/change_password/"
	urlSetPrivate       = "accounts/set_private/"
	urlSetPublic        = "accounts/set_public/"
	urlRemoveProfPic    = "accounts/remove_profile_picture/"
	urlChangeProfPic    = "accounts/change_profile_picture/"
	urlFeedSaved        = "feed/saved/"
	urlSetBiography     = "accounts/set_biography/"
	urlEditProfile      = "accounts/edit_profile/"
	urlFeedLiked        = "feed/liked/"
	urlConsent          = "consent/existing_user_flow/"
	urlNotifBadge       = "notifications/badge/"
	urlFeaturedAccounts = "multiple_accounts/get_featured_accounts/"

	// account and profile
	urlFollowers = "friendships/%d/followers/"
	urlFollowing = "friendships/%d/following/"

	// users
	urlUserArchived      = "feed/only_me_feed/"
	urlUserByName        = "users/%s/usernameinfo/"
	urlUserByID          = "users/%d/info/"
	urlUserBlock         = "friendships/block/%d/"
	urlUserUnblock       = "friendships/unblock/%d/"
	urlUserMute          = "friendships/mute_posts_or_story_from_follow/"
	urlUserUnmute        = "friendships/unmute_posts_or_story_from_follow/"
	urlUserFollow        = "friendships/create/%d/"
	urlUserUnfollow      = "friendships/destroy/%d/"
	urlUserFeed          = "feed/user/%d/"
	urlFriendship        = "friendships/show/%d/"
	urlFriendshipPending = "friendships/pending/"
	urlUserStories       = "feed/user/%d/reel_media/"
	urlUserTags          = "usertags/%d/feed/"
	urlBlockedList       = "users/blocked_list/"
	urlUserInfo          = "users/%d/info/"
	urlUserHighlights    = "highlights/%d/highlights_tray/"

	// timeline
	urlTimeline  = "feed/timeline/"
	urlStories   = "feed/reels_tray/"
	urlReelMedia = "feed/reels_media/"

	// search
	urlSearchTop       = "fbsearch/topsearch_flat/"
	urlSearchUser      = "users/search/"
	urlSearchTag       = "tags/search/"
	urlSearchLocation  = "fbsearch/places/"
	urlSearchRecent    = "fbsearch/recent_searches/"
	urlSearchNullState = "fbsearch/nullstate_dynamic_sections/"

	// feeds
	urlFeedLocationID = "feed/location/%d/"
	urlFeedLocations  = "locations/%d/sections/"
	urlFeedTag        = "feed/tag/%s/"

	// media
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
	urlIGTVChannel = "/api/v1/igtv/channel/"

	// Discover
	urlDiscoverExplore = "discover/topical_explore/"

	// comments
	urlCommentAdd     = "media/%d/comment/"
	urlCommentDelete  = "media/%s/comment/%s/delete/"
	urlCommentSync    = "media/%s/comments/"
	urlCommentDisable = "media/%s/disable_comments/"
	urlCommentEnable  = "media/%s/enable_comments/"
	urlCommentLike    = "media/%s/comment_like/"
	urlCommentUnlike  = "media/%s/comment_unlike/"

	// activity
	urlActivityFollowing = "news/"
	urlActivityRecent    = "news/inbox/"

	// inbox
	urlInbox         = "direct_v2/inbox/"
	urlInboxPending  = "direct_v2/pending_inbox/"
	urlInboxSend     = "direct_v2/threads/broadcast/text/"
	urlInboxSendLike = "direct_v2/threads/broadcast/like/"
	urlReplyStory    = "direct_v2/threads/broadcast/reel_share/"
	urlInboxThread   = "direct_v2/threads/%s/"
	urlInboxMute     = "direct_v2/threads/%s/mute/"
	urlInboxUnmute   = "direct_v2/threads/%s/unmute/"

	// tags
	urlTagSync    = "tags/%s/info/"
	urlTagStories = "tags/%s/story/"
	urlTagContent = "tags/%s/ranked_sections/"

	// upload
	urlUploadStory = "https://i.instagram.com/rupload_igphoto/103079408575885_0_-1340379573"
)
