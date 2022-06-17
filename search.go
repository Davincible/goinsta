package goinsta

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
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
	Results []*TopSearchItem `json:"list"`
	History []SearchHistory

	// User search results
	Users []*User `json:"users"`

	// Loaction search results
	Places []Place `json:"items"`

	// Tag search results
	Tags []*Hashtag `json:"results"`

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
	Title    string    `json:"title"`
	Subtitle string    `json:"subtitle"`
	Location *Location `json:"location"`
}

type TopSearchItem struct {
	insta *Instagram

	Position int      `json:"position"`
	User     *User    `json:"user"`
	Hashtag  *Hashtag `json:"hashtag"`
	Place    Place    `json:"place"`
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
// Search will perform a topsearch query returning users, locations and tags,
//  just like the app would.
//
// By default search behavior will be mimicked by sending a search request per
//  added letter, and waiting a few millis in between, just as if you were to
//  type anything into the search bar. However, if you only want to make one
//  search request passing in the full query immediately, you can turn on quick
//  search by passing in one bool:true parameter, like so:
//
//  Search("myquery", true)  // this will perform a quick search
func (insta *Instagram) Search(query string, p ...bool) (*SearchResult, error) {
	return insta.Searchbar.Search(query, p...)
}

// Search will perform a topsearch query returning users, locations and tags,
//  just like the app would.
//
// By default search behavior will be mimicked by sending a search request per
//  added letter, and waiting a few millis in between, just as if you were to
//  type anything into the search bar. However, if you only want to make one
//  search request passing in the full query immediately, you can turn on quick
//  search by passing in one bool:true parameter, like so:
//
//  Search("myquery", true)  // this will perform a quick search
func (sb *Search) Search(query string, p ...bool) (*SearchResult, error) {
	if isQuickSearch(p) {
		return sb.topsearch(query)
	}
	return sb.search(query, sb.topsearch)
}

// SearchUser will perorm a user search with the provided query.
//
// By default search behavior will be mimicked by sending a search request per
//  added letter, and waiting a few millis in between, just as if you were to
//  type anything into the search bar. However, if you only want to make one
//  search request passing in the full query immediately, you can turn on quick
//  search by passing in one bool:true parameter, like so:
//
//  SearchUser("myquery", true)  // this will perform a quick search
func (sb *Search) SearchUser(query string, p ...bool) (*SearchResult, error) {
	if isQuickSearch(p) {
		return sb.user(query)
	}
	return sb.search(query, sb.user)
}

// SearchHashtag will perform a hashtag search with the provided query.
//
// By default search behavior will be mimicked by sending a search request per
//  added letter, and waiting a few millis in between, just as if you were to
//  type anything into the search bar. However, if you only want to make one
//  search request passing in the full query immediately, you can turn on quick
//  search by passing in one bool:true parameter, like so:
//
//  SearchHashtag("myquery", true)  // this will perform a quick search
func (sb *Search) SearchHashtag(query string, p ...bool) (*SearchResult, error) {
	if isQuickSearch(p) {
		return sb.tags(query)
	}
	return sb.search(query, sb.tags)
}

// SearchLocation will perform a location search with the provided query.
//
// By default search behavior will be mimicked by sending a search request per
//  added letter, and waiting a few millis in between, just as if you were to
//  type anything into the search bar. However, if you only want to make one
//  search request passing in the full query immediately, you can turn on quick
//  search by passing in one bool:true parameter, like so:
//
//  SearchLocation("myquery", true)  // this will perform a quick search
func (sb *Search) SearchLocation(query string, p ...bool) (*SearchResult, error) {
	if isQuickSearch(p) {
		return sb.places(query)
	}
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
		if errIsFatal(err) {
			return nil, err
		}
		sb.insta.warnHandler("Non fatal error while setting search null state", err)
	}
	return h, nil
}

func (sr *TopSearchItem) RegisterClick() error {
	insta := sr.insta

	var entityType string
	var id int64
	if id = sr.User.ID; id != 0 {
		entityType = "user"
	} else if id = sr.Hashtag.ID; id != 0 {
		entityType = "hashtag"
	} else if id = sr.Place.Location.ID; id != 0 {
		entityType = "place"
	}

	err := insta.sendSearchRegisterRequest(
		map[string]string{
			"entity_id":   toString(id),
			"_uuid":       insta.uuid,
			"entity_type": entityType,
		},
	)
	return err
}

func (sr *SearchResult) RegisterUserClick(user *User) error {
	present := false
	for _, u := range sr.Users {
		if u.ID == user.ID {
			present = true
			break
		}
	}
	if !present {
		return ErrSearchUserNotFound
	}

	err := sr.insta.sendSearchRegisterRequest(
		map[string]string{
			"entity_id":   toString(user.ID),
			"_uuid":       sr.insta.uuid,
			"entity_type": "user",
		},
	)
	return err
}

// RegisterHashtagClick send a register click request, and calls Hashtag.Info()
func (sr *SearchResult) RegisterHashtagClick(h *Hashtag) error {
	present := false
	for _, x := range sr.Tags {
		if x.ID == h.ID {
			present = true
			break
		}
	}
	if !present {
		return ErrSearchUserNotFound
	}

	err := sr.insta.sendSearchRegisterRequest(
		map[string]string{
			"entity_id":   toString(h.ID),
			"_uuid":       sr.insta.uuid,
			"entity_type": "hashtag",
		},
	)
	if err != nil {
		return err
	}
	err = h.Info()
	if err != nil {
		return err
	}

	err = h.Stories()
	return err
}

// RegisterLocationClick send a register click request
func (sr *SearchResult) RegisterLocationClick(l *Location) error {
	present := false
	for _, x := range sr.Places {
		if x.Location.ID == l.ID {
			present = true
			break
		}
	}
	if !present {
		return ErrSearchUserNotFound
	}

	err := sr.insta.sendSearchRegisterRequest(
		map[string]string{
			"entity_id":   toString(l.ID),
			"_uuid":       sr.insta.uuid,
			"entity_type": "place",
		},
	)
	return err
}

func (insta *Instagram) sendSearchRegisterRequest(query map[string]string) error {
	_, _, err := insta.sendRequest(&reqOptions{
		Endpoint: urlSearchRegisterClick,
		IsPost:   true,
		Query:    query,
	})
	return err
}

func (sb *Search) search(query string, fn func(string) (*SearchResult, error)) (*SearchResult, error) {
	insta := sb.insta
	result := &SearchResult{}

	if insta.Discover.NumResults == 0 {
		sb.insta.Discover.Next()
	}
	h, err := sb.history()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get search history")
	}
	result.History = *h
	if err := sb.NullState(); err != nil {
		if errIsFatal(err) {
			return nil, err
		}
		insta.warnHandler("Non fatal error while setting search null state", err)
	}

	var q string
	for _, char := range query {
		q += string(char)
		result, err = fn(q)
		if err != nil {
			return nil, err
		}
		// If the query is a username, and in the top 10, return
		if len(result.Results) >= 10 {
			for _, r := range result.Results[:10] {
				if r.User != nil && r.User.Username == query {
					return result, nil
				}
			}
		}

		s := random(150, 500)
		time.Sleep(time.Duration(s) * time.Millisecond)
	}
	return result, nil
}

func (search *Search) topsearch(query string) (*SearchResult, error) {
	insta := search.insta
	res := &SearchResult{
		insta:         insta,
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
	if err != nil {
		return nil, err
	}
	res.setValues()
	return res, nil
}

func (sr *SearchResult) Error() error {
	return sr.err
}

func (sr *SearchResult) setValues() {
	for _, r := range sr.Results {
		r.insta = sr.insta
		if r.User != nil {
			r.User.insta = sr.insta
		}

		if r.Hashtag != nil {
			r.Hashtag.insta = sr.insta
			r.Hashtag.setValues()
		}
	}

	for _, u := range sr.Users {
		u.insta = sr.insta
	}

	for _, t := range sr.Tags {
		t.insta = sr.insta
		t.setValues()
	}

	for _, l := range sr.Places {
		l.Location.insta = sr.insta
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
		insta:         insta,
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
	if err != nil {
		return nil, err
	}
	res.setValues()
	return res, nil
}

func (search *Search) places(location string) (*SearchResult, error) {
	insta := search.insta
	res := &SearchResult{
		insta:         insta,
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
	if err != nil {
		return nil, err
	}
	res.setValues()
	return res, nil
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

func isQuickSearch(p []bool) bool {
	if len(p) > 0 {
		return p[0]
	}
	return false
}
