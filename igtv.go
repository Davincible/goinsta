package goinsta

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Do i need to extract the rank token?

// Methods to create:
// Broadcasts.Discover

// All items with interface{} I have only seen a null response
type IGTV struct {
	insta *Instagram
	err   error

	// Shared between the endpoints
	DestinationClientConfigs interface{} `json:"destination_client_configs"`
	MaxID                    string      `json:"max_id"`
	MoreAvailable            bool        `json:"more_available"`
	SeenState                interface{} `json:"seen_state"`
	NumResults               int         `json:"num_results"`
	Status                   string      `json:"status"`

	// Specific to igtv/discover
	Badging          interface{}   `json:"badging"`
	BannerToken      interface{}   `json:"banner_token"`
	BrowseItems      interface{}   `json:"browser_items"`
	Channels         []IGTVChannel `json:"channels"`
	Composer         interface{}   `json:"composer"`
	Items            []*Item       `json:"items"`
	DestinationItems []IGTVItem    `json:"destination_items"`
	MyChannel        struct{}      `json:"my_channel"`

	// Specific to igtv/suggested_searches
	RankToken int `json:"rank_token"`
}

// IGTVItem is a media item that can be found inside the IGTV struct, from the
//   IGTV Discover endpoint.
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

// IGTVChannel can represent a single user's collection of IGTV posts, or it can
//   e.g. represent a user's IGTV series.
//
// It's called a channel, however the Items inside the Channel struct can, but
//   don't have to, belong to the same account, depending on the request. It's a bit dubious
//
type IGTVChannel struct {
	insta *Instagram
	id    string // user id parameter
	err   error

	ApproxTotalVideos        interface{}  `json:"approx_total_videos"`
	ApproxVideosFormatted    interface{}  `json:"approx_videos_formatted"`
	CoverPhotoUrl            string       `json:"cover_photo_url"`
	Description              string       `json:"description"`
	ID                       string       `json:"id"`
	Items                    []*Item      `json:"items"`
	NumResults               int          `json:"num_results"`
	Broadcasts               []*Broadcast `json:"live_items"`
	Title                    string       `json:"title"`
	Type                     string       `json:"type"`
	User                     *User        `json:"user_dict"`
	DestinationClientConfigs interface{}  `json:"destination_client_configs"`
	NextID                   interface{}  `json:"max_id"`
	MoreAvailable            bool         `json:"more_available"`
	SeenState                interface{}  `json:"seen_state"`
}

func newIGTV(insta *Instagram) *IGTV {
	return &IGTV{insta: insta}
}

// IGTV returns the IGTV items of a user
//
// Use IGTVChannel.Next for pagination.
//
func (user *User) IGTV() (*IGTVChannel, error) {
	insta := user.insta
	igtv := &IGTVChannel{
		insta: insta,
		id:    fmt.Sprintf("user_%d", user.ID),
	}
	if !igtv.Next() {
		return igtv, igtv.Error()
	}
	return igtv, nil
}

// IGTVSeries will fetch the igtv series of a user. Usually the slice length
//   of the return value is 1, as there is one channel, which contains multiple
//   series.
//
func (user *User) IGTVSeries() ([]*IGTVChannel, error) {
	if !user.HasIGTVSeries {
		return nil, ErrIGTVNoSeries
	}
	insta := user.insta

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlIGTVSeries, user.ID),
			Query:    generateSignature("{}"),
		},
	)
	if err != nil {
		return nil, err
	}

	var res struct {
		Channels []*IGTVChannel `json:"channels"`
		Status   string         `json:"status"`
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	for _, chann := range res.Channels {
		chann.setValues()
	}
	return res.Channels, nil
}

// Live will return a list of current broadcasts
func (igtv *IGTV) Live() (*IGTVChannel, error) {
	return igtv.insta.callIGTVChannel("live", "")
}

// Live test method to see if Live can paginate
func (igtv *IGTVChannel) Live() (*IGTVChannel, error) {
	return igtv.insta.callIGTVChannel("live", igtv.GetNextID())
}

// GetNexID returns the max id used for pagination.
func (igtv *IGTVChannel) GetNextID() string {
	return formatID(igtv.NextID)
}

// Next allows you to paginate the IGTV feed of a channel.
// returns false when list reach the end.
// if FeedMedia.Error() is ErrNoMore no problems have occurred.
func (igtv *IGTVChannel) Next(params ...interface{}) bool {
	if igtv.err != nil {
		return false
	}
	insta := igtv.insta

	query := map[string]string{
		"id":    igtv.id,
		"_uuid": igtv.insta.uuid,
		"count": "10",
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
			Query:    query,
		})
	if err != nil {
		igtv.err = err
		return false
	}

	oldItems := igtv.Items
	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber()
	err = d.Decode(igtv)
	if err != nil {
		igtv.err = err
		return false
	}

	if !igtv.MoreAvailable {
		igtv.err = ErrNoMore
	}
	igtv.setValues()
	igtv.Items = append(oldItems, igtv.Items...)
	return true
}

// Latest will return the results from the latest fetch
func (igtv *IGTVChannel) Latest() []*Item {
	return igtv.Items[len(igtv.Items)-igtv.NumResults:]
}

func (insta *Instagram) callIGTVChannel(id, nextID string) (*IGTVChannel, error) {
	query := map[string]string{
		"id":    id,
		"_uuid": insta.uuid,
		"count": "10",
	}
	if nextID != "" {
		query["max_id"] = nextID
	}
	body, _, err := insta.sendRequest(&reqOptions{
		Endpoint: urlIGTVChannel,
		IsPost:   true,
		Query:    query,
	})
	if err != nil {
		return nil, err
	}
	igtv := &IGTVChannel{insta: insta, id: id}
	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber()
	err = d.Decode(igtv)
	igtv.setValues()
	return igtv, err
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
	if igtv.User != nil {
		igtv.User.insta = insta
	}
	for _, i := range igtv.Items {
		setToItem(i, igtv)
	}
	for _, br := range igtv.Broadcasts {
		br.setValues(insta)
	}
}

// Next allows you to paginate the IGTV Discover page.
func (igtv *IGTV) Next(params ...interface{}) bool {
	if igtv.err != nil {
		return false
	}
	insta := igtv.insta

	query := map[string]string{}
	if igtv.MaxID != "" {
		query["max_id"] = igtv.MaxID
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlIGTVDiscover,
			Query:    query,
		},
	)
	if err != nil {
		igtv.err = err
		return false
	}

	err = json.Unmarshal(body, igtv)
	if err != nil {
		igtv.err = err
		return false
	}
	igtv.setValues()
	count := 0
	for _, item := range igtv.DestinationItems {
		if item.Item != nil {
			igtv.Items = append(igtv.Items, item.Item)
			count++
		}
	}
	igtv.NumResults = count
	if !igtv.MoreAvailable {
		igtv.err = ErrNoMore
		return false
	}
	return true
}

func (igtv *IGTV) setValues() {
	for _, item := range igtv.DestinationItems {
		if item.Item != nil {
			setToItem(item.Item, igtv)
		}
	}
	for _, chann := range igtv.Channels {
		chann.setValues()
	}
}

// Delete does nothing, is only a place holder
func (igtv *IGTV) Delete() error {
	return nil
}

// Error return the error of IGTV, if one has occured
func (igtv *IGTV) Error() error {
	return igtv.err
}

func (igtv *IGTV) getInsta() *Instagram {
	return igtv.insta
}

// Error return the error of IGTV, if one has occured
func (igtv *IGTV) GetNextID() string {
	return igtv.MaxID
}

// Latest returns the last fetched items, by slicing IGTV.Items with IGTV.NumResults
func (igtv *IGTV) Latest() []*Item {
	return igtv.Items[len(igtv.Items)-igtv.NumResults:]
}
