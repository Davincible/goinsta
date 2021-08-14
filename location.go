package goinsta

import (
	"encoding/json"
	"fmt"
)

type LocationInstance struct {
	insta *Instagram
}

func newLocation(insta *Instagram) *LocationInstance {
	return &LocationInstance{insta: insta}
}

type LayoutSection struct {
	LayoutType    string `json:"layout_type"`
	LayoutContent struct {
		Medias []struct {
			Media Item `json:"media"`
		} `json:"medias"`
	} `json:"layout_content"`
	FeedType        string `json:"feed_type"`
	ExploreItemInfo struct {
		NumColumns      int     `json:"num_columns"`
		TotalNumColumns int     `json:"total_num_columns"`
		AspectRatio     float64 `json:"aspect_ratio"`
		Autoplay        bool    `json:"autoplay"`
	} `json:"explore_item_info"`
}

type Section struct {
	Sections      []LayoutSection `json:"sections"`
	MoreAvailable bool            `json:"more_available"`
	NextPage      int             `json:"next_page"`
	NextMediaIds  []int64         `json:"next_media_ids"`
	NextID        string          `json:"next_max_id"`
	Status        string          `json:"status"`
}

func (l *LocationInstance) Feeds(locationID int64) (*Section, error) {
	// TODO: use pagination for location feeds.
	insta := l.insta
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlFeedLocations, locationID),
			IsPost:   true,
			Query: map[string]string{
				"rank_token":     insta.rankToken,
				"ranked_content": "true",
				"_uid":           toString(insta.Account.ID),
				"_uuid":          insta.uuid,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	section := &Section{}
	err = json.Unmarshal(body, section)
	return section, err
}

func (l *Location) Feed() (*Section, error) {
	// TODO: use pagination for location feeds.
	insta := l.insta
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlFeedLocations, l.ID),
			IsPost:   true,
			Query: map[string]string{
				"rank_token":     insta.rankToken,
				"ranked_content": "true",
				"_uid":           toString(insta.Account.ID),
				"_uuid":          insta.uuid,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	section := &Section{}
	err = json.Unmarshal(body, section)
	return section, err
}
