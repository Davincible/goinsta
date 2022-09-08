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
	appVersion         = "250.0.0.21.109"
	appVersionCode     = "394071253"

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
	omitAPIHeadersExclude = []string{
		"X-Ig-Bandwidth-Speed-Kbps",
		"Ig-U-Shbts",
		"X-Ig-Mapped-Locale",
		"X-Ig-Family-Device-Id",
		"X-Ig-Android-Id",
		"X-Ig-Timezone-Offset",
		"X-Ig-Device-Locale",
		"X-Ig-Device-Id",
		"Ig-Intended-User-Id",
		"X-Ig-App-Locale",
		"X-Bloks-Is-Layout-Rtl",
		"X-Pigeon-Rawclienttime",
		"X-Bloks-Version-Id",
		"X-Ig-Bandwidth-Totalbytes-B",
		"X-Ig-Bandwidth-Totaltime-Ms",
		"X-Ig-App-Startup-Country",
		"X-Ig-Www-Claim",
		"X-Bloks-Is-Panorama-Enabled",
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
	urlUserArchived           = "feed/only_me_feed/"
	urlUserByName             = "users/%s/usernameinfo/"
	urlUserByID               = "users/%s/info/"
	urlUserBlock              = "friendships/block/%d/"
	urlUserUnblock            = "friendships/unblock/%d/"
	urlUserMute               = "friendships/mute_posts_or_story_from_follow/"
	urlUserUnmute             = "friendships/unmute_posts_or_story_from_follow/"
	urlUserFollow             = "friendships/create/%d/"
	urlUserUnfollow           = "friendships/destroy/%d/"
	urlUserFeed               = "feed/user/%d/"
	urlFriendship             = "friendships/show/%d/"
	urlFriendshipShowMany     = "friendships/show_many/"
	urlFriendshipPending      = "friendships/pending/"
	urlFriendshipPendingCount = "friendships/pending_follow_requests_count/"
	urlFriendshipApprove      = "friendships/approve/%d/"
	urlFriendshipIgnore       = "friendships/ignore/%d/"
	urlUserStories            = "feed/user/%d/story/"
	urlUserTags               = "usertags/%d/feed/"
	urlBlockedList            = "users/blocked_list/"
	urlUserInfo               = "users/%d/info/"
	urlUserHighlights         = "highlights/%d/highlights_tray/"

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
	urlFeedLocationID    = "feed/location/%d/"
	urlFeedLocations     = "locations/%d/sections/"
	urlFeedTag           = "feed/tag/%s/"
	urlFeedNewPostsExist = "feed/new_feed_posts_exist/"

	// Media
	urlMediaInfo         = "media/%s/info/"
	urlMediaDelete       = "media/%s/delete/"
	urlMediaLike         = "media/%s/like/"
	urlMediaUnlike       = "media/%s/unlike/"
	urlMediaSave         = "media/%s/save/"
	urlMediaUnsave       = "media/%s/unsave/"
	urlMediaSeen         = "media/seen/"
	urlMediaLikers       = "media/%s/likers/"
	urlMediaBlocked      = "media/blocked/"
	urlMediaCommentInfos = "media/comment_infos/"

	// Broadcasts
	urlLiveInfo      = "live/%d/info/"
	urlLiveComments  = "live/%d/get_comment/"
	urlLiveLikeCount = "live/%d/get_like_count/"
	urlLiveHeartbeat = "live/%d/heartbeat_and_get_viewer_count/"
	urlLiveChaining  = "live/get_live_chaining/"

	// IGTV
	urlIGTVDiscover = "igtv/discover/"
	urlIGTVChannel  = "igtv/channel/"
	urlIGTVSeries   = "igtv/series/all_user_series/%d/"
	urlIGTVSeen     = "igtv/write_seen_state/"

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
	urlActivitySeen      = "news/inbox_seen/"

	// Inbox
	urlInbox             = "direct_v2/inbox/"
	urlInboxPending      = "direct_v2/pending_inbox/"
	urlInboxSend         = "direct_v2/threads/broadcast/text/"
	urlInboxSendLike     = "direct_v2/threads/broadcast/like/"
	urlReplyStory        = "direct_v2/threads/broadcast/reel_share/"
	urlGetByParticipants = "direct_v2/threads/get_by_participants/"
	urlInboxThread       = "direct_v2/threads/%s/"
	urlInboxMute         = "direct_v2/threads/%s/mute/"
	urlInboxUnmute       = "direct_v2/threads/%s/unmute/"
	urlInboxGetItems     = "direct_v2/threads/%s/get_items/"
	urlInboxMsgSeen      = "direct_v2/threads/%s/items/%s/seen/"
	urlInboxApprove      = "direct_v2/threads/%s/approve/"
	urlInboxHide         = "direct_v2/threads/%s/hide/"

	// Tags
	urlTagInfo    = "tags/%s/info/"
	urlTagStories = "tags/%s/story/"
	urlTagContent = "tags/%s/sections/"

	// Upload
	urlUploadPhoto      = "rupload_igphoto/%s"
	urlUploadVideo      = "rupload_igvideo/%s"
	urlUploadFinishVid  = "media/upload_finish/?video=1"
	urlConfigure        = "media/configure/"
	urlConfigureClip    = "media/configure_to_clips/"
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
	ErrBadPassword     = errors.New("password is incorrect")
	ErrTooManyRequests = errors.New("too many requests, please wait a few minutes before you try again")
	ErrLoggedOut       = errors.New("you have been logged out, please log back in")
	ErrLoginRequired   = errors.New("you are not logged in, please login")
	ErrSessionNotSet   = errors.New("session identifier is not set, please log in again to set it")
	ErrLogoutFailed    = errors.New("failed to logout")

	ErrChallengeRequired  = errors.New("challenge required")
	ErrCheckpointRequired = errors.New("checkpoint required")
	ErrCheckpointPassed   = errors.New("a checkpoint was thrown, but goinsta managed to solve it. Please call the function again")
	ErrChallengeFailed    = errors.New("failed to solve challenge automatically")

	Err2FARequired = errors.New("two Factor Autentication required. Please call insta.TwoFactorInfo.Login2FA(code)")
	Err2FANoCode   = errors.New("2FA seed is not set, and no code was provided. Please do atleast one of them")
	ErrInvalidCode = errors.New("the security code provided is incorrect")

	// Upload Errors
	ErrInvalidFormat      = errors.New("invalid file type, please use one of jpeg, jpg, mp4")
	ErrInvalidImage       = errors.New("invalid file type, please use a jpeg or jpg image")
	ErrCarouselType       = ErrInvalidImage
	ErrCarouselMediaLimit = errors.New("carousel media limit of 10 exceeded")
	ErrStoryBadMediaType  = errors.New("when uploading multiple items to your story at once, all have to be mp4")
	ErrStoryMediaTooLong  = errors.New("story media must not exceed 15 seconds per item")

	// Search Errors
	ErrSearchUserNotFound = errors.New("User not found in search result")

	// IGTV
	ErrIGTVNoSeries = errors.New(
		"User has no IGTV series, unable to fetch. If you think this was a mistake please update the user",
	)

	// Feed Errors
	ErrInvalidTab   = errors.New("invalid tab, please select top or recent")
	ErrNoMore       = errors.New("no more posts availible, page end has been reached")
	ErrNotHighlight = errors.New("unable to sync, Reel is not of type highlight")
	ErrMediaDeleted = errors.New("sorry, this media has been deleted")

	// Inbox
	ErrConvNotPending = errors.New("unable to perform action, conversation is not pending")

	// Misc
	ErrByteIndexNotFound = errors.New("failed to index byte slice, delim not found")
	ErrNoMedia           = errors.New("failed to download, no media found")
	ErrInstaNotDefined   = errors.New(
		"insta has not been defined, this is most likely a bug in the code. Please backtrack which call this error came from, and open an issue detailing exactly how you got to this error",
	)
	ErrNoValidLogin    = errors.New("no valid login found")
	ErrNoProfilePicURL = errors.New("no profile picture url was found. Please fetch the profile first")

	// Users
	ErrNoPendingFriendship = errors.New("unable to approve or ignore friendship for user, as there is no pending friendship request")

	// Headless
	ErrChromeNotFound = errors.New("to solve challenges a (headless) Chrome browser is used, but none was found. Please install Chromium or Google Chrome, and try again")
)
