package goinsta

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

// Search is the object for all searches like Facebook, Location or Tag search.
type Search struct {
	insta *Instagram
}

type SearchFunc interface{}

// SearchResult handles the data for the results given by each type of Search.
type SearchResult struct {
	insta *Instagram
	err   error

	HasMore       bool   `json:"has_more"`
	PageToken     string `json:"page_token"`
	RankToken     string `json:"rank_token"`
	Status        string `json:"status"`
	NumResults    int64  `json:"num_results"`
	Query         string
	SearchSurface string
	context       string
	queryParam    string
	entityType    string

	// Regular Search Results
	Results []TopSearchItem `json:"list"`
	History []SearchHistory

	// User search results
	Users []*User `json:"users"`

	// Loaction search results
	Places []Place `json:"items"`

	// Tag search results
	Tags []struct {
		ID               int64       `json:"id"`
		Name             string      `json:"name"`
		MediaCount       int         `json:"media_count"`
		FollowStatus     interface{} `json:"follow_status"`
		Following        interface{} `json:"following"`
		AllowFollowing   interface{} `json:"allow_following"`
		AllowMutingStory interface{} `json:"allow_muting_story"`
		ProfilePicURL    interface{} `json:"profile_pic_url"`
		NonViolating     interface{} `json:"non_violating"`
		RelatedTags      interface{} `json:"related_tags"`
		DebugInfo        interface{} `json:"debug_info"`
	} `json:"results"`

	// Location search result
	RequestID string `json:"request_id"`
	Venues    []struct {
		ExternalIDSource string  `json:"external_id_source"`
		ExternalID       string  `json:"external_id"`
		Lat              float64 `json:"lat"`
		Lng              float64 `json:"lng"`
		Address          string  `json:"address"`
		Name             string  `json:"name"`
	} `json:"venues"`

	ClearClientCache bool `json:"clear_client_cache"`
}

type Place struct {
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle"`
	Location Location `json:"location"`
}

type TopSearchItem struct {
	insta *Instagram

	Position int     `json:"position"`
	User     User    `json:"user"`
	Hashtag  Hashtag `json:"hashtag"`
	Place    Place   `json:"place"`
}

type SearchHistory struct {
	Time int64 `json:"client_time"`
	User User  `json:"user"`
}

// newSearch creates new Search structure
func newSearch(insta *Instagram) *Search {
	search := &Search{
		insta: insta,
	}
	return search
}

// Search is a wrapper for insta.Searchbar.Search()
func (insta *Instagram) Search(query string) (*SearchResult, error) {
	return insta.Searchbar.Search(query)
}

func (sb *Search) Search(query string) (*SearchResult, error) {
	return sb.search(query, sb.topsearch)
}

func (sb *Search) SearchUser(query string) (*SearchResult, error) {
	return sb.search(query, sb.user)
}

func (sb *Search) SearchHashtag(query string) (*SearchResult, error) {
	return sb.search(query, sb.tags)
}

func (sb *Search) SearchLocation(query string) (*SearchResult, error) {
	return sb.search(query, sb.places)
}

func (sr *SearchResult) Next() bool {
	if !sr.HasMore || sr.RankToken == "" || sr.PageToken == "" {
		sr.err = errors.New("No more results available, or rank or page token have not been set")
		return false
	}

	insta := sr.insta
	query := map[string]string{
		"search_surface":  sr.SearchSurface,
		"timezone_offset": timeOffset,
		"count":           "30",
		sr.queryParam:     sr.Query,
		"rank_token":      sr.RankToken,
		"page_token":      sr.PageToken,
	}
	if sr.context != "" {
		query["context"] = sr.context
	}
	if sr.SearchSurface == "places_search_page" {
		query["lon"] = ""
		query["lng"] = ""
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlSearchTop,
			Query:    query,
		},
	)
	if err != nil {
		sr.err = err
		return false
	}
	res := &SearchResult{}
	err = json.Unmarshal(body, res)
	if err != nil {
		sr.err = err
		return false
	}
	res.setValues()
	sr.Results = append(sr.Results, res.Results...)
	sr.HasMore = res.HasMore
	sr.RankToken = res.RankToken
	sr.ClearClientCache = res.ClearClientCache
	return true
}

func (sb *Search) History() (*[]SearchHistory, error) {
	sb.insta.Discover.Next()
	h, err := sb.history()
	if err != nil {
		return nil, err
	}
	if err := sb.NullState(); err != nil {
		sb.insta.ErrHandler("Non fatal error while setting search null state", err)
	}
	return h, nil
}

func (sr *TopSearchItem) RegisterClick() error {
	var entityType string
	var id int64
	if id = sr.User.ID; id != 0 {
		entityType = "user"
	} else if id = sr.Hashtag.ID; id != 0 {
		entityType = "hashtag"
	} else if id = sr.Place.Location.Pk; id != 0 {
		entityType = "place"
	}

	_, _, err := sr.insta.sendRequest(&reqOptions{
		Endpoint: urlSearchRegisterClick,
		IsPost:   true,
		Query: map[string]string{
			"entity_id":   strconv.Itoa(int(id)),
			"_uuid":       sr.insta.uuid,
			"entity_type": entityType,
		},
	})
	return err
}

func (sb *Search) search(query string, fn func(string) (*SearchResult, error)) (*SearchResult, error) {
	sb.insta.Discover.Next()
	h, err := sb.history()
	if err != nil {
		sb.insta.ErrHandler("Non fatal error while fetcihng recent search results",
			err)
	}
	if err := sb.NullState(); err != nil {
		sb.insta.ErrHandler("Non fatal error while setting search null state", err)
	}

	var result *SearchResult
	var q string
	for _, char := range query {
		q += string(char)
		result, err = fn(q)
		if err != nil {
			return nil, err
		}
		s := random(150, 500)
		time.Sleep(time.Duration(s) * time.Millisecond)
	}
	result.History = *h
	return result, nil
}

func (search *Search) topsearch(query string) (*SearchResult, error) {
	insta := search.insta
	res := &SearchResult{
		Query:         query,
		SearchSurface: "top_search_page",
		context:       "blended",
		queryParam:    "query",
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlSearchTop,
			Query: map[string]string{
				"search_surface":  res.SearchSurface,
				"timezone_offset": timeOffset,
				"count":           "30",
				res.queryParam:    query,
				"context":         res.context,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, res)
	res.setValues()
	return res, err
}

func (sr *SearchResult) Error() error {
	return sr.err
}

func (sr *SearchResult) setValues() {
	for _, r := range sr.Results {
		r.insta = sr.insta
		r.User.insta = sr.insta
		r.Hashtag.insta = sr.insta
	}
	for _, u := range sr.Users {
		u.insta = sr.insta
	}
}

func (search *Search) user(user string) (*SearchResult, error) {
	insta := search.insta
	res := &SearchResult{
		insta:         insta,
		SearchSurface: "user_search_page",
		queryParam:    "q",
		Query:         user,
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlSearchUser,
			Query: map[string]string{
				"search_surface":  res.SearchSurface,
				"timezone_offset": timeOffset,
				"count":           "30",
				res.queryParam:    user,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, res)
	res.setValues()
	return res, err
}

func (search *Search) tags(tag string) (*SearchResult, error) {
	insta := search.insta
	res := &SearchResult{
		SearchSurface: "hashtag_search_page",
		queryParam:    "q",
		Query:         tag,
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlSearchTag,
			Query: map[string]string{
				"search_surface":  res.SearchSurface,
				"timezone_offset": timeOffset,
				"count":           "30",
				res.queryParam:    tag,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, res)
	return res, err
}

func (search *Search) places(location string) (*SearchResult, error) {
	insta := search.insta
	res := &SearchResult{
		SearchSurface: "places_search_page",
		queryParam:    "query",
		Query:         location,
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlSearchLocation,
			Query: map[string]string{
				"search_surface":  res.SearchSurface,
				"timezone_offset": timeOffset,
				"count":           "30",
				res.queryParam:    location,
				"lat":             "",
				"lng":             "",
			},
		},
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, res)
	return res, err
}

func (search *Search) NullState() error {
	_, _, err := search.insta.sendRequest(&reqOptions{
		Endpoint: urlSearchNullState,
		Query:    map[string]string{"type": "blended"},
	})
	return err
}

func (search *Search) history() (*[]SearchHistory, error) {
	body, err := search.insta.sendSimpleRequest(urlSearchRecent)
	if err != nil {
		return nil, err
	}
	s := struct {
		Recent []SearchHistory `json:"recent"`
		Status string          `json:"status"`
	}{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	for _, i := range s.Recent {
		i.User.insta = search.insta
	}
	return &s.Recent, nil
}
