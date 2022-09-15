package goinsta

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"time"
)

type FollowOrder string

const (
	DefaultOrder  FollowOrder = "default"
	LatestOrder   FollowOrder = "date_followed_latest"
	EarliestOrder FollowOrder = "date_followed_earliest"
)

// Users is a struct that stores many user's returned by many different methods.
type Users struct {
	insta *Instagram

	// It's a bit confusing to have the same structure
	// in the Instagram strucure and in the multiple users
	// calls

	err      error
	endpoint string
	query    map[string]string

	Status    string          `json:"status"`
	BigList   bool            `json:"big_list"`
	Users     []*User         `json:"users"`
	PageSize  int             `json:"page_size"`
	RawNextID json.RawMessage `json:"next_max_id"`
	NextID    string          `json:"-"`
}

// SetInstagram sets new instagram to user structure
func (users *Users) SetInstagram(insta *Instagram) {
	users.insta = insta
}

// Next allows to paginate after calling:
// Account.Follow* and User.Follow*
//
// New user list is stored inside Users
//
// returns false when list reach the end.
func (users *Users) Next() bool {
	if users.err != nil {
		return false
	}

	insta := users.insta
	endpoint := users.endpoint

	query := map[string]string{}
	if users.NextID != "" {
		query["max_id"] = users.NextID
	}

	if _, ok := users.query["rank_token"]; !ok {
		users.query["rank_token"] = generateUUID()
	}

	for key, value := range users.query {
		query[key] = value
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: endpoint,
			Query:    query,
		},
	)
	if err != nil {
		users.err = err
		return false
	}

	var newUsers Users
	if err := json.Unmarshal(body, &newUsers); err != nil {
		users.err = err
		return false
	}

	// check whether the nextID contains quotes (string type) or not (int64 type)
	if len(newUsers.RawNextID) > 0 && newUsers.RawNextID[0] == '"' && newUsers.RawNextID[len(newUsers.RawNextID)-1] == '"' {
		if err := json.Unmarshal(newUsers.RawNextID, &users.NextID); err != nil {
			users.err = err
			return false
		}
	} else if newUsers.RawNextID != nil {
		var nextID int64
		if err := json.Unmarshal(newUsers.RawNextID, &nextID); err != nil {
			users.err = err
			return false
		}
		users.NextID = strconv.FormatInt(nextID, 10)
	}

	users.Status = newUsers.Status
	users.BigList = newUsers.BigList
	users.Users = newUsers.Users
	users.PageSize = newUsers.PageSize
	users.RawNextID = newUsers.RawNextID

	users.setValues()

	// Dont't return false on first error otherwise for loop won't run
	if users.NextID == "" {
		users.err = ErrNoMore
	}

	return true
}

// Error returns users error
func (users *Users) Error() error {
	return users.err
}

func (users *Users) setValues() {
	for i := range users.Users {
		users.Users[i].insta = users.insta
	}
}

type userResp struct {
	Status string `json:"status"`
	User   User   `json:"user"`
}

// User is the representation of instagram's user profile
type User struct {
	insta       *Instagram
	Collections *Collections

	// User info
	ID                     int64  `json:"pk"`
	Username               string `json:"username"`
	FullName               string `json:"full_name,omitempty"`
	Email                  string `json:"email,omitempty"`
	PhoneNumber            string `json:"phone_number,omitempty"`
	WhatsappNumber         string `json:"whatsapp_number,omitempty"`
	Gender                 int    `json:"gender,omitempty"`
	PublicEmail            string `json:"public_email,omitempty"`
	PublicPhoneNumber      string `json:"public_phone_number,omitempty"`
	PublicPhoneCountryCode string `json:"public_phone_country_code,omitempty"`
	ContactPhoneNumber     string `json:"contact_phone_number,omitempty"`

	// Profile visible properties
	IsPrivate                  bool   `json:"is_private"`
	IsVerified                 bool   `json:"is_verified"`
	ExternalURL                string `json:"external_url,omitempty"`
	ExternalLynxURL            string `json:"external_lynx_url,omitempty"`
	FollowerCount              int    `json:"follower_count"`
	FollowingCount             int    `json:"following_count"`
	ProfilePicID               string `json:"profile_pic_id,omitempty"`
	ProfilePicURL              string `json:"profile_pic_url,omitempty"`
	HasAnonymousProfilePicture bool   `json:"has_anonymous_profile_picture"`
	Biography                  string `json:"biography,omitempty"`
	BiographyWithEntities      struct {
		RawText  string        `json:"raw_text"`
		Entities []interface{} `json:"entities"`
	} `json:"biography_with_entities"`
	BiographyProductMentions []interface{} `json:"biography_products_mentions"`

	// Profile hidden properties
	IsNeedy                        bool          `json:"is_needy"`
	IsInterestAccount              bool          `json:"is_interest_account"`
	IsVideoCreator                 bool          `json:"is_video_creator"`
	IsBusiness                     bool          `json:"is_business"`
	BestiesCount                   int           `json:"besties_count"`
	ShowBestiesBadge               bool          `json:"show_besties_badge"`
	RecentlyBestiedByCount         int           `json:"recently_bestied_by_count"`
	AccountType                    int           `json:"account_type"`
	AccountBadges                  []interface{} `json:"account_badges,omitempty"`
	FbIdV2                         int64         `json:"fbid_v2"`
	IsUnpublished                  bool          `json:"is_unpublished"`
	UserTagsCount                  int           `json:"usertags_count"`
	UserTagReviewEnabled           bool          `json:"usertag_review_enabled"`
	FollowingTagCount              int           `json:"following_tag_count"`
	MutualFollowersID              []int64       `json:"profile_context_mutual_follow_ids,omitempty"`
	FollowFrictionType             int           `json:"follow_friction_type"`
	ProfileContext                 string        `json:"profile_context,omitempty"`
	HasBiographyTranslation        bool          `json:"has_biography_translation"`
	HasSavedItems                  bool          `json:"has_saved_items"`
	Nametag                        Nametag       `json:"nametag,omitempty"`
	HasChaining                    bool          `json:"has_chaining"`
	IsFavorite                     bool          `json:"is_favorite"`
	IsFavoriteForStories           bool          `json:"is_favorite_for_stories"`
	IsFavoriteForHighlights        bool          `json:"is_favorite_for_highlights"`
	IsProfileActionNeeded          bool          `json:"is_profile_action_needed"`
	CanBeReportedAsFraud           bool          `json:"can_be_reported_as_fraud"`
	CanSeeSupportInbox             bool          `json:"can_see_support_inbox"`
	CanSeeSupportInboxV1           bool          `json:"can_see_support_inbox_v1"`
	CanSeePrimaryCountryInsettings bool          `json:"can_see_primary_country_in_settings"`
	CanFollowHashtag               bool          `json:"can_follow_hashtag"`

	// Business profile properies
	CanBoostPosts                  bool   `json:"can_boost_posts"`
	CanSeeOrganicInsights          bool   `json:"can_see_organic_insights"`
	CanConvertToBusiness           bool   `json:"can_convert_to_business"`
	CanCreateSponsorTags           bool   `json:"can_create_sponsor_tags"`
	CanCreateNewFundraiser         bool   `json:"can_create_new_standalone_fundraiser"`
	CanCreateNewPersonalFundraiser bool   `json:"can_create_new_standalone_personal_fundraiser"`
	CanBeTaggedAsSponsor           bool   `json:"can_be_tagged_as_sponsor"`
	PersonalAccountAdsPageName     string `json:"personal_account_ads_page_name,omitempty"`
	PersonalAccountAdsId           string `json:"personal_account_ads_page_id,omitempty"`
	Category                       string `json:"category,omitempty"`

	// Shopping properties
	ShowShoppableFeed           bool `json:"show_shoppable_feed"`
	CanTagProductsFromMerchants bool `json:"can_tag_products_from_merchants"`

	// Miscellaneous
	IsMutedWordsGlobalEnabled bool   `json:"is_muted_words_global_enabled"`
	IsMutedWordsCustomEnabled bool   `json:"is_muted_words_custom_enabled"`
	AllowedCommenterType      string `json:"allowed_commenter_type,omitempty"`

	// Media properties
	MediaCount          int   `json:"media_count"`
	IGTVCount           int   `json:"total_igtv_videos"`
	HasIGTVSeries       bool  `json:"has_igtv_series"`
	HasVideos           bool  `json:"has_videos"`
	TotalClipCount      int   `json:"total_clips_count"`
	TotalAREffects      int   `json:"total_ar_effects"`
	GeoMediaCount       int   `json:"geo_media_count"`
	HasProfileVideoFeed bool  `json:"has_profile_video_feed"`
	LiveBroadcastID     int64 `json:"live_broadcast_id"`

	HasPlacedOrders bool `json:"has_placed_orders"`

	ShowInsightTerms           bool `json:"show_insights_terms"`
	ShowConversionEditEntry    bool `json:"show_conversion_edit_entry"`
	ShowPostsInsightEntryPoint bool `json:"show_post_insights_entry_point"`
	ShoppablePostsCount        int  `json:"shoppable_posts_count"`
	RequestContactEnabled      bool `json:"request_contact_enabled"`
	FeedPostReshareDisabled    bool `json:"feed_post_reshare_disabled"`
	CreatorShoppingInfo        struct {
		LinkedMerchantAccounts []interface{} `json:"linked_merchant_accounts,omitempty"`
	} `json:"creator_shopping_info,omitempty"`
	StandaloneFundraiserInfo struct {
		HasActiveFundraiser                 bool        `json:"has_active_fundraiser"`
		FundraiserId                        int64       `json:"fundraiser_id"`
		FundraiserTitle                     string      `json:"fundraiser_title"`
		FundraiserType                      interface{} `json:"fundraiser_type"`
		FormattedGoalAmount                 string      `json:"formatted_goal_amount"`
		BeneficiaryUsername                 string      `json:"beneficiary_username"`
		FormattedFundraiserProgressInfoText string      `json:"formatted_fundraiser_progress_info_text"`
		PercentRaised                       interface{} `json:"percent_raised"`
	} `json:"standalone_fundraiser_info"`
	AggregatePromoteEngagement   bool         `json:"aggregate_promote_engagement"`
	AllowMentionSetting          string       `json:"allow_mention_setting,omitempty"`
	AllowTagSetting              string       `json:"allow_tag_setting,omitempty"`
	LimitedInteractionsEnabled   bool         `json:"limited_interactions_enabled"`
	ReelAutoArchive              string       `json:"reel_auto_archive,omitempty"`
	HasHighlightReels            bool         `json:"has_highlight_reels"`
	HightlightReshareDisabled    bool         `json:"highlight_reshare_disabled"`
	IsMemorialized               bool         `json:"is_memorialized"`
	HasGuides                    bool         `json:"has_guides"`
	HasAffiliateShop             bool         `json:"has_active_affiliate_shop"`
	CityID                       int64        `json:"city_id"`
	CityName                     string       `json:"city_name,omitempty"`
	AddressStreet                string       `json:"address_street,omitempty"`
	DirectMessaging              string       `json:"direct_messaging,omitempty"`
	Latitude                     float64      `json:"latitude"`
	Longitude                    float64      `json:"longitude"`
	BusinessContactMethod        string       `json:"business_contact_method"`
	IncludeDirectBlacklistStatus bool         `json:"include_direct_blacklist_status"`
	HdProfilePicURLInfo          PicURLInfo   `json:"hd_profile_pic_url_info,omitempty"`
	HdProfilePicVersions         []PicURLInfo `json:"hd_profile_pic_versions,omitempty"`
	School                       School       `json:"school"`
	Byline                       string       `json:"byline"`
	SocialContext                string       `json:"social_context,omitempty"`
	SearchSocialContext          string       `json:"search_social_context,omitempty"`
	MutualFollowersCount         float64      `json:"mutual_followers_count"`
	LatestReelMedia              int64        `json:"latest_reel_media,omitempty"`
	IsCallToActionEnabled        bool         `json:"is_call_to_action_enabled"`
	IsPotentialBusiness          bool         `json:"is_potential_business"`
	FbPageCallToActionID         string       `json:"fb_page_call_to_action_id,omitempty"`
	FbPayExperienceEnabled       bool         `json:"fbpay_experience_enabled"`
	Zip                          string       `json:"zip,omitempty"`
	Friendship                   Friendship   `json:"friendship_status"`
	AutoExpandChaining           bool         `json:"auto_expand_chaining"`

	AllowedToCreateNonprofitFundraisers        bool          `json:"is_allowed_to_create_standalone_nonprofit_fundraisers"`
	AllowedToCreatePersonalFundraisers         bool          `json:"is_allowed_to_create_standalone_personal_fundraisers"`
	IsElegibleToShowFbCrossSharingNux          bool          `json:"is_eligible_to_show_fb_cross_sharing_nux"`
	PageIdForNewSumaBizAccount                 interface{}   `json:"page_id_for_new_suma_biz_account"`
	ElegibleShoppingSignupEntrypoints          []interface{} `json:"eligible_shopping_signup_entrypoints"`
	IsIgdProductPickerEnabled                  bool          `json:"is_igd_product_picker_enabled"`
	IsElegibleForAffiliateShopOnboarding       bool          `json:"is_eligible_for_affiliate_shop_onboarding"`
	IsElegibleForSMBSupportFlow                bool          `json:"is_eligible_for_smb_support_flow"`
	ElegibleShoppingFormats                    []interface{} `json:"eligible_shopping_formats"`
	NeedsToAcceptShoppingSellerOnboardingTerms bool          `json:"needs_to_accept_shopping_seller_onboarding_terms"`
	SellerShoppableFeedType                    string        `json:"seller_shoppable_feed_type"`
	VisibleProducts                            int           `json:"num_visible_storefront_products"`
	IsShoppingCatalogSettingsEnabled           bool          `json:"is_shopping_settings_enabled"`
	IsShoppingCommunityContentEnabled          bool          `json:"is_shopping_community_content_enabled"`
	IsShoppingAutoHighlightEnabled             bool          `json:"is_shopping_auto_highlight_eligible"`
	IsShoppingCatalogSourceSelectionEnabled    bool          `json:"is_shopping_catalog_source_selection_enabled"`
	ProfessionalConversionSuggestedAccountType int           `json:"professional_conversion_suggested_account_type"`
	InteropMessagingUserfbid                   int64         `json:"interop_messaging_user_fbid"`
	LinkedFbInfo                               struct{}      `json:"linked_fb_info"`
	HasElegibleWhatsappLinkingCategory         struct{}      `json:"has_eligible_whatsapp_linking_category"`
	ExistingUserAgeCollectionEnabled           bool          `json:"existing_user_age_collection_enabled"`
	AboutYourAccountBloksEntrypointEnabled     bool          `json:"about_your_account_bloks_entrypoint_enabled"`
	OpenExternalUrlWithInAppBrowser            bool          `json:"open_external_url_with_in_app_browser"`
	MerchantCheckoutStyle                      string        `json:"merchant_checkout_style"`

	// Profile picture as raw bytes, to populate call User.DownloadProfilePic()
	ProfilePic []byte
}

// SetInstagram will update instagram instance for selected User.
func (user *User) SetInstagram(insta *Instagram) {
	user.insta = insta
}

// NewUser returns prepared user to be used with his functions.
func (insta *Instagram) NewUser() *User {
	return &User{insta: insta}
}

// Info updates user info
// extra query arguments can be passes one after another as func(key, value).
// Only if an even number of string arguements will be passed, they will be
//   used in the query.
//
// See example: examples/user/friendship.go
func (user *User) Info(params ...interface{}) error {
	insta := user.insta
	query := map[string]string{}
	if len(params)%2 == 0 {
		for i := 0; i < len(params); i = i + 2 {
			query[params[i].(string)] = params[i+1].(string)
		}
	}

	body, _, err := insta.sendRequest(&reqOptions{
		Endpoint: fmt.Sprintf(urlUserInfo, user.ID),
		Query:    query,
	})
	if err != nil {
		return err
	}
	result := struct {
		User   *User  `json:"user"`
		Status string `json:"status"`
	}{
		User: user,
	}

	err = json.Unmarshal(body, &result)
	return err
}

// Sync wraps User.Info() 1:1
func (user *User) Sync(params ...interface{}) error {
	return user.Info(params...)
}

// Following returns a list of user following.
//
// Query can be used to search for a specific user.
// Be aware that it only matches from the start, e.g.
// "theprimeagen" will only match "theprime" not "prime".
// To fetch all user an empty string "".
//
// Users.Next can be used to paginate
func (user *User) Following(query string, order FollowOrder) *Users {
	return user.followList(urlFollowing, query, order)
}

// Followers returns a list of user followers.
//
// Query can be used to search for a specific user.
// Be aware that it only matches from the start, e.g.
// "theprimeagen" will only match "theprime" not "prime".
// To fetch all user an empty string "".
//
// Users.Next can be used to paginate
func (user *User) Followers(query string) *Users {
	return user.followList(urlFollowers, query, DefaultOrder)
}

func (user *User) followList(url, query string, order FollowOrder) *Users {
	users := Users{
		insta:    user.insta,
		endpoint: fmt.Sprintf(url, user.ID),
		query: map[string]string{
			"search_surface": "follow_list_page",
			"query":          query,
			"enable_groups":  "true",
		},
	}

	if order != DefaultOrder {
		users.query["order"] = string(order)
	}

	if url == urlFollowing {
		users.query["includes_hashtags"] = "true"
	}

	return &users
}

// Block blocks user
//
// This function updates current User.Friendship structure.
// Param: autoBlock - automatically block accounts registered on the same email/number
//
// See example: examples/user/block.go
func (user *User) Block(autoBlock bool) error {
	insta := user.insta
	data, err := json.Marshal(map[string]string{
		"surface":              "profile",
		"is_autoblock_enabled": strconv.FormatBool(autoBlock),
		"user_id":              strconv.FormatInt(user.ID, 10),
		"_uid":                 strconv.FormatInt(insta.Account.ID, 10),
		"_uuid":                insta.uuid,
	})
	if err != nil {
		return err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlUserBlock, user.ID),
			IsPost:   true,
			Query: map[string]string{
				"signed_body": "SIGNATURE." + string(data),
			},
		},
	)
	if err != nil {
		return err
	}
	resp := friendResp{}
	err = json.Unmarshal(body, &resp)
	user.Friendship = resp.Friendship
	if err != nil {
		return err
	}

	return nil
}

// Unblock unblocks user
//
// This function updates current User.Friendship structure.
//
// See example: examples/user/unblock.go
func (user *User) Unblock() error {
	insta := user.insta
	data, err := json.Marshal(
		map[string]interface{}{
			"user_id":          user.ID,
			"_uid":             insta.Account.ID,
			"_uuid":            insta.uuid,
			"container_module": "blended_search",
		},
	)
	if err != nil {
		return err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlUserUnblock, user.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err != nil {
		return err
	}
	resp := friendResp{}
	err = json.Unmarshal(body, &resp)
	user.Friendship = resp.Friendship
	if err != nil {
		return err
	}

	return nil
}

// Mute mutes user from appearing in the feed or story reel
//
// Use one of the pre-defined constants to choose what exactly to mute:
// goinsta.MuteAll, goinsta.MuteStory, goinsta.MuteFeed
// This function updates current User.Friendship structure.
func (user *User) Mute(opt muteOption) error {
	if opt == MuteAll {
		err := user.muteOrUnmute(MuteStory, urlUserMute)
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		return user.muteOrUnmute(MutePosts, urlUserMute)
	}
	return user.muteOrUnmute(opt, urlUserMute)
}

// Unmute unmutes user so it appears in the feed or story reel again
//
// Use one of the pre-defined constants to choose what exactly to unmute:
// goinsta.MuteAll, goinsta.MuteStory, goinsta.MuteFeed
// This function updates current User.Friendship structure.
func (user *User) Unmute(opt muteOption) error {
	if opt == MuteAll {
		err := user.muteOrUnmute(MuteStory, urlUserUnmute)
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		return user.muteOrUnmute(MutePosts, urlUserUnmute)
	}
	return user.muteOrUnmute(opt, urlUserUnmute)
}

func (user *User) muteOrUnmute(opt muteOption, endpoint string) error {
	insta := user.insta
	data, err := json.Marshal(generateMuteData(user, opt))
	if err != nil {
		return err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: endpoint,
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err != nil {
		return err
	}
	resp := friendResp{}
	err = json.Unmarshal(body, &resp)
	user.Friendship = resp.Friendship
	if err != nil {
		return err
	}

	return nil
}

func generateMuteData(user *User, opt muteOption) map[string]string {
	insta := user.insta
	data := map[string]string{
		"_uid":             toString(insta.Account.ID),
		"_uuid":            insta.uuid,
		"container_module": "media_mute_sheet",
	}

	switch opt {
	case MuteStory:
		data["target_reel_author_id"] = toString(user.ID)
	case MutePosts:
		data["target_posts_author_id"] = toString(user.ID)
	}

	return data
}

// Follow started following some user
//
// This function performs a follow call. If user is private
// you have to wait until he/she accepts you.
//
// If the account is public User.Friendship will be updated
//
// See example: examples/user/follow.go
func (user *User) Follow() error {
	insta := user.insta
	data, err := json.Marshal(
		map[string]string{
			"user_id":    toString(user.ID),
			"radio_type": "wifi-none",
			"_uid":       toString(insta.Account.ID),
			"device_id":  insta.dID,
			"_uuid":      insta.uuid,
		},
	)
	if err != nil {
		return err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlUserFollow, user.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err != nil {
		return err
	}
	resp := friendResp{}
	err = json.Unmarshal(body, &resp)
	user.Friendship = resp.Friendship
	if err != nil {
		return err
	}

	return nil
}

// Unfollow unfollows user
//
// User.Friendship will be updated
//
// See example: examples/user/unfollow.go
func (user *User) Unfollow() error {
	insta := user.insta
	data, err := json.Marshal(
		map[string]string{
			"user_id":          toString(user.ID),
			"radio_type":       "wifi-none",
			"_uid":             toString(insta.Account.ID),
			"device_id":        insta.dID,
			"_uuid":            insta.uuid,
			"container_module": "following_sheet",
		},
	)
	if err != nil {
		return err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlUserUnfollow, user.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err != nil {
		return err
	}
	resp := friendResp{}
	err = json.Unmarshal(body, &resp)
	user.Friendship = resp.Friendship
	if err != nil {
		return err
	}

	return nil
}

// GetFriendship allows user to get friend relationship.
//
// The result is stored in user.Friendship
func (user *User) GetFriendship() (fr *Friendship, err error) {
	insta := user.insta
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlFriendship, user.ID),
		},
	)
	if err == nil {
		fr = &user.Friendship
		err = json.Unmarshal(body, fr)
	}
	return
}

// GetFeaturedAccounts will call the featured accounts enpoint.
func (user *User) GetFeaturedAccounts() ([]*User, error) {
	body, _, err := user.insta.sendRequest(&reqOptions{
		Endpoint: urlFeaturedAccounts,
		Query: map[string]string{
			"target_user_id": strconv.FormatInt(user.ID, 10),
		},
	})
	if err != nil {
		return nil, err
	}

	d := struct {
		Accounts []*User `json:"accounts"`
		Status   string  `json:"status"`
	}{}
	err = json.Unmarshal(body, &d)
	return d.Accounts, err
}

// Feed returns user feeds (media)
//
// 	params can be:
// 		string: timestamp of the minimum media timestamp.
//
// For pagination use FeedMedia.Next()
//
// See example: examples/user/feed.go
func (user *User) Feed(params ...interface{}) *FeedMedia {
	insta := user.insta

	media := &FeedMedia{
		insta:    insta,
		endpoint: urlUserFeed,
		uid:      user.ID,
	}

	for _, param := range params {
		switch s := param.(type) {
		case string:
			media.timestamp = s
		}
	}

	return media
}

// Tags returns media where user is tagged in
//
// For pagination use FeedMedia.Next()
//
// See example: examples/user/tags.go
func (user *User) Tags(minTimestamp []byte) (*FeedMedia, error) {
	insta := user.insta

	timestamp := string(minTimestamp)
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlUserTags, user.ID),
			Query: map[string]string{
				"max_id":         "",
				"rank_token":     user.insta.rankToken,
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
		uid:      user.ID,
	}
	err = json.Unmarshal(body, media)
	if err != nil {
		return nil, err
	}
	media.setValues()
	return media, nil
}

// DownloadProfilePic will download a user's profile picture if available, and
//   return it as a byte slice.
func (user *User) DownloadProfilePic() ([]byte, error) {
	if user.ProfilePicURL == "" {
		return nil, ErrNoProfilePicURL
	}
	insta := user.insta
	b, err := insta.download(user.ProfilePicURL)
	if err != nil {
		return nil, err
	}
	user.ProfilePic = b
	return b, nil
}

// DownloadProfilePicTo will download the user profile picture to the provided
//   path. If path does not include a file name, one will be extracted automatically.
// File extention does not need to be set, and will be set automatically.
func (user *User) DownloadProfilePicTo(dst string) error {
	folder, fn := path.Split(dst)
	b, err := user.DownloadProfilePic()
	if err != nil {
		return nil
	}
	fn, err = getDownloadName(user.ProfilePicURL, fn)
	if err != nil {
		return err
	}
	err = saveToFolder(folder, fn, b)
	return err
}

func (user *User) ApprovePending() error {
	return user.changePending(urlFriendshipApprove)
}

func (user *User) IgnorePending() error {
	return user.changePending(urlFriendshipIgnore)
}

func (user *User) changePending(endpoint string) error {
	insta := user.insta
	if !user.Friendship.IncomingRequest {
		return ErrNoPendingFriendship
	}

	query := map[string]string{
		"user_id":          toString(user.ID),
		"radio_type":       "wifi-none",
		"_uid":             toString(insta.Account.ID),
		"_uuid":            insta.uuid,
		"nav_chain":        "EAS:newsfeed_you:52,DmK:follow_requests:61",
		"container_module": "follow_requests",
	}

	data, err := json.Marshal(query)
	if err != nil {
		return err
	}
	resp, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(endpoint, user.ID),
			IsPost:   true,
			Query:    generateSignature(data),
		},
	)
	if err != nil {
		return err
	}

	var result friendResp
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return err
	}

	if result.Status != "ok" {
		return fmt.Errorf("bad status: %s", result.Status)
	}
	user.Friendship = result.Friendship
	return nil
}
