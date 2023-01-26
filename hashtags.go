package goinsta

import (
	"encoding/json"
	"errors"
	"fmt"
)

type PageTab string

var (
	// PageTop fetches the top items.
	PageTop PageTab = "top"
	// PageRecent fetches the recent items.
	PageRecent PageTab = "recent"
	// PageReels fetches the reels items.
	PageReels PageTab = "clips"
)

// Hashtag is used for getting the media that matches a hashtag on instagram.
type Hashtag struct {
	insta *Instagram
	err   error

	Name                string `json:"name"`
	ID                  int64  `json:"id"`
	MediaCount          int    `json:"media_count"`
	FormattedMediaCount string `json:"formatted_media_count,omitempty"`
	FollowStatus        any    `json:"follow_status,omitempty"`
	Subtitle            string `json:"subtitle,omitempty"`
	Description         string `json:"description,omitempty"`
	Following           any    `json:"following,omitempty"`
	AllowFollowing      any    `json:"allow_following,omitempty"`
	AllowMutingStory    any    `json:"allow_muting_story,omitempty"`
	ProfilePicURL       any    `json:"profile_pic_url,omitempty"`
	NonViolating        any    `json:"non_violating,omitempty"`
	RelatedTags         any    `json:"related_tags,omitempty"`
	DebugInfo           any    `json:"debug_info,omitempty"`
	// All Top Items
	Items      []*Item     `json:"items,omitempty"`
	Story      *StoryMedia `json:"story,omitempty"`
	NumResults int

	// Sections will always contain the last fetched sections, regardless of tab
	Sections            []hashtagSection `json:"sections,omitempty"`
	PageInfo            map[string]hashtagPageInfo
	AutoLoadMoreEnabled bool    `json:"auto_load_more_enabled,omitempty"`
	MoreAvailable       bool    `json:"more_available,omitempty"`
	NextID              string  `json:"next_max_id,omitempty"`
	NextPage            int     `json:"next_page,omitempty"`
	NextMediaIds        []int64 `json:"next_media_ids,omitempty"`
	Status              string  `json:"status,omitempty"`
}

type hashtagSection struct {
	LayoutType    string `json:"layout_type"`
	LayoutContent struct {
		// TODO: misses onebyoneitem etc., like discover page
		FillItems []mediaItem `json:"fill_items"`
		Medias    []mediaItem `json:"medias"`
	} `json:"layout_content"`
	FeedType        string `json:"feed_type"`
	ExploreItemInfo struct {
		NumColumns      int     `json:"num_columns"`
		TotalNumColumns int     `json:"total_num_columns"`
		AspectRatio     float32 `json:"aspect_ratio"`
		Autoplay        bool    `json:"autoplay"`
	} `json:"explore_item_info"`
}

type mediaItem struct {
	Item *Item `json:"media"`
}

type hashtagPageInfo struct {
	MoreAvailable bool    `json:"more_available"`
	NextID        string  `json:"next_max_id"`
	NextPage      int     `json:"next_page"`
	NextMediaIds  []int64 `json:"next_media_ids"`
	Status        string  `json:"status"`
}

func (h *Hashtag) setValues() {
	if h.PageInfo == nil {
		h.PageInfo = make(map[string]hashtagPageInfo)
	}

	for _, s := range h.Sections {
		for _, m := range s.LayoutContent.Medias {
			setToItem(m.Item, h)
		}

		for _, m := range s.LayoutContent.FillItems {
			setToItem(m.Item, h)
		}
	}
}

// Delete only a place holder, does nothing.
func (h *Hashtag) Delete() error {
	return nil
}

func (h *Hashtag) GetNextID() string {
	return ""
}

// NewHashtag returns initialized hashtag structure.
// Name parameter is hashtag name.
func (insta *Instagram) NewHashtag(name string) *Hashtag {
	return &Hashtag{
		insta:    insta,
		Name:     name,
		PageInfo: make(map[string]hashtagPageInfo),
	}
}

// Sync wraps Hashtag.Info().
func (h *Hashtag) Sync() error {
	return h.Info()
}

// Info updates Hashtag information.
func (h *Hashtag) Info() error {
	insta := h.insta

	body, err := insta.sendSimpleRequest(urlTagInfo, h.Name)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, h)
}

// Next paginates over hashtag top pages.
func (h *Hashtag) Next(tab ...interface{}) bool {
	p := PageTop

	if len(tab) > 0 {
		page, ok := tab[0].(PageTab)
		if !ok {
			h.err = errors.New("can only provide type PageTab")
			return false
		}

		if len(page) > 0 {
			p = page
		}
	}

	return h.next(p)
}

func (h *Hashtag) next(tab PageTab) bool {
	pageInfo, ok := h.PageInfo[string(tab)]

	if h.err != nil && errors.Is(h.err, ErrNoMore) {
		return false
	} else if errors.Is(h.err, ErrNoMore) && ok && !pageInfo.MoreAvailable {
		return false
	}

	if tab != "top" && tab != "recent" && tab != "clips" {
		h.err = ErrInvalidTab
		return false
	}

	insta := h.insta
	name := h.Name

	query := map[string]string{
		"tab":                string(tab),
		"_uuid":              insta.uuid,
		"include_persistent": "false",
		"rank_token":         insta.rankToken,
	}

	if ok {
		nextMediaIds, err := json.Marshal(pageInfo.NextMediaIds)
		if err != nil {
			h.err = err
			return false
		}

		query["max_id"] = pageInfo.NextID
		query["page"] = toString(pageInfo.NextPage)
		query["next_media_ids"] = string(nextMediaIds)
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlTagContent, name),
			IsPost:   true,
			Query:    query,
		},
	)
	if err != nil {
		h.err = err
		return false
	}

	res := &Hashtag{}
	if err := json.Unmarshal(body, res); err != nil {
		h.err = err
		return false
	}

	h.fillItems(res, tab)

	if !h.MoreAvailable {
		h.err = ErrNoMore
		return false
	}

	return true
}

func (h *Hashtag) fillItems(res *Hashtag, tab PageTab) {
	h.AutoLoadMoreEnabled = res.AutoLoadMoreEnabled
	h.NextID = res.NextID
	h.MoreAvailable = res.MoreAvailable
	h.NextPage = res.NextPage
	h.Status = res.Status
	h.Sections = res.Sections

	h.PageInfo[string(tab)] = hashtagPageInfo{
		MoreAvailable: res.MoreAvailable,
		NextID:        res.NextID,
		NextPage:      res.NextPage,
		NextMediaIds:  res.NextMediaIds,
	}

	h.setValues()

	count := 0

	for _, s := range res.Sections {
		for _, m := range s.LayoutContent.Medias {
			count++

			h.Items = append(h.Items, m.Item)
		}

		for _, m := range s.LayoutContent.FillItems {
			count++

			h.Items = append(h.Items, m.Item)
		}
	}

	h.NumResults = count
}

// Latest will return the last fetched items.
func (h *Hashtag) Latest() []*Item {
	var res []*Item

	for _, s := range h.Sections {
		for _, m := range s.LayoutContent.Medias {
			res = append(res, m.Item)
		}
	}

	return res
}

// Error returns hashtag error.
func (h *Hashtag) Error() error {
	return h.err
}

// Clears the Hashtag.err error.
func (h *Hashtag) ClearError() {
	h.err = nil
}

func (h *Hashtag) getInsta() *Instagram {
	return h.insta
}

// Stories returns hashtag stories.
func (h *Hashtag) Stories() error {
	body, err := h.insta.sendSimpleRequest(
		urlTagStories, h.Name,
	)
	if err != nil {
		return err
	}

	var resp struct {
		Story  *StoryMedia `json:"story"`
		Status string      `json:"status"`
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return err
	}

	h.Story = resp.Story

	return err
}
