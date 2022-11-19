package goinsta

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrAllSaved = errors.New("unable to call function for collection all posts")

// MediaItem defines a item media for the
// SavedMedia struct
type MediaItem struct {
	Media Item `json:"media"`
}

// SavedMedia stores information about ALL your saved posts, regardless of their collection.
// This is the same to vising your saved posts page, and clicking "All Posts".
// If you want to view a single collection, use the Collections type.
type SavedMedia struct {
	insta    *Instagram
	endpoint string

	err error

	Items []MediaItem `json:"items"`

	NumResults          int    `json:"num_results"`
	MoreAvailable       bool   `json:"more_available"`
	AutoLoadMoreEnabled bool   `json:"auto_load_more_enabled"`
	Status              string `json:"status"`

	NextID interface{} `json:"next_max_id"`
}

// Collections stores information about all your collection.
// The first collection will always be the "All Posts" collcetion.
// To fetch all collections, you can paginate with Collections.Next()
type Collections struct {
	insta *Instagram
	err   error

	AutoLoadMoreEnabled bool          `json:"auto_load_more_enabled"`
	Items               []*Collection `json:"items"`
	MoreAvailable       bool          `json:"more_available"`
	NextID              string        `json:"next_max_id"`
	NumResults          int
	Status              string `json:"status"`
}

// Collection represents a single collection. All collections will not load
//   any posts by default. To load the posts you need to call Collection.Next()
// You can edit your collections with their respective methods. e.g. edit the name
//   or change the thumbnail.
type Collection struct {
	insta *Instagram
	err   error
	all   *SavedMedia

	ID         string `json:"collection_id"`
	MediaCount int    `json:"collection_media_count"`
	Name       string `json:"collection_name"`
	Type       string `json:"collection_type"`
	Cover      struct {
		ID             string `json:"id"`
		Images         Images `json:"image_versions2"`
		OriginalWidth  int    `json:"original_width"`
		OriginalHeight int    `json:"original_height"`
		MediaType      int    `json:"media_type"`
	} `json:"cover_media"`

	Items         []Item
	NumResults    int         `json:"num_results"`
	MoreAvailable bool        `json:"more_available"`
	NextID        interface{} `json:"next_max_id"`
}

type collectionSync struct {
	Clips  SavedMedia `json:"saved_clips_response"`
	IGTV   SavedMedia `json:"saved_igtv_response"`
	Media  SavedMedia `json:"saved_media_response"`
	Status string     `json:"status"`
}

func newCollections(insta *Instagram) *Collections {
	return &Collections{
		insta: insta,
	}
}

// Next allows you to fetch and paginate your list of collections.
// This method will cumulatively add to the collections. To get the latest
//   fetched collections, call Collections.Latest(), or index with Collections.LastCount
func (c *Collections) Next() bool {
	// Check if prev returned error
	if c.err != nil {
		return false
	}

	// No more available
	if len(c.Items) > 0 && !c.MoreAvailable {
		return false
	}
	insta := c.insta
	query := map[string]string{
		"collection_types": "[\"ALL_MEDIA_AUTO_COLLECTION\",\"PRODUCT_AUTO_COLLECTION\",\"MEDIA\",\"AUDIO_AUTO_COLLECTION\",\"GUIDES_AUTO_COLLECTION\"]",
	}
	if c.NextID != "" {
		query["max_id"] = c.NextID
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlCollectionsList,
			Query:    query,
		},
	)
	if err != nil {
		c.err = err
		return false
	}
	tmp := Collections{}
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		c.err = err
		return false
	}
	c.MoreAvailable = tmp.MoreAvailable
	c.NextID = tmp.NextID
	c.AutoLoadMoreEnabled = tmp.AutoLoadMoreEnabled
	c.Status = tmp.Status
	c.NumResults = len(tmp.Items)

	for _, i := range tmp.Items {
		i.insta = insta
	}
	c.Items = append(c.Items, tmp.Items...)
	if !c.MoreAvailable {
		c.err = ErrNoMore
	}

	return c.MoreAvailable
}

// Latest will return the last fetched items by indexing with Collections.LastCount.
// Collections.Next keeps adding to the items, this method only returns the latest items.
func (c *Collections) Latest() []*Collection {
	return c.Items[len(c.Items)-c.NumResults:]
}

// Error returns the error if one occured in Collections.Next()
func (c *Collections) Error() error {
	return c.err
}

// Create allows you to create a new collection.
// :param: name (*required) - the name of the collection
// :param: ...Items (optional) - posts to add to the collection
func (c *Collections) Create(name string, items ...Item) (*Collection, error) {
	insta := c.insta

	module := "collection_create"
	if len(items) == 1 {
		module = "feed_timeline"
	}

	ids := []string{}
	for _, i := range items {
		ids = append(ids, i.GetID())
	}
	mediaIDs, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(map[string]string{
		"module_name":     module,
		"added_media_ids": string(mediaIDs),
		"_uid":            toString(insta.Account.ID),
		"name":            name,
		"_uuid":           insta.uuid,
	})
	if err != nil {
		return nil, err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlCollectionsCreate,
			IsPost:   true,
			Query:    generateSignature(data),
		},
	)
	if err != nil {
		return nil, err
	}
	n := Collection{insta: insta}
	err = json.Unmarshal(body, &n)
	return &n, err
}

// Delete will permanently delete a collection
func (c *Collection) Delete() error {
	if c.Name == "ALL_MEDIA_AUTO_COLLECTION" {
		return ErrAllSaved
	}
	insta := c.insta

	data, err := json.Marshal(
		map[string]string{
			"_uid":  toString(insta.Account.ID),
			"_uuid": insta.uuid,
		},
	)
	if err != nil {
		return err
	}
	_, _, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlCollectionDelete, c.ID),
			IsPost:   true,
			Query:    generateSignature(data),
		},
	)
	return err
}

// ChangeCover will change to cover of the collection. The item parameter must
//   must be an item that is present inside the collection.
func (c *Collection) ChangeCover(item Item) error {
	if c.Name == "ALL_MEDIA_AUTO_COLLECTION" {
		return ErrAllSaved
	}
	insta := c.insta

	data, err := json.Marshal(
		map[string]string{
			"cover_media_id":         item.GetID(),
			"_uid":                   toString(insta.Account.ID),
			"name":                   c.Name,
			"_uuid":                  insta.uuid,
			"added_collaborator_ids": "[]",
		},
	)
	if err != nil {
		return err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlCollectionEdit, c.ID),
			IsPost:   true,
			Query:    generateSignature(data),
		},
	)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, c); err != nil {
		return err
	}

	return err
}

// ChangeName will change the name of the collection
func (c *Collection) ChangeName(name string) error {
	if c.Name == "ALL_MEDIA_AUTO_COLLECTION" {
		return ErrAllSaved
	}
	insta := c.insta

	data, err := json.Marshal(
		map[string]string{
			"cover_media_id":         toString(c.Cover.ID),
			"_uid":                   toString(insta.Account.ID),
			"name":                   name,
			"_uuid":                  insta.uuid,
			"added_collaborator_ids": "[]",
		},
	)
	if err != nil {
		return err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlCollectionEdit, c.ID),
			IsPost:   true,
			Query:    generateSignature(data),
		},
	)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, c); err != nil {
		return err
	}

	return err
}

// AddCollaborators should in theory add collaborators. Untested.
func (c *Collection) AddCollaborators(colab ...User) error {
	if c.Name == "ALL_MEDIA_AUTO_COLLECTION" {
		return ErrAllSaved
	}
	insta := c.insta

	ids := []string{}
	for _, u := range colab {
		ids = append(ids, toString(u.ID))
	}
	colabIDs, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	data, err := json.Marshal(
		map[string]string{
			"cover_media_id":         toString(c.Cover.ID),
			"_uid":                   toString(insta.Account.ID),
			"name":                   c.Name,
			"_uuid":                  insta.uuid,
			"added_collaborator_ids": string(colabIDs),
		},
	)
	if err != nil {
		return err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlCollectionEdit, c.ID),
			IsPost:   true,
			Query:    generateSignature(data),
		},
	)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, c); err != nil {
		return err
	}

	return err
}

// RemoveMedia will remove media from a collection. The items provided must be
//   inside the collection
func (c *Collection) RemoveMedia(items ...Item) error {
	if c.Name == "ALL_MEDIA_AUTO_COLLECTION" {
		return ErrAllSaved
	}
	insta := c.insta

	ids := []string{}
	for _, u := range items {
		ids = append(ids, u.GetID())
	}
	itemIDs, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	data, err := json.Marshal(
		map[string]string{
			"module_name":       "feed_saved_collections",
			"_uid":              toString(insta.Account.ID),
			"_uuid":             insta.uuid,
			"removed_media_ids": string(itemIDs),
		},
	)
	if err != nil {
		return err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlCollectionEdit, c.ID),
			IsPost:   true,
			Query:    generateSignature(data),
		},
	)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, c); err != nil {
		return err
	}

	return err
}

// Sync will fetch the initial items inside a collection.
// The first call to fetch posts will always be sync, however you can also only
//   call Collection.Next() as Sync() will automatically be called if required.
func (c *Collection) Sync() error {
	if c.Name == "ALL_MEDIA_AUTO_COLLECTION" {
		if !c.Next() {
			return c.err
		}
		return nil
	}
	insta := c.insta
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlCollectionFeedAll, c.ID),
			Query: map[string]string{
				"include_igtv_preview":    "true",
				"include_igtv_tab":        "true",
				"include_clips_subtab":    "false", // default is false, but could be set to true
				"clips_subtab_first":      "false",
				"show_igtv_first":         "false",
				"include_collection_info": "false",
			},
		},
	)
	if err != nil {
		return err
	}

	tmp := collectionSync{}
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		return err
	}

	c.NextID = tmp.Media.NextID
	c.MoreAvailable = tmp.Media.MoreAvailable
	c.NumResults = tmp.Media.NumResults

	c.Items = []Item{}
	for _, i := range tmp.Media.Items {
		c.Items = append(c.Items, i.Media)
	}
	c.setValues()

	return nil
}

// Next allows you to fetch and paginate over all items present inside a collection.
func (c *Collection) Next(params ...interface{}) bool {
	if c.err != nil {
		return false
	}

	if c.Name == "ALL_MEDIA_AUTO_COLLECTION" {
		if c.all == nil {
			c.all = &SavedMedia{
				insta:    c.insta,
				endpoint: urlFeedSavedPosts,
			}
		}
		r := c.all.Next()
		c.NextID = c.all.NextID
		c.MoreAvailable = c.all.MoreAvailable
		c.NumResults = c.all.NumResults
		c.Items = []Item{}
		for _, i := range c.all.Items {
			c.Items = append(c.Items, i.Media)
		}

		c.err = c.all.err
		return r
	}
	if len(c.Items) == 0 {
		if err := c.Sync(); err != nil {
			c.err = err
			return false
		}
	}
	insta := c.insta
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlCollectionFeedPosts, c.ID),
			Query: map[string]string{
				"max_id": c.GetNextID(),
			},
		},
	)
	if err != nil {
		c.err = err
		return false
	}

	tmp := SavedMedia{}
	if err := json.Unmarshal(body, &tmp); err != nil {
		c.err = err
		return false
	}

	c.NextID = tmp.NextID
	c.MoreAvailable = tmp.MoreAvailable
	c.NumResults = tmp.NumResults

	for _, i := range tmp.Items {
		c.Items = append(c.Items, i.Media)
	}
	c.setValues()

	if !c.MoreAvailable {
		c.err = ErrNoMore
	}
	return c.MoreAvailable
}

func (c *Collection) setValues() {
	for i := range c.Items {
		c.Items[i].insta = c.insta
		c.Items[i].User.insta = c.insta
		setToItem(&c.Items[i], c)
	}
}

// Error will return the error if one occured during Collection.Next().
func (c *Collection) Error() error {
	return c.err
}

func (media *Collection) getInsta() *Instagram {
	return media.insta
}

// GetNextID will return the pagination ID as a string
func (c *Collection) GetNextID() string {
	return formatID(c.NextID)
}

// Save saves media item.
//
// You can get saved media using Account.Saved()
func (item *Item) Save() error {
	return item.changeSave(urlMediaSave)
}

// Saveto allows you to save a media item to a specific collection
func (item *Item) SaveTo(c *Collection) error {
	insta := item.insta

	// Call genereral save method first, as in the app this will always happen
	err := item.Save()
	if err != nil {
		if errIsFatal(err) {
			return err
		}
		insta.warnHandler(errors.New("non fatal error, failed to save post to all"))
	}

	data, err := json.Marshal(
		map[string]string{
			"module_name":          "feed_timeline",
			"radio_type":           "wifi-none",
			"_uid":                 toString(insta.Account.ID),
			"_uuid":                insta.uuid,
			"added_collection_ids": fmt.Sprintf("[%s]", c.ID),
		},
	)
	if err != nil {
		return err
	}
	_, _, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlMediaSave, item.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// Unsave unsaves media item.
func (item *Item) Unsave() error {
	return item.changeSave(urlMediaUnsave)
}

func (item *Item) changeSave(endpoint string) error {
	insta := item.insta
	query := map[string]string{
		"module_name":     "feed_timeline",
		"client_position": toString(item.Index),
		"nav_chain":       "",
		"_uid":            toString(insta.Account.ID),
		"_uuid":           insta.uuid,
		"radio_type":      "wifi-none",
	}
	if item.IsCommercial {
		query["delivery_class"] = "ad"
	} else {
		query["delivery_class"] = "organic"
	}
	if item.InventorySource != "" {
		query["inventory_source"] = item.InventorySource
	}
	if len(item.CarouselMedia) > 0 || item.CarouselParentID != "" {
		query["carousel_index"] = "0"
	}
	data, err := json.Marshal(query)
	if err != nil {
		return err
	}

	_, _, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(endpoint, item.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// Sync will fetch the initial saved items.
// The first call to fetch posts will always be sync, however you can also only
//   call SavedMedia.Next() as Sync() will automatically be called if required.
func (media *SavedMedia) Sync() error {
	insta := media.insta
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlFeedSaved,
			Query: map[string]string{
				"include_igtv_preview":    "true",
				"include_igtv_tab":        "true",
				"include_clips_subtab":    "false", // default is false, but could be set to true
				"clips_subtab_first":      "false",
				"show_igtv_first":         "false",
				"include_collection_info": "false",
			},
		},
	)
	if err != nil {
		return err
	}

	tmp := collectionSync{}
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		return err
	}

	media.NextID = tmp.Media.NextID
	media.MoreAvailable = tmp.Media.MoreAvailable
	media.NumResults = tmp.Media.NumResults

	media.Items = tmp.Media.Items
	media.setValues()

	return nil
}

// Next allows pagination of "All Posts" collection
func (media *SavedMedia) Next(params ...interface{}) bool {
	// Return if last pagination had errors
	if media.err != nil {
		return false
	}
	if len(media.Items) == 0 {
		if err := media.Sync(); err != nil {
			media.err = err
			return false
		}
	}

	insta := media.insta
	endpoint := media.endpoint

	opts := &reqOptions{
		Endpoint: endpoint,
		Query: map[string]string{
			"max_id":               media.GetNextID(),
			"include_igtv_preview": "true",
		},
	}

	body, _, err := insta.sendRequest(opts)
	if err != nil {
		media.err = err
		return false
	}

	m := SavedMedia{}
	if err := json.Unmarshal(body, &m); err != nil {
		media.err = err
		return false
	}

	media.MoreAvailable = m.MoreAvailable
	media.NextID = m.NextID
	media.AutoLoadMoreEnabled = m.AutoLoadMoreEnabled
	media.Status = m.Status
	media.NumResults = m.NumResults
	media.Items = append(media.Items, m.Items...)
	media.setValues()

	if m.NextID == 0 || !m.MoreAvailable {
		media.err = ErrNoMore
	}

	return m.MoreAvailable
}

// Error returns the SavedMedia error
func (media *SavedMedia) Error() error {
	return media.err
}

// ID returns the Sav
func (media *SavedMedia) GetNextID() string {
	return formatID(media.NextID)
}

// Delete will unsave ALL saved items.
func (media *SavedMedia) Delete() error {
	for _, i := range media.Items {
		err := i.Media.Unsave()
		if err != nil {
			return err
		}
	}
	return nil
}

func (media *SavedMedia) getInsta() *Instagram {
	return media.insta
}

// setValues set the SavedMedia items values
func (media *SavedMedia) setValues() {
	for i := range media.Items {
		media.Items[i].Media.insta = media.insta
		media.Items[i].Media.User.insta = media.insta
		setToMediaItem(&media.Items[i], media)
	}
}
