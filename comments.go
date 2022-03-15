package goinsta

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Comments allows user to interact with media (item) comments.
// You can Add or Delete by index or by user name.
type Comments struct {
	item     *Item
	endpoint string
	err      error

	Items                          []Comment       `json:"comments"`
	CommentCount                   int64           `json:"comment_count"`
	Caption                        Caption         `json:"caption"`
	CaptionIsEdited                bool            `json:"caption_is_edited"`
	HasMoreComments                bool            `json:"has_more_comments"`
	HasMoreHeadloadComments        bool            `json:"has_more_headload_comments"`
	ThreadingEnabled               bool            `json:"threading_enabled"`
	MediaHeaderDisplay             string          `json:"media_header_display"`
	InitiateAtTop                  bool            `json:"initiate_at_top"`
	InsertNewCommentToTop          bool            `json:"insert_new_comment_to_top"`
	PreviewComments                []Comment       `json:"preview_comments"`
	NextID                         json.RawMessage `json:"next_max_id,omitempty"`
	NextMinID                      json.RawMessage `json:"next_min_id,omitempty"`
	CommentLikesEnabled            bool            `json:"comment_likes_enabled"`
	DisplayRealtimeTypingIndicator bool            `json:"display_realtime_typing_indicator"`
	Status                         string          `json:"status"`
}

func (comments *Comments) setValues() {
	for i := range comments.Items {
		comments.Items[i].item = comments.item
		comments.Items[i].setValues(comments.item.insta)
	}
}

func newComments(item *Item) *Comments {
	c := &Comments{
		item: item,
	}
	return c
}

func (comments Comments) Error() error {
	return comments.err
}

// Disable disables comments in FeedMedia.
//
// See example: examples/media/commentDisable.go
func (comments *Comments) Disable() error {
	return comments.toggleComments(urlCommentDisable)
}

// Enable enables comments in FeedMedia
//
// See example: examples/media/commentEnable.go
func (comments *Comments) Enable() error {
	return comments.toggleComments(urlCommentEnable)
}

func (comments *Comments) toggleComments(endpoint string) error {
	// Use something else to repel stories
	// switch comments.item.media.(type) {
	// case *StoryMedia:
	// 	return fmt.Errorf("Incompatible type. Cannot use Enable() with StoryMedia type")
	// }
	insta := comments.item.insta

	_, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(endpoint, comments.item.ID),
			Query:    map[string]string{"_uuid": insta.uuid},
			IsPost:   true,
		},
	)
	return err
}

// Next allows comment pagination.
//
// This function support concurrency methods to get comments using Last and Next ID
//
// New comments are stored in Comments.Items
func (comments *Comments) Next() bool {
	if comments.err != nil {
		return false
	}

	item := comments.item
	insta := item.insta
	endpoint := comments.endpoint
	query := map[string]string{
		// "can_support_threading": "true",
	}
	if comments.NextID != nil {
		next, _ := strconv.Unquote(string(comments.NextID))
		query["max_id"] = next
	} else if comments.NextMinID != nil {
		next, _ := strconv.Unquote(string(comments.NextMinID))
		query["min_id"] = next
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint:   endpoint,
			Connection: "keep-alive",
			Query:      query,
		},
	)
	if err == nil {
		c := Comments{}
		err = json.Unmarshal(body, &c)
		if err == nil {
			*comments = c
			comments.endpoint = endpoint
			comments.item = item
			if (!comments.HasMoreComments || comments.NextID == nil) &&
				(!comments.HasMoreHeadloadComments || comments.NextMinID == nil) {
				comments.err = ErrNoMore
			}
			comments.setValues()
			return true
		}
	}
	comments.err = err
	return false
}

// Sync prepare Comments to receive comments.
// Use Next to receive comments.
//
// See example: examples/media/commentsSync.go
func (comments *Comments) Sync() {
	endpoint := fmt.Sprintf(urlCommentSync, comments.item.ID)
	comments.endpoint = endpoint
	return
}

// Add push a comment in media.
//
// See example: examples/media/commentsAdd.go
func (comments *Comments) Add(text string) (err error) {
	return comments.item.Comment(text)
}

// Delete deletes a single comment.
func (c *Comment) Delete() error {
	return c.item.insta.bulkDelComments([]*Comment{c})
}

// BulkDelete allows you to select and delete multiple comments on a single post.
func (comments *Comments) BulkDelete(c []*Comment) error {
	return comments.item.insta.bulkDelComments(c)
}

func (insta *Instagram) bulkDelComments(c []*Comment) error {
	if len(c) == 0 {
		return nil
	}
	cIDs := toString(c[0].ID)
	pID := c[0].item.ID
	if len(c) > 1 {
		for _, i := range c[1:] {
			if i.item.ID != pID {
				return errors.New("All comments have to belong to the same post")
			}
			cIDs += "," + toString(i.ID)
		}
	}

	data, err := json.Marshal(
		map[string]string{
			"comment_ids_to_delete": cIDs,
			"_uid":                  toString(insta.Account.ID),
			"_uuid":                 insta.uuid,
			"container_module":      "feed_timeline",
		},
	)
	if err != nil {
		return err
	}
	_, _, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlCommentBulkDelete, pID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// DeleteMine removes all of your comments limited by parsed parameter.
//
// If limit is <= 0 DeleteMine will delete all your comments.
// Be careful with using this on posts with a large number of comments,
//  as a large number of requests will be made to index all comments.This
//  can result in a ratelimiter being tripped.
//
// See example: examples/media/commentsDelMine.go
func (comments *Comments) DeleteMine(limit int) error {
	comments.Sync()
	cList := make([]*Comment, 1)

	insta := comments.item.insta
floop:
	for i := 0; comments.Next(); i++ {
		for _, c := range comments.Items {
			if c.UserID == insta.Account.ID || c.User.ID == insta.Account.ID {
				if limit > 0 && i >= limit {
					break floop
				}
				cList = append(cList, &c)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err := comments.Error(); err != nil && err != ErrNoMore {
		return err
	}
	err := comments.BulkDelete(cList)
	return err
}

// Comment is a type of Media retrieved by the Comments methods
type Comment struct {
	insta *Instagram
	idstr string
	item  *Item

	ID                             interface{} `json:"pk"`
	Text                           string      `json:"text"`
	Type                           interface{} `json:"type"`
	User                           User        `json:"user"`
	UserID                         int64       `json:"user_id"`
	BitFlags                       int         `json:"bit_flags"`
	ChildCommentCount              int         `json:"child_comment_count"`
	CommentIndex                   int         `json:"comment_index"`
	CommentLikeCount               int         `json:"comment_like_count"`
	ContentType                    string      `json:"content_type"`
	CreatedAt                      int64       `json:"created_at"`
	CreatedAtUtc                   int64       `json:"created_at_utc"`
	DidReportAsSpam                bool        `json:"did_report_as_spam"`
	HasLikedComment                bool        `json:"has_liked_comment"`
	InlineComposerDisplayCondition string      `json:"inline_composer_display_condition"`
	OtherPreviewUsers              []*User     `json:"other_preview_users"`
	PreviewChildComments           []Comment   `json:"preview_child_comments"`
	NextMaxChildCursor             string      `json:"next_max_child_cursor,omitempty"`
	HasMoreTailChildComments       bool        `json:"has_more_tail_child_comments,omitempty"`
	NextMinChildCursor             string      `json:"next_min_child_cursor,omitempty"`
	HasMoreHeadChildComments       bool        `json:"has_more_head_child_comments,omitempty"`
	NumTailChildComments           int         `json:"num_tail_child_comments,omitempty"`
	NumHeadChildComments           int         `json:"num_head_child_comments,omitempty"`
	ShareEnabled                   bool        `json:"share_enabled"`
	IsCovered                      bool        `json:"is_covered"`
	PrivateReplyStatus             int64       `json:"private_reply_status"`
	SupporterInfo                  struct {
		SupportTier string `json:"support_tier"`
		BadgesCount int    `json:"badges_count"`
	} `json:"supporter_info"`
	Status string `json:"status"`
}

func (c *Comment) setValues(insta *Instagram) {
	c.User.insta = insta
	for i := range c.OtherPreviewUsers {
		c.OtherPreviewUsers[i].insta = insta
	}
	for i := range c.PreviewChildComments {
		c.PreviewChildComments[i].setValues(insta)
	}
}

func (c Comment) getid() string {
	return toString(c.ID)
}

// Like likes comment.
func (c *Comment) Like() error {
	return c.changeLike(urlCommentLike)
}

// Unlike unlikes comment.
func (c *Comment) Unlike() error {
	return c.changeLike(urlCommentUnlike)
}

func (c *Comment) changeLike(endpoint string) error {
	insta := c.insta
	item := c.item
	query := map[string]string{
		"feed_position":           "0",
		"container_module":        "feed_timeline",
		"nav_chain":               "",
		"_uid":                    toString(insta.Account.ID),
		"_uuid":                   insta.uuid,
		"radio_type":              "wifi-none",
		"is_carousel_bumped_post": "false", // not sure when this would be true
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

	_, _, err = c.insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(endpoint, c.getid()),
			IsPost:   true,
			Query:    generateSignature(data),
		},
	)
	return err
}
