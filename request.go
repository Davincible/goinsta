package goinsta

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type reqOptions struct {
	// Connection is connection header. Default is "close".
	Connection string

	// Endpoint is the request path of instagram api
	Endpoint string

	// IsPost set to true will send request with POST method.
	//
	// By default this option is false.
	IsPost bool

	// Compress post form data with gzip
	Gzip bool

	// UseV2 is set when API endpoint uses v2 url.
	UseV2 bool

	// Use b.i.instagram.com
	Useb bool

	// Query is the parameters of the request
	//
	// This parameters are independents of the request method (POST|GET)
	Query map[string]string

	// List of headers to ignore
	IgnoreHeaders []string

	// Extra headers to add
	ExtraHeaders map[string]string

	// Timestamp
	Timestamp string
}

func (insta *Instagram) sendSimpleRequest(uri string, a ...interface{}) (body []byte, err error) {
	body, _, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(uri, a...),
		},
	)
	return body, err
}

func (insta *Instagram) sendRequest(o *reqOptions) (body []byte, h http.Header, err error) {
	checkCookieExpiry(insta)

	method := "GET"
	if o.IsPost {
		method = "POST"
	}
	if o.Connection == "" {
		o.Connection = "close"
	}

	if o.Timestamp == "" {
		o.Timestamp = strconv.Itoa(int(time.Now().Unix()))
	}

	var nu string
	if o.Useb {
		nu = goInstaAPIUrlb
	} else {
		nu = goInstaAPIUrl
	}
	if o.UseV2 && !o.Useb {
		nu = goInstaAPIUrlv2
	} else if o.UseV2 && o.Useb {
		nu = goInstaAPIUrlv2b
	}

	u, err := url.Parse(nu + o.Endpoint)
	if err != nil {
		return nil, nil, err
	}

	vs := url.Values{}
	bf := bytes.NewBuffer([]byte{})

	for k, v := range o.Query {
		vs.Add(k, v)
	}

	var contentEncoding string
	if o.IsPost && o.Gzip {
		zw := gzip.NewWriter(bf)
		if _, err := zw.Write([]byte(vs.Encode())); err != nil {
			return nil, nil, err
		}
		if err := zw.Close(); err != nil {
			return nil, nil, err
		}
		contentEncoding = "gzip"
	} else if o.IsPost {
		bf.WriteString(vs.Encode())
	} else {
		for k, v := range u.Query() {
			vs.Add(k, strings.Join(v, " "))
		}

		u.RawQuery = vs.Encode()
	}

	var req *http.Request
	req, err = http.NewRequest(method, u.String(), bf)
	if err != nil {
		return
	}

	ignoreHeader := func(h string) bool {
		for _, k := range o.IgnoreHeaders {
			if k == h {
				return true
			}
		}
		return false
	}

	setHeaders := func(r *http.Request, h map[string]string) {
		for k, v := range h {
			if v != "" && !ignoreHeader(k) {
				r.Header.Set(k, v)
			}
		}
	}

	accID := "0"
	if insta.Account != nil {
		accID = strconv.Itoa(int(insta.Account.ID))
	}

	headers := map[string]string{
		"Accept-Language":             locale,
		"Accept-Encoding":             "gzip,deflate",
		"Authorization":               insta.headerOptions["Authorization"],
		"Connection":                  o.Connection,
		"Content-Type":                "application/x-www-form-urlencoded; charset=UTF-8",
		"Ig-Intended-User-Id":         accID,
		"Ig-U-Shbid":                  insta.headerOptions["Ig-U-Shbid"],
		"Ig-U-Shbts":                  insta.headerOptions["Ig-U-Shbts"],
		"Ig-U-Ds-User-Id":             insta.headerOptions["Ig-U-Ds-User-Id"],
		"Ig-U-Rur":                    insta.headerOptions["Ig-U-Rur"],
		"Ig-U-Ig-Direct-Region-Hint":  insta.headerOptions["Ig-U-Ig-Direct-Region-Hint"],
		"User-Agent":                  userAgent,
		"X-Ig-App-Locale":             locale,
		"X-Ig-Device-Locale":          locale,
		"X-Ig-Mapped-Locale":          locale,
		"X-Ig-App-Id":                 fbAnalytics,
		"X-Ig-Device-Id":              insta.uuid,
		"X-Ig-Family-Device-Id":       insta.fID,
		"X-Ig-Android-Id":             insta.dID,
		"X-Ig-Timezone-Offset":        timeOffset,
		"X-Ig-Capabilities":           igCapabilities,
		"X-Ig-Connection-Type":        connType,
		"X-Pigeon-Session-Id":         insta.psID,
		"X-Pigeon-Rawclienttime":      fmt.Sprintf("%s.%d", o.Timestamp, acquireRand(100, 900)),
		"X-Ig-Bandwidth-Speed-KBPS":   fmt.Sprintf("%d.000", acquireRand(1000, 9000)),
		"X-Ig-Bandwidth-TotalBytes-B": strconv.Itoa(acquireRand(1000000, 5000000)),
		"X-Ig-Bandwidth-Totaltime-Ms": strconv.Itoa(acquireRand(200, 800)),
		"X-Ig-App-Startup-Country":    "unkown",
		"X-Ig-Www-Claim":              insta.headerOptions["X-Ig-Www-Claim"],
		"X-Bloks-Version-Id":          goInstaBloksVerID,
		"X-Bloks-Is-Layout-Rtl":       "false",
		"X-Bloks-Is-Panorama-Enabled": "true",
		"X-Mid":                       insta.headerOptions["X-Mid"],
		"X-Fb-Http-Engine":            "Liger",
		"X-Fb-Client-Ip":              "True",
		"X-Fb-Server-Cluster":         "True",
	}
	if contentEncoding != "" {
		headers["Content-Encoding"] = contentEncoding
	}

	setHeaders(req, headers)
	setHeaders(req, o.ExtraHeaders)

	resp, err := insta.c.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		err = isError(resp.StatusCode, body)
	}
	extractHeaders(insta, resp.Header)

	return body, resp.Header.Clone(), err
}

func checkCookieExpiry(insta *Instagram) {
	if insta.xmidExpiry != -1 && time.Now().Unix() > insta.xmidExpiry-10 {
		insta.xmidExpiry = -1
		insta.zrToken()
	}
}

func extractHeaders(insta *Instagram, h http.Header) {
	extract := func(in string, out string) {
		x := h[in]
		if len(x) > 0 && x[0] != "" {
			insta.headerOptions[out] = x[0]
		}
	}

	extract("Ig-Set-Authorization", "Authorization")
	extract("Ig-Set-X-Mid", "X-Mid")
	extract("X-Ig-Set-Www-Claim", "X-Ig-Www-Claim")
	extract("Ig-Set-Ig-U-Ig-Direct-Region-Hint", "Ig-U-Ig-Direct-Region-Hint")
	extract("Ig-Set-Ig-U-Shbid", "Ig-U-Shbid")
	extract("Ig-Set-Ig-U-Shbts", "Ig-U-Shbts")
	extract("Ig-Set-Ig-U-Rur", "Ig-U-Rur")
	extract("Ig-Set-Ig-U-Ds-User-Id", "Ig-U-Ds-User-Id")
}

func isError(code int, body []byte) (err error) {
	switch code {
	case 200:
	case 503:
		return Error503{
			Message: "Instagram API error. Try it later.",
		}
	case 400:
		ierr := Error400{}
		err = json.Unmarshal(body, &ierr)
		if err != nil {
			return err
		}

		if ierr.Message == "challenge_required" {
			return ierr.ChallengeError
		}

		if err == nil && ierr.Message != "" {
			return ierr
		}
	default:
		ierr := ErrorN{}
		err = json.Unmarshal(body, &ierr)
		if err != nil {
			return err
		}
		return ierr
	}
	return nil
}

func (insta *Instagram) prepareData(other ...map[string]interface{}) (string, error) {
	data := map[string]interface{}{
		"_uuid":      insta.uuid,
		"_csrftoken": insta.token,
	}
	if insta.Account != nil && insta.Account.ID != 0 {
		data["_uid"] = strconv.FormatInt(insta.Account.ID, 10)
	}

	for i := range other {
		for key, value := range other[i] {
			data[key] = value
		}
	}
	b, err := json.Marshal(data)
	if err == nil {
		return b2s(b), err
	}
	return "", err
}

func (insta *Instagram) prepareDataQuery(other ...map[string]interface{}) map[string]string {
	data := map[string]string{
		"_uuid":      insta.uuid,
		"_csrftoken": insta.token,
	}
	for i := range other {
		for key, value := range other[i] {
			data[key] = toString(value)
		}
	}
	return data
}

func acquireRand(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
