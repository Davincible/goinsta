package goinsta

import (
	"encoding/json"
	"image"
	// Required for getImageDimensionFromReader in jpg and png format
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strconv"
	"time"
	"unsafe"
)

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func toString(i interface{}) string {
	switch s := i.(type) {
	case string:
		return s
	case bool:
		return strconv.FormatBool(s)
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32)
	case int:
		return strconv.Itoa(s)
	case int64:
		return strconv.FormatInt(s, 10)
	case int32:
		return strconv.Itoa(int(s))
	case int16:
		return strconv.FormatInt(int64(s), 10)
	case int8:
		return strconv.FormatInt(int64(s), 10)
	case uint:
		return strconv.FormatInt(int64(s), 10)
	case uint64:
		return strconv.FormatInt(int64(s), 10)
	case uint32:
		return strconv.FormatInt(int64(s), 10)
	case uint16:
		return strconv.FormatInt(int64(s), 10)
	case uint8:
		return strconv.FormatInt(int64(s), 10)
	case []byte:
		return b2s(s)
	case error:
		return s.Error()
	}
	return ""
}

func prepareRecipients(cc interface{}) (bb string, err error) {
	var b []byte
	ids := make([][]int64, 0)
	switch c := cc.(type) {
	case *Conversation:
		for i := range c.Users {
			ids = append(ids, []int64{c.Users[i].ID})
		}
	case *Item:
		ids = append(ids, []int64{c.User.ID})
	case int64:
		ids = append(ids, []int64{c})
	}
	b, err = json.Marshal(ids)
	bb = b2s(b)
	return
}

// getImageDimensionFromReader return image dimension , types is .jpg and .png
func getImageDimensionFromReader(rdr io.Reader) (int, int, error) {
	image, _, err := image.DecodeConfig(rdr)
	if err != nil {
		return 0, 0, err
	}
	return image.Width, image.Height, nil
}

func getTimeOffset() string {
	_, offset := time.Now().Zone()
	return strconv.Itoa(offset)
}

func Jazoest(str string) string {
	b := []byte(str)
	var s int
	for v := range b {
		s += v
	}
	return "2" + strconv.Itoa(s)
}

func createUserAgent() string {
	// Instagram 195.0.0.31.123 Android (28/9; 560dpi; 1440x2698; LGE/lge; LG-H870DS; lucye; lucye; en_GB; 302733750)
	// Instagram 195.0.0.31.123 Android (28/9; 560dpi; 1440x2872; Genymotion/Android; Samsung Galaxy S10; vbox86p; vbox86; en_US; 302733773)  # version_code: 302733773
	// Instagram 195.0.0.31.123 Android (30/11; 560dpi; 1440x2898; samsung; SM-G975F; beyond2; exynos9820; en_US; 302733750)
	return fmt.Sprintf("Instagram %s Android (%d/%s; %s; %s; %s; %s; %s; %s; %s; %s)",
		appVersion,
		goInstaDeviceSettings["android_version"],
		goInstaDeviceSettings["android_release"],
		goInstaDeviceSettings["screen_dpi"],
		goInstaDeviceSettings["screen_resolution"],
		goInstaDeviceSettings["manufacturer"],
		goInstaDeviceSettings["model"],
		goInstaDeviceSettings["code_name"],
		goInstaDeviceSettings["chipset"],
		locale,
		appVersionCode,
	)
}
