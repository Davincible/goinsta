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

	// Omit API omit the /api/v1/ part of the url
	OmitAPI bool

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

	// DataBytes can be used to pass raw data to a request, instead of a
	//   form using the Query param. This is used for e.g. photo and vieo uploads.
	DataBytes *bytes.Buffer

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
	insta.checkXmidExpiry()

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
		nu = instaAPIUrlb
	} else {
		nu = instaAPIUrl
	}
	if o.UseV2 && !o.Useb {
		nu = instaAPIUrlv2
	} else if o.UseV2 && o.Useb {
		nu = instaAPIUrlv2b
	}
	if o.OmitAPI {
		nu = baseUrl
		o.IgnoreHeaders = append(o.IgnoreHeaders, omitAPIHeadersExclude...)
	}

	u, err := url.Parse(nu + o.Endpoint)
	if err != nil {
		return nil, nil, err
	}

	vs := url.Values{}
	bf := bytes.NewBuffer([]byte{})
	reqData := bytes.NewBuffer([]byte{})

	for k, v := range o.Query {
		vs.Add(k, v)
	}

	// If DataBytes has been passed, use that as data, else use Query
	if o.DataBytes != nil {
		reqData = o.DataBytes
	} else {
		reqData.WriteString(vs.Encode())
	}

	var contentEncoding string
	if o.IsPost && o.Gzip {
		// If gzip encoding needs to be applied
		zw := gzip.NewWriter(bf)
		defer zw.Close()
		if _, err := zw.Write(reqData.Bytes()); err != nil {
			return nil, nil, err
		}
		if err := zw.Close(); err != nil {
			return nil, nil, err
		}
		contentEncoding = "gzip"
	} else if o.IsPost {
		// use post form if POST request
		bf = reqData
	} else {
		// append query to url if GET request
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

	setHeaders := func(h map[string]string) {
		for k, v := range h {
			if v != "" && !ignoreHeader(k) {
				req.Header.Set(k, v)
			}
		}
	}

	headers := map[string]string{
		"Accept-Language":             locale,
		"Accept-Encoding":             "gzip,deflate",
		"Connection":                  o.Connection,
		"Content-Type":                "application/x-www-form-urlencoded; charset=UTF-8",
		"User-Agent":                  insta.userAgent,
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
		"X-Pigeon-Rawclienttime":      fmt.Sprintf("%s.%d", o.Timestamp, random(100, 900)),
		"X-Ig-Bandwidth-Speed-KBPS":   fmt.Sprintf("%d.000", random(1000, 9000)),
		"X-Ig-Bandwidth-TotalBytes-B": strconv.Itoa(random(1000000, 5000000)),
		"X-Ig-Bandwidth-Totaltime-Ms": strconv.Itoa(random(200, 800)),
		"X-Ig-App-Startup-Country":    "unkown",
		"X-Bloks-Version-Id":          bloksVerID,
		"X-Bloks-Is-Layout-Rtl":       "false",
		"X-Bloks-Is-Panorama-Enabled": "true",
		"X-Fb-Http-Engine":            "Liger",
		"X-Fb-Client-Ip":              "True",
		"X-Fb-Server-Cluster":         "True",
	}
	if insta.Account != nil {
		req.Header.Set("Ig-Intended-User-Id", strconv.Itoa(int(insta.Account.ID)))
	} else {
		req.Header.Set("Ig-Intended-User-Id", "0")
	}
	if contentEncoding != "" {
		headers["Content-Encoding"] = contentEncoding
	}

	setHeaders(headers)
	setHeaders(o.ExtraHeaders)
	setHeaders(insta.headerOptions)

	resp, err := insta.c.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"Status: '%s', Status Code: '%d', Err: '%v'",
			resp.Status,
			resp.StatusCode,
			err,
		)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		err = isError(resp.StatusCode, body)
	}
	insta.extractHeaders(resp.Header)

	// Decode gzip encoded responses
	encoding := resp.Header.Get("Content-Encoding")
	if encoding != "" && encoding == "gzip" {
		buf := bytes.NewBuffer(body)
		zr, err := gzip.NewReader(buf)
		if err != nil {
			return nil, nil, err
		}
		body, err = ioutil.ReadAll(zr)
		if err != nil {
			return nil, nil, err
		}
		if err := zr.Close(); err != nil {
			return nil, nil, err
		}
	}
	return body, resp.Header.Clone(), err
}

func (insta *Instagram) checkXmidExpiry() {
	if insta.xmidExpiry != -1 && time.Now().Unix() > insta.xmidExpiry-10 {
		insta.xmidExpiry = -1
		insta.zrToken()
	}
}

func (insta *Instagram) extractHeaders(h http.Header) {
	extract := func(in string, out string) {
		x := h[in]
		if len(x) > 0 && x[0] != "" {
			// prevent from auth being set without token post login
			if in == "Ig-Set-Authorization" && len(insta.headerOptions[out]) != 0 {
				current := strings.Split(insta.headerOptions[out], ":")
				newHeader := strings.Split(x[0], ":")
				if len(current[2]) > len(newHeader[2]) {
					return
				}
			}
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
	case 202:
	case 400:
		ierr := Error400{}
		err = json.Unmarshal(body, &ierr)
		if err == nil {
			if ierr.Message == "challenge_required" {
				return ierr.ChallengeError
			}

			if err == nil {
				return ierr
			}
		}

	case 429:
		return ErrTooManyRequests
	case 503:
		return Error503{
			Message: "Instagram API error. Try it later.",
		}
	default:
		ierr := ErrorN{}
		err = json.Unmarshal(body, &ierr)
		if err != nil {
			return err
		}
		if ierr.Message == "Transcode not finished yet." {
			return nil
		}
		return ierr
	}
	return nil
}

func (insta *Instagram) prepareData(other ...map[string]interface{}) (string, error) {
	data := map[string]interface{}{
		"_uuid": insta.uuid,
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
		return string(b), err
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

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
