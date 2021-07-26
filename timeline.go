package goinsta

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Timeline is the object to represent the main feed on instagram, the first page that shows the latest feeds of my following contacts.
type Timeline struct {
	insta       *Instagram
	err         error
	errChan     chan error
	lastRequest int64
	pullRefresh bool
	sessionID   string
	prevReason  string
	fetchExtra  bool

	endpoint string
	Items    []Item
	Tray     Tray

	More_available            bool
	NextID                    string
	Num_results               float64
	Preload_distance          float64
	Pull_to_refresh_window_ms float64
	Request_id                string
	Session_id                string
}

type FeedCache struct {
	Items []struct {
		Media_or_ad Item `json:"media_or_ad"`
		EndOfFeed   struct {
			Pause    bool   `json:"pause"`
			Title    string `json:"title"`
			Subtitle string `json:"subtitle"`
		} `json:"end_of_feed_demarcator"`
	} `json:"feed_items"`

	More_available                 bool    `json:"more_available"`
	NextID                         string  `json:"next_max_id"`
	Num_results                    float64 `json:"num_results"`
	Pull_to_refresh_window_ms      float64 `json:"pull_to_refresh_window_ms"`
	Request_id                     string  `json:"request_id"`
	Session_id                     string  `json:"session_id"`
	View_state_version             string  `json:"view_state_version"`
	Auto_load_more_enabled         bool    `json:"auto_load_more_enabled"`
	Is_direct_v2_enabled           bool    `json:"is_direct_v2_enabled"`
	Client_feed_changelist_applied bool    `json:"client_feed_changelist_applied"`
	Preload_distance               float64 `json:"preload_distance"`
	Status                         string  `json:"status"`
	Feed_pill_text                 string  `json:"feed_pill_text"`
	Startup_prefetch_configs       struct {
		Explore struct {
			Containermodule            string `json:"containermodule"`
			Should_prefetch            bool   `json:"should_prefetch"`
			Should_prefetch_thumbnails bool   `json:"should_prefetch_thumbnails"`
		} `json:"explore"`
	} `json:"startup_prefetch_configs"`
	Use_aggressive_first_tail_load bool    `json:"use_aggressive_first_tail_load"`
	Hide_like_and_view_counts      float64 `json:"hide_like_and_view_counts"`
}

func newTimeline(insta *Instagram) *Timeline {
	time := &Timeline{
		insta:    insta,
		endpoint: urlTimeline,
		errChan:  make(chan error, 1),
	}
	return time
}

// Next allows pagination after calling:
// User.Feed
// Params: ranked_content is set to "true" by default, you can set it to false by either passing "false" or false as parameter.
// returns false when list reach the end.
// if FeedMedia.Error() is ErrNoMore no problem have been occurred.
// starts first request will be a cold start
func (tl *Timeline) Next() bool {
	if tl.err != nil {
		return false
	}

	insta := tl.insta
	endpoint := tl.endpoint

	// make sure at least 4 sec after last request, at most 6 sec
	var th int64 = 4
	var thR float64 = 2

	// if fetching extra, no big timeout is needed
	if tl.fetchExtra {
		th = 2
		thR = 1
	}

	if delta := time.Now().Unix() - tl.lastRequest; delta < th {
		s := time.Duration(rand.Float64()*thR + float64(th-delta))
		time.Sleep(s * time.Second)
	}
	t := time.Now().Unix()

	var reason string
	isPullToRefresh := "0"
	query := map[string]string{
		"feed_view_info":      "[]",
		"timezone_offset":     timeOffset,
		"device_id":           insta.uuid,
		"request_id":          generateUUID(),
		"_uuid":               insta.uuid,
		"bloks_versioning_id": goInstaBloksVerID,
	}

	catchErr := func() {
		// not the best error handling
		e := <-tl.errChan
		if e != nil {
			fmt.Println("Failed to fetch stories:", e)
			tl.err = e
		}
	}

	var tWarm int64 = 10
	if tl.pullRefresh || (!tl.More_available && t-tl.lastRequest < tWarm*60) {
		reason = "pull_to_refresh"
		isPullToRefresh = "1"
		tl.sessionID = generateUUID()
		go tl.fetchTray("pull_to_refresh")
		defer catchErr()
	} else if tl.lastRequest == 0 || (tl.fetchExtra && tl.prevReason == "warm_start_fetch") {
		reason = "cold_start_fetch"
		tl.sessionID = generateUUID()
		go tl.fetchTray("cold_start")
		defer catchErr()
	} else if t-tl.lastRequest > tWarm*60 { // 10 min
		reason = "warm_start_fetch"
		tl.sessionID = generateUUID()
		go tl.fetchTray("warm_start_with_feed")
		defer catchErr()
	} else if tl.fetchExtra || tl.More_available && tl.NextID != "" {
		reason = "pagination"
		query["max_id"] = tl.NextID
	}

	query["reason"] = reason
	query["is_pull_to_refresh"] = isPullToRefresh
	query["session_id"] = tl.sessionID
	tl.prevReason = reason

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: endpoint,
			IsPost:   true,
			Gzip:     true,
			Query:    query,
			ExtraHeaders: map[string]string{
				"X-Ads-Opt-Out":  "0",
				"X-Google-AD-ID": insta.adid,
				"X-Fb":           "1",
			},
		},
	)
	if err == nil {
		tl.lastRequest = t

		// Decode json
		tmp := FeedCache{}
		d := json.NewDecoder(bytes.NewReader(body))
		d.UseNumber()
		err = d.Decode(&tmp)

		// Add posts to Timeline object
		if err == nil {
			// copy constants over
			tl.NextID = tmp.NextID
			tl.More_available = tmp.More_available
			if tl.fetchExtra {
				tl.Num_results += tmp.Num_results
			} else {
				tl.Num_results = tmp.Num_results
			}
			tl.Preload_distance = tmp.Preload_distance
			tl.Pull_to_refresh_window_ms = tmp.Pull_to_refresh_window_ms
			tl.fetchExtra = false

			// copy post items over
			for i := range tmp.Items {
				if tmp.Items[i].EndOfFeed.Title != "" {
					tl.err = errors.New(
						tmp.Items[i].EndOfFeed.Title + " | " +
							tmp.Items[i].EndOfFeed.Subtitle)
				}
				populateItem(&tmp.Items[i].Media_or_ad, insta)
				tl.Items = append(tl.Items, tmp.Items[i].Media_or_ad)
			}

			// fetch more posts if not enough posts were returned, mimick apk behvaior
			if tmp.Num_results < tmp.Preload_distance && tmp.More_available {
				tl.fetchExtra = true
				tl.Next()
			}

			return true
		} else {
			tl.err = err
		}
	} else {
		tl.err = err
	}

	return false
}

func (tl *Timeline) SetPullRefresh() {
	tl.pullRefresh = true
}

func (tl *Timeline) UnsetPullRefresh() {
	tl.pullRefresh = false
}

func (tl *Timeline) ClearPosts() {
	tl.Items = []Item{}
}

func populateItem(item *Item, insta *Instagram) {
	item.User.insta = insta
	item.Comments = newComments(item)
	for i := range item.CarouselMedia {
		item.CarouselMedia[i].User = item.User
		item.CarouselMedia[i].Comments = newComments(&item.CarouselMedia[i])
	}
}

func (tl *Timeline) fetchTray(reason string) {
	body, _, err := tl.insta.sendRequest(&reqOptions{
		Endpoint: urlStories,
		IsPost:   true,
		Query: map[string]string{
			"supported_capabilities_new": `[{"name":"SUPPORTED_SDK_VERSIONS","value":"100.0,101.0,102.0,103.0,104.0,105.0,106.0,107.0,108.0,109.0,110.0,111.0,112.0,113.0,114.0,115.0,116.0,117.0"},{"name":"FACE_TRACKER_VERSION","value":"14"},{"name":"segmentation","value":"segmentation_enabled"},{"name":"COMPRESSION","value":"ETC2_COMPRESSION"},{"name":"world_tracker","value":"world_tracker_enabled"},{"name":"gyroscope","value":"gyroscope_enabled"}]`,
			"reason":                     reason,
			"timezone_offset":            timeOffset,
			"tray_session_id":            generateUUID(),
			"request_id":                 generateUUID(),
			"_uuid":                      tl.insta.uuid,
		},
	})
	if err != nil {
		tl.errChan <- err
		return
	}

	tray := Tray{}
	err = json.Unmarshal(body, &tray)
	if err != nil {
		tl.errChan <- err
		return
	}

	tray.set(tl.insta, urlStories)
	tl.Tray = tray
	tl.errChan <- nil
}

func (tl *Timeline) Refresh() error {
	tl.SetPullRefresh()
	if !tl.Next() {
		return tl.err
	}
	return nil
}

// helper function to get the stories
func (tl *Timeline) Stories() []StoryMedia {
	return tl.Tray.Stories
}

// helper function to get the Broadcasts
func (tl *Timeline) Broadcasts() []Broadcast {
	return tl.Tray.Broadcasts
}
