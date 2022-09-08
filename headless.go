package goinsta

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type headlessOptions struct {
	// seconds
	timeout int64

	showBrowser bool

	tasks chromedp.Tasks
}

// Wait until page gets redirected to instagram home page
func waitForInstagram(b *bool) chromedp.ActionFunc {
	return chromedp.ActionFunc(
		func(ctx context.Context) error {
			for {
				select {
				case <-time.After(time.Millisecond * 250):
					var l string
					err := chromedp.Location(&l).Do(ctx)
					if err != nil {
						return err
					}
					if l == "https://www.instagram.com/" {
						*b = true
						return nil
					}
				case <-ctx.Done():
					return nil
				}
			}
		})
}

// Wait until page gets redirected to instagram home page
func printButtons(insta *Instagram) chromedp.Action {
	return chromedp.ActionFunc(
		func(ctx context.Context) error {
			var nodes []*cdp.Node
			err := chromedp.Nodes("button", &nodes, chromedp.ByQueryAll).Do(ctx)
			if err != nil {
				return err
			}
			for _, p := range nodes {
				if len(p.Children) > 0 {
					insta.infoHandler(
						fmt.Sprintf("Found button on challenge page: %s\n",
							p.Children[0].NodeValue,
						))
				}
			}
			return nil
		})
}

func takeScreenshot(fn string) chromedp.Action {
	return chromedp.ActionFunc(
		func(ctx context.Context) error {
			var buf []byte
			err := chromedp.FullScreenshot(&buf, 90).Do(ctx)
			if err != nil {
				return err
			}
			if err := os.WriteFile(fn, buf, 0o644); err != nil {
				return err
			}
			return nil
		})
}

func (insta *Instagram) acceptPrivacyCookies(url string) error {
	// Looks for the "Allow All Cookies button"
	selector := `//button[contains(text(),"Allow All Cookies")]`

	// This value is not actually used, since its headless, the browser cannot
	//  be closed easily. If the process is unsuccessful, it will return a timeout error.
	success := false

	return insta.runHeadless(
		&headlessOptions{
			timeout:     60,
			showBrowser: false,
			tasks: chromedp.Tasks{
				chromedp.Navigate(url),

				// wait a second after elemnt is visible, does not work otherwise
				chromedp.WaitVisible(selector),
				chromedp.Sleep(time.Second * 1),
				chromedp.Click(selector, chromedp.BySearch),

				waitForInstagram(&success),
			},
		},
	)
}

func (insta *Instagram) openChallenge(url string) error {
	fname := fmt.Sprintf("challenge-screenshot-%d.png", time.Now().Unix())

	success := false

	err := insta.runHeadless(
		&headlessOptions{
			timeout:     300,
			showBrowser: true,
			tasks: chromedp.Tasks{
				chromedp.Navigate(url),

				// Wait for a few seconds, and screenshot the page after
				chromedp.Sleep(time.Second * 5),
				printButtons(insta),
				takeScreenshot(fname),

				// Wait until page gets redirected to instagram home page
				waitForInstagram(&success),
			},
		},
	)
	if err != nil {
		return err
	}

	insta.infoHandler(
		fmt.Sprintf(
			"Saved a screenshot of the challenge '%s' to %s, please report it in a github issue so the challenge can be solved automatiaclly.\n",
			url,
			fname,
		))

	if !success {
		return ErrChallengeFailed
	}
	return nil
}

// runHeadless takes a list of chromedp actions to perform, wrapped around default
//   actions that will need to be run for every headless request, such as setting
//   the cookies and user-agent.
func (insta *Instagram) runHeadless(options *headlessOptions) error {
	if insta.privacyCalled.Get() {
		return errors.New("Accept privacy cookie method has already been called. Did it not work? please report on a github issue")
	}

	if options.timeout <= 0 {
		options.timeout = 60
	}

	// Extract required headers as cookies
	cookies := map[string]string{}
	cookie_list := []string{
		"x-mid",
		"authorization",
		"ig-u-shbid",
		"ig-u-shbts",
		"ig-u-ds-user-id",
		"ig-u-rur",
	}
	insta.headerOptions.Range(
		func(key, value interface{}) bool {
			header := strings.ToLower(key.(string))
			for _, cookie_name := range cookie_list {
				if cookie_name == header {
					cookies[cookie_name] = value.(string)
				}
			}
			return true
		},
	)

	userAgent := fmt.Sprintf(
		"Mozilla/5.0 (Linux; Android %d; %s/%s; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/95.0.4638.50 Mobile Safari/537.36 %s",
		insta.device.AndroidRelease,
		insta.device.Model,
		insta.device.Chipset,
		insta.userAgent,
	)

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(userAgent),
	)

	if insta.proxy != "" {
		opts = append(opts, chromedp.ProxyServer(insta.proxy))
	}
	if insta.proxyInsecure {
		opts = append(opts, chromedp.Flag("ignore-certificate-errors", true))
	}
	if options.showBrowser {
		opts = append(opts, chromedp.Flag("headless", false))
	}

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create chrome instance
	ctx, cancel = chromedp.NewContext(
		ctx,
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, time.Duration(options.timeout)*time.Second)
	defer cancel()

	// Size for custom device
	// res := strings.Split(strings.ToLower(insta.device.ScreenResolution), "x")
	// width, err := strconv.Atoi(res[0])
	// if err != nil {
	// 	return err
	// }
	// height, err := strconv.Atoi(res[1])
	// if err != nil {
	// 	return err
	// }

	default_actions := chromedp.Tasks{
		// Set custom device type
		chromedp.Tasks{
			emulation.SetUserAgentOverride(userAgent),
			// emulation.SetDeviceMetricsOverride(int64(width), int64(height), 1.000000, true).
			// 	WithScreenOrientation(&emulation.ScreenOrientation{
			// 		Type:  emulation.OrientationTypePortraitPrimary,
			// 		Angle: 0,
			// 	}),
			emulation.SetTouchEmulationEnabled(true),
		},

		// Set custom cookie
		chromedp.ActionFunc(func(ctx context.Context) error {
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			for key, val := range cookies {
				err := network.SetCookie(key, val).
					WithExpires(&expr).
					WithDomain("i.instagram.com").
					// WithHTTPOnly(true).
					Do(ctx)
				if err != nil {
					return err
				}
			}
			return nil
		}),

		// Set custom headers
		network.Enable(),
		network.SetExtraHTTPHeaders(
			network.Headers(
				map[string]interface{}{
					"X-Requested-With": "com.instagram.android",
				},
			),
		),
	}

	err := chromedp.Run(ctx, append(default_actions, options.tasks))
	return err
}
