package goinsta

import "encoding/json"

type Discover struct {
	insta      *Instagram
	sessionId  string
	err        error
	Items      []DiscoverSectionalItem
	NumResults int

	AutoLoadMoreEnabled bool `json:"auto_load_more_enabled"`
	Clusters            []struct {
		CanMute     bool        `json:"can_mute"`
		Context     string      `json:"context"`
		DebugInfo   string      `json:"debug_info"`
		Description string      `json:"description"`
		ID          interface{} `json:"id"`
		IsMuted     bool        `json:"is_muted"`
		Labels      interface{} `json:"labels"`
		Name        string      `json:"name"`
		Title       string      `json:"title"`
		Type        string      `json:"type"`
	} `json:"clusters"`
	MaxID              string                  `json:"max_id"`
	MoreAvailable      bool                    `json:"more_available"`
	NextID             string                  `json:"next_max_id"`
	RankToken          string                  `json:"rank_token"`
	SectionalItems     []DiscoverSectionalItem `json:"sectional_items"`
	SessionPagingToken string                  `json:"session_paging_token"`
	Status             string                  `json:"status"`
}

type DiscoverMediaItem struct {
	Media Item `json:"media"`
}

type DiscoverSectionalItem struct {
	ExploreItemInfo struct {
		AspectRatio     float64 `json:"aspect_ratio"`
		Autoplay        bool    `json:"autoplay"`
		NumColumns      int     `json:"num_columns"`
		TotalNumColumns int     `json:"total_num_columns"`
	} `json:"explore_item_info"`
	FeedType      string `json:"feed_type"`
	LayoutContent struct {
		// Usually not all of these are filled
		// I have often seen the 1 out of 5 items being the ThreeByFour
		//   and the third item the TwoByTwoItems + Fill Items,
		//   the others tend to be Medias, but this is not always the case
		LayoutType string              `json:"layout_type"`
		Medias     []DiscoverMediaItem `json:"medias"`

		FillItems    []DiscoverMediaItem `json:"fill_items"`
		OneByOneItem DiscoverMediaItem   `json:"one_by_one"`
		TwoByTwoItem DiscoverMediaItem   `json:"two_by_two_item"`

		ThreeByFourItem struct {
			// TODO: this is a reels section, which you can paginate on its own
			Clips struct {
				ContentSource string              `json:"content_source"`
				Design        string              `json:"design"`
				ID            string              `json:"id"`
				Items         []DiscoverMediaItem `json:"items"`
				Label         string              `json:"label"`
				MaxID         string              `json:"max_id"`
				MoreAvailable bool                `json:"more_available"`
				Type          string              `json:"type"`
			} `json:"clips"`
		} `json:"three_by_four_item"`
	} `json:"layout_content"`
	LayoutType string `json:"layout_type"`
}

func newDiscover(insta *Instagram) *Discover {
	return &Discover{
		insta:     insta,
		sessionId: generateUUID(),
	}
}

// Next allows you to paginate explore page results. Also use this for your
//  first fetch
func (disc *Discover) Next() bool {
	if disc.sessionId == "" {
		disc.sessionId = generateUUID()
	}

	query := map[string]string{
		"omit_cover_media":           "true",
		"reels_configuration":        "default",
		"use_sectional_payload":      "true",
		"timezone_offset":            timeOffset,
		"session_id":                 disc.sessionId,
		"include_fixed_destinations": "true",
	}

	if disc.NextID != "" {
		query["max_id"] = disc.NextID
		query["is_prefetch"] = "false"
	} else {
		query["is_prefetch"] = "true"
	}

	body, _, err := disc.insta.sendRequest(&reqOptions{
		Endpoint: urlDiscoverExplore,
		Query:    query,
	})
	if err != nil {
		disc.err = err
		return false
	}

	err = json.Unmarshal(body, disc)
	if err != nil {
		disc.err = err
		return false
	}
	disc.setValues()
	disc.Items = append(disc.Items, disc.SectionalItems...)
	disc.NumResults = len(disc.SectionalItems)
	return true
}

// Error will return the error, if one is present
func (disc *Discover) Error() error {
	return disc.err
}

// Refresh will remove the session token, and frefresh the results, like a pull down
func (disc *Discover) Refresh() bool {
	disc.sessionId = generateUUID()
	disc.NextID = ""
	return disc.Next()
}

func (disc *Discover) setValues() {
	for _, sec := range disc.SectionalItems {
		for _, i := range sec.LayoutContent.Medias {
			i.Media.insta = disc.insta
			i.Media.User.insta = disc.insta
		}
		for _, i := range sec.LayoutContent.FillItems {
			i.Media.insta = disc.insta
			i.Media.User.insta = disc.insta
		}
		for _, i := range sec.LayoutContent.ThreeByFourItem.Clips.Items {
			i.Media.insta = disc.insta
			i.Media.User.insta = disc.insta
		}
		sec.LayoutContent.OneByOneItem.Media.insta = disc.insta
		sec.LayoutContent.OneByOneItem.Media.User.insta = disc.insta
		sec.LayoutContent.TwoByTwoItem.Media.insta = disc.insta
		sec.LayoutContent.TwoByTwoItem.Media.User.insta = disc.insta
	}
}
