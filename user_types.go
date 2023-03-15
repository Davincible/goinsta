package goinsta

type User struct {
	insta       *Instagram
	Collections *Collections

	// Profile picture as raw bytes, to populate call User.DownloadProfilePic()
	ProfilePicBytes []byte

	*UserInfo
	*UserStats
	*BiographyInfo
	*ContactInfo
	*Friendship `json:"friendship_status,omitempty"`
	*ProfilePic
	*Shop
	*Favorites
	*Business
	*AdsAccount
	*Category
	*ProfileContext
	*Fundraiser
	*FanClub `json:"fan_club_info,omitempty"`
	*ExtraUnknown
}

// HasUserInfo returns true if user info has been set.
func (u *User) HasUserInfo() bool {
	return u.UserInfo != nil
}

func (u *User) HasBiographyData() bool {
	return u.BiographyInfo != nil
}

func (u *User) HasUserStats() bool {
	return u.UserStats != nil
}

func (u *User) HasContactInfo() bool {
	return u.ContactInfo != nil
}

// HasProflePicData returns true if profile picture data has been set.
func (u *User) HasProflePicData() bool {
	return u.ProfilePic != nil
}

// HasShopData returns true if shop data has been set
func (u *User) HasShopData() bool {
	return u.Shop != nil
}

func (u *User) HasBusinessData() bool {
	return u.Business != nil
}

func (u *User) HasCategoryData() bool {
	return u.Category != nil
}

func (u *User) HasProfileContextData() bool {
	return u.ProfileContext != nil
}

func (u *User) HasFanClubData() bool {
	return u.FanClub != nil
}

func (u *User) HasAdsAccountData() bool {
	return u.AdsAccount != nil
}

func (u *User) HasFundraiserData() bool {
	return u.Fundraiser != nil
}

func (u *User) HasFavoritesData() bool {
	return u.Favorites != nil
}

func (u *User) HasExtraUnknownData() bool {
	return u.ExtraUnknown != nil
}

type UserInfo struct {
	ID         int64  `json:"pk"`
	IDstr      string `json:"pk_id"`
	Username   string `json:"username"`
	FullName   string `json:"full_name"`
	IsPrivate  bool   `json:"is_private"`
	IsVerified bool   `json:"is_verified"`

	*UserInfoExtended
}

type UserInfoExtended struct {
	StrongID        string `json:"strong_id__"`
	PageID          int64  `json:"page_id"`
	PageName        string `json:"page_name"`
	IsFavorite      bool   `json:"is_favorite"`
	IsBusiness      bool   `json:"is_business"`
	IsUnpublished   bool   `json:"is_unpublished"`
	AccountBadges   []any  `json:"account_badges"`
	LatestReelMedia int    `json:"latest_reel_media"`
	HasChaining     bool   `json:"has_chaining"`
	AccountType     int    `json:"account_type"`

	MutualFollowersCount      int  `json:"mutual_followers_count"`
	IsProfileAudioCallEnabled bool `json:"is_profile_audio_call_enabled"`
}

func (u *UserInfo) HasUserInfoExtendedData() bool {
	return u.UserInfoExtended != nil
}

type ProfilePic struct {
	ProfilePicID               string `json:"profile_pic_id"`
	ProfilePicURL              string `json:"profile_pic_url"`
	HasAnonymousProfilePicture bool   `json:"has_anonymous_profile_picture"`

	*ProfilePicExtended
}

type ProfilePicExtended struct {
	HdProfilePicURLInfo struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"hd_profile_pic_url_info"`
	HdProfilePicVersions []struct {
		Width  int    `json:"width"`
		Height int    `json:"height"`
		URL    string `json:"url"`
	} `json:"hd_profile_pic_versions"`
}

func (p *ProfilePic) HasProfilePicExtended() bool {
	return p.ProfilePicExtended != nil
}

type ProfileContext struct {
	// ProfileContext is the social context displayed at the bottom of a profile
	// E.g. 'Followed by xxx, yyy'.
	ProfileContext              string  `json:"profile_context"`
	MutualFollowIDs             []int64 `json:"profile_context_mutual_follow_ids"`
	ProfileContextFacepileUsers []struct {
		*UserInfo
		*ProfilePic
	} `json:"profile_context_facepile_users"`
	ProfileContextLinksWithUserIds []struct {
		Start    int    `json:"start"`
		End      int    `json:"end"`
		Username string `json:"username,omitempty"`
	} `json:"profile_context_links_with_user_ids"`
}

type Category struct {
	Category           string `json:"category"`
	IsCategoryTappable bool   `json:"is_category_tappable"`

	*CategoryExtended
}

type CategoryExtended struct {
	CategoryID         int64 `json:"category_id"`
	CanHideCategory    bool  `json:"can_hide_category"`
	ShouldShowCategory bool  `json:"should_show_category"`
}

func (c *Category) HasCategoryExtendedData() bool {
	return c.CategoryExtended != nil
}

type Shop struct {
	// TODO: fetch info from a real shop and see if it has extra fields
	ShowShoppableFeed       bool   `json:"show_shoppable_feed"`
	MerchantCheckoutStyle   string `json:"merchant_checkout_style"`
	SellerShoppableFeedType string `json:"seller_shoppable_feed_type"`
}

type FanClub struct {
	FanClubID                            any `json:"fan_club_id"`
	FanClubName                          any `json:"fan_club_name"`
	IsFanClubReferralEligible            any `json:"is_fan_club_referral_eligible"`
	FanConsiderationPageRevampEligiblity any `json:"fan_consideration_page_revamp_eligiblity"`
	IsFanClubGiftingEligible             any `json:"is_fan_club_gifting_eligible"`
}

type UserStats struct {
	FollowerCount     int `json:"follower_count"`
	FollowingCount    int `json:"following_count"`
	FollowingTagCount int `json:"following_tag_count"`
	UsertagsCount     int `json:"usertags_count"`

	MediaCount      int `json:"media_count"`
	TotalClipsCount int `json:"total_clips_count"`
	TotalIGTVVideos int `json:"total_igtv_videos"`

	HasIGTVSeries bool `json:"has_igtv_series"`
	HasGuides     bool `json:"has_guides"`
}

type AdsAccount struct {
	AdsPageID                  int64  `json:"ads_page_id"`
	AdsPageName                string `json:"ads_page_name"`
	IsExperiencedAdvertiser    bool   `json:"is_experienced_advertiser"`
	AdsIncentiveExpirationDate any    `json:"ads_incentive_expiration_date"`
}

type Fundraiser struct {
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
	ActiveStandaloneFundraisers struct {
		TotalCount  int   `json:"total_count"`
		Fundraisers []any `json:"fundraisers"`
	} `json:"active_standalone_fundraisers"`
}

type BiographyInfo struct {
	Biography             string `json:"biography"`
	BiographyWithEntities struct {
		RawText  string `json:"raw_text"`
		Entities []struct {
			User struct {
				ID       int    `json:"id"`
				Username string `json:"username"`
			} `json:"user"`
		} `json:"entities"`
	} `json:"biography_with_entities"`

	ExternalURL            string `json:"external_url"`
	ExternalLynxURL        string `json:"external_lynx_url"`
	PrimaryProfileLinkType int    `json:"primary_profile_link_type"`
	ShowFbLinkOnProfile    bool   `json:"show_fb_link_on_profile"`

	*BiographyExtended
}

type BiographyExtended struct {
	HasBiographyTranslation bool `json:"has_biography_translation"`
}

func (b *BiographyInfo) HasBiographyExtendedData() bool {
	return b.BiographyExtended != nil
}

type ContactInfo struct {
	BusinessContactMethod  string `json:"business_contact_method"`
	PublicEmail            string `json:"public_email"`
	PublicPhoneCountryCode string `json:"public_phone_country_code"`
	PublicPhoneNumber      string `json:"public_phone_number"`
	ContactPhoneNumber     string `json:"contact_phone_number"`
}

type Friendship struct {
	Following       bool `json:"following"`
	OutgoingRequest bool `json:"outgoing_request"`
	IsMutingReel    bool `json:"is_muting_reel"`
	IsBestie        bool `json:"is_bestie"`
	IsRestricted    bool `jsoN:"is_restricted"`
	IsFeedFavorite  bool `json:"is_feed_favorite"`

	*FriendshipExtended
}

type FriendshipExtended struct {
	Blocking              bool `json:"blocking"`
	FollowedBy            bool `json:"followed_by"`
	IncomingRequest       bool `json:"incoming_request"`
	IsBlockingReel        bool `json:"is_blocking_reel"`
	IsPrivate             bool `json:"is_private"`
	Muting                bool `json:"muting"`
	Subscribed            bool `json:"subscribed"`
	IsEligibleToSubscribe bool `json:"is_eligible_to_subscribe"`
	IsSupervisedByViewer  bool `json:"is_supervised_by_viewer"`
	IsGuardianOfViewer    bool `json:"is_guardian_of_viewer"`
}

func (f *Friendship) HasFriendshipExtended() bool {
	return f.FriendshipExtended != nil
}

type Address struct {
	CityID        int     `json:"city_id"`
	CityName      string  `json:"city_name"`
	Zip           string  `json:"zip"`
	AddressStreet string  `json:"address_street"`
	LocationID    string  `json:"instagram_location_id"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
}

type Favorites struct {
	IsFavoriteForStories    bool `json:"is_favorite_for_stories"`
	IsFavoriteForIGTV       bool `json:"is_favorite_for_igtv"`
	IsFavoriteForClips      bool `json:"is_favorite_for_clips"`
	IsFavoriteForHighlights bool `json:"is_favorite_for_highlights"`
}

type Business struct {
	// is_potential_business
	IsEligibleForSmbSupportFlow bool `json:"is_eligible_for_smb_support_flow"`
	IsEligibleForLeadCenter     bool `json:"is_eligible_for_lead_center"`
	SmbDeliveryPartner          any  `json:"smb_delivery_partner"`
	SmbSupportDeliveryPartner   any  `json:"smb_support_delivery_partner"`
	SmbSupportPartner           any  `json:"smb_support_partner"`
	IsCallToActionEnabled       bool `json:"is_call_to_action_enabled"`
}

type ExtraUnknown struct {
	MiniShopSellerOnboardingStatus any `json:"mini_shop_seller_onboarding_status,omitempty"`
	ShoppingPostOnboardNuxType     any `json:"shopping_post_onboard_nux_type,omitempty"`
	NumOfAdminedPages              any `json:"num_of_admined_pages,omitempty"`

	DisplayedActionButtonPartner any    `json:"displayed_action_button_partner"`
	DisplayedActionButtonType    string `json:"displayed_action_button_type"`

	CurrentCatalogID any    `json:"current_catalog_id"`
	DirectMessaging  string `json:"direct_messaging"`

	ProfessionalConversionSuggestedAccountType int `json:"professional_conversion_suggested_account_type"`

	CanHidePublicContacts    bool   `json:"can_hide_public_contacts"`
	ShouldShowPublicContacts bool   `json:"should_show_public_contacts"`
	LeadDetailsAppID         string `json:"lead_details_app_id"`

	LiveSubscriptionStatus string `json:"live_subscription_status"`

	UpcomingEvents []any `json:"upcoming_events"`
}
