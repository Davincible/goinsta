package goinsta

import (
	"bytes"
	"encoding/json"
)

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
	Item       *Item       `json:"item"`
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
	insta *Instagram
	id    string // user id parameter
	err   error

	ApproxTotalVideos        interface{} `json:"approx_total_videos"`
	ApproxVideosFormatted    interface{} `json:"approx_videos_formatted"`
	CoverPhotoUrl            string      `json:"cover_photo_url"`
	Description              string      `json:"description"`
	ID                       string      `json:"id"`
	Items                    []*Item     `json:"items"`
	LiveItems                []Broadcast `json:"live_items"`
	Title                    string      `json:"title"`
	Type                     string      `json:"user"`
	User                     *User       `json:"user_dict"`
	DestinationClientConfigs interface{} `json:"destination_client_configs"`
	NextID                   interface{} `json:"max_id"`
	MoreAvailable            bool        `json:"more_available"`
	SeenState                interface{} `json:"seen_state"`
}

// GetNexID returns the max id used for pagination.
func (igtv *IGTVChannel) GetNextID() string {
	return formatID(igtv.NextID)
}

// Next allows pagination after calling:
// User.Feed
// extra query arguments can be passes one after another as func(key, value).
// Only if an even number of string arguements will be passed, they will be
//   used in the query.
// returns false when list reach the end.
// if FeedMedia.Error() is ErrNoMore no problems have occurred.
func (igtv *IGTVChannel) Next(params ...interface{}) bool {
	if igtv.err != nil {
		return false
	}

	insta := igtv.insta

	query := map[string]string{
		"exclude_comment":                 "true",
		"only_fetch_first_carousel_media": "false",
	}
	if len(params)%2 == 0 {
		for i := 0; i < len(params); i = i + 2 {
			query[params[i].(string)] = params[i+1].(string)
		}
	}

	if next := igtv.GetNextID(); next != "" {
		query["max_id"] = next
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlIGTVChannel,
			IsPost:   true,
			Query: map[string]string{
				"id":    igtv.id,
				"_uuid": igtv.insta.uuid,
				"count": "10",
			},
		})
	if err == nil {
		m := IGTVChannel{insta: igtv.insta, id: igtv.id}
		d := json.NewDecoder(bytes.NewReader(body))
		d.UseNumber()
		err = d.Decode(&igtv)
		if err == nil {
			igtv.NextID = m.NextID
			igtv.MoreAvailable = m.MoreAvailable
			if m.NextID == 0 || !m.MoreAvailable {
				igtv.err = ErrNoMore
			}
			m.setValues()
			igtv.Items = append(igtv.Items, m.Items...)
			return true
		}
	}
	return false
}

func (igtv *IGTVChannel) Delete() error {
	return nil
}

func (igtv *IGTVChannel) Error() error {
	return igtv.err
}

func (media *IGTVChannel) getInsta() *Instagram {
	return media.insta
}

func (igtv *IGTVChannel) setValues() {
	insta := igtv.insta
	igtv.User.insta = insta
	for _, i := range igtv.Items {
		setToItem(i, igtv)
	}
}
