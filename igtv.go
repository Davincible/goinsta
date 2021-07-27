package goinsta

// Do i need to extract the rank token?

// Methods to create:
// User.GetIGTV
// Broadcasts.Discover

// All items with interface{} I have only seen a null response

type IGTV struct {
	// Shared between the endpoints
	DestinationClientConfigs interface{} `json:"destination_client_configs"`
	MaxID                    string      `json:"max_id"`
	MoreAvailable            bool        `json:"more_available"`
	SeenState                interface{} `json:"seen_state"`
	Status                   string      `json:"status"`

	// Specific to igtv/discover
	Badging          interface{}   `json:"badging"`
	BannerToken      interface{}   `json:"banner_token"`
	BrowseItems      interface{}   `json:"browser_items"`
	Channels         []IGTVChannel `json:"channels"`
	Composer         interface{}   `json:"composer"`
	DestinationItems []IGTVItem    `json:"destination_items"`
	MyChannel        struct{}      `json:"my_channel"`

	// Specific to igtv/suggested_searches
	NumResults int `json:"num_results"`
	RankToken  int `json:"rank_token"`
}

type IGTVItem struct {
	Title      string      `json:"title"`
	Type       string      `json:"type"`
	Channel    IGTVChannel `json:"channel"`
	Item       Item        `json:"item"`
	LogingInfo struct {
		SourceChannelType string `json:"source_channel_type"`
	} `json:"logging_info"`

	// Specific to igtv/suggested_searches
	Hashtag interface{} `json:"hashtag"`
	Keyword interface{} `json:"keyword"`
	User    User        `json:"user"`
}

// It's called a channel, however the Items can, but don't have to,
//   belong to the same account, depending on the request. It's a bit dubious
type IGTVChannel struct {
	ApproxTotalVideos        interface{} `json:"approx_total_videos"`
	ApproxVideosFormatted    interface{} `json:"approx_videos_formatted"`
	CoverPhotoUrl            string      `json:"cover_photo_url"`
	Description              string      `json:"description"`
	ID                       string      `json:"id"`
	Items                    []Item      `json:"items"`
	LiveItems                []Broadcast `json:"live_items"`
	Title                    string      `json:"title"`
	Type                     string      `json:"user"`
	User                     User        `json:"user_dict"`
	DestinationClientConfigs interface{} `json:"destination_client_configs"`
	MaxID                    string      `json:"max_id"`
	MoreAvailable            bool        `json:"more_available"`
	SeenState                interface{} `json:"seen_state"`
}
