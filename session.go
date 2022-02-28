package tiktok_go

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func strToPng(b64 string) {
	log.Println("STR TO PNG")
	log.Println(b64)
	unbased, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic("Cannot decode b64")
	}

	r := bytes.NewReader(unbased)
	im, err := png.Decode(r)
	if err != nil {
		panic("Bad png")
	}

	f, err := os.OpenFile("qrCodeFromBase64.png", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic("Cannot open file")
	}

	png.Encode(f, im)
	f.Sync()
	f.Close()
}

func QrCode() chromedp.Tasks {
	var qrcode []byte
	var page []byte
	var html string

	i := 0

	printCookies := func(ctx context.Context) error {
		log.Println("reading network cookies")
		cookies, err := network.GetAllCookies().Do(ctx)
		if err != nil {
			return err
		}

		for _, cookie := range cookies {
			json, err := cookie.MarshalJSON()
			if err != nil {
				log.Println(err)
			} else {
				fmt.Printf("%s,\n", string(json))
			}
		}

		return nil
	}

	capture := func() chromedp.Tasks {
		var buffer []byte

		return chromedp.Tasks{
			chromedp.FullScreenshot(&buffer, 100),
			chromedp.ActionFunc(func(ctx context.Context) error {
				if err := ioutil.WriteFile(fmt.Sprintf("captures/capture%d.png", i), buffer, 0o644); err != nil {
					log.Println(err)
				}

				log.Println(i)
				i = i + 1
				return nil
			}),
		}
	}

	var ok bool

	return chromedp.Tasks{
		chromedp.Navigate("https://www.tiktok.com"),
		capture(),
		chromedp.WaitVisible("a[href=\"https://www.tiktok.com/legal/cookie-settings?lang=en\"] + button"),
		capture(),
		chromedp.Click("a[href=\"https://www.tiktok.com/legal/cookie-settings?lang=en\"] + button"),
		capture(),
		chromedp.WaitNotVisible("a[href=\"https://www.tiktok.com/legal/cookie-settings?lang=en\"] + button"),
		capture(),
		chromedp.Navigate("https://www.tiktok.com/login/qrcode?lang=en&redirect_url=https%3A%2F%2Fwww.tiktok.com%2F"),
		//chromedp.EmulateViewport(1920, 1080),
		chromedp.Sleep(2500 * time.Millisecond),
		capture(),
		chromedp.Click("img[alt=qrcode]"),
		capture(),
		chromedp.Sleep(1 * time.Second),
		chromedp.Screenshot("[alt=qrcode]", &qrcode),
		chromedp.FullScreenshot(&page, 100),
		chromedp.AttributeValue("[alt=qrcode]", "src", &html, &ok),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if err := ioutil.WriteFile("qrcode.png", qrcode, 0o644); err != nil {
				log.Fatal(err)
			}
			if err := ioutil.WriteFile("page.png", page, 0o644); err != nil {
				log.Fatal(err)
			}
			fmt.Println("=================================")
			fmt.Println(ok)
			fmt.Println(html)
			fmt.Println("=================================")

			if ok {
				start := strings.Index(html, ",") + 1

				strToPng(html[start:])
			}

			return nil
		}),
		capture(),
		chromedp.Sleep(60 * time.Second),
		capture(),
		chromedp.Navigate("https://www.tiktok.com"),
		capture(),
		chromedp.Sleep(30 * time.Second),
		capture(),
		chromedp.ActionFunc(printCookies),
	}
}

func SplitCookies(docCookies string) []string {
	// 'tt_csrf_token=6TXVKD_7BfxO2Gw2-G68K10k; _abck=911D90C52E205CBDD90B9BFCAD29697A~-1~YAAQHM0QAoAgFAp/AQAAFFkMDAfoyAZ0JlMAEPHC+ofOXrBlHg5DHkveXUJ7YIZBEyqklx8d3m2nFK9rJj0APZ0lk+uNSaNtqZBL2RF5UZ4HaYGAqpPA8C8DwY9AnNvaNJlJiwDPDG5pmmFgmiKXfV2TLoX1gsSwokgokNDALXC5rYvnMFGEedXT3R2JzQZukH17wJRu1YgGEQxkegjEStaQhBP6zF8I3r3L+51OWHcvjhcM/awN1U65bvdKVm5cxHZN348dSvVuwD1T/WcF8ZLvTSVMoOS6yyWJH+dVFzgHAKrPg4ODKXzYPqueMtsQJg0qNf+ipqvmRe+VHCLM7KHX4+QqmajYY31D5tNzY1VK/owpH+IovAI59A==~-1~-1~-1; bm_sz=E54ECF6F8F2CC7A5A5ADD82F4DFACF8B~YAAQHM0QAoIgFAp/AQAAFFkMDA7FeaV1UurIhG32Brh9Tql4lxB+bUa1a2Hnq3rzlUk9hjUDdohA2Zv/wMz9zDZuUI6fVMVWIkhP+9xxtFYXcQgcbnDKsfzSyFMec01r4GCGDAl/pG77ZRBm3gE31psfo2LMQcaxgq3EaG+Wbkg3jtPHrXk+GLO0irfedT73iEZ2F1QlUnBU5YbcQ/4eqJYUSrfVXJxLutwZ3y8+XnLHiNKs8RX9bGY3+oGG36qz6WDGwJinwPxfKlGWQ2JxjsaMYR79KJDXa0ieC1h5kwosdzU=~3289410~3290679; passport_csrf_token=aaf33dd863320c8bd307bcab4a1a4588; passport_csrf_token_default=aaf33dd863320c8bd307bcab4a1a4588; s_v_web_id=verify_17e7cc821ca81cc4456edbe51f8ac8fc; passport_fe_beating_status=true; msToken=JAzzNHDUUqEwtnqGyIRo8RwOqahYuJDo4h5Zp31MatfIaNP4B0xrca6YSuAYJZTZeAFnQz5uQMynTryB-eQSieWg6ENvPfPBlPQPdJS0HBNMxi10iwN5xQH_9g4lKGfaBBWsibYb; msToken=TGqDWe7Z1UhiIf8_3xSggnn7TMnrcxJ1gcN6tiJ5ct-tKILr-Tfkh5JWDrKG6EefzzdgtwoyBgjXXoOMfsQ0iOMACZNkzOZJuSjfoDXjXt7iO4C9HQbxaxylIZO8h1Euu7oB9BNc0J5HMQktmw=='

	result := make([]string, 0)

	cookies := strings.Split(docCookies, ";")
	for _, cookie := range cookies {
		div := strings.Index(cookie, "=")
		key := cookie[:div]
		val := cookie[div+1:]

		result = append(result, key, val)
		fmt.Println(key, val)
	}

	return result
}

type CP struct {
	Secure   bool
	HttpOnly bool
}

func Props() map[string]CP {
	cm := make(map[string]CP)
	cm["msToken"] = CP{
		Secure:   true,
		HttpOnly: false,
	}
	cm["passport_fe_beating_status"] = CP{
		Secure:   false,
		HttpOnly: false,
	}
	cm["tt-target-idc"] = CP{
		Secure:   false,
		HttpOnly: true,
	}
	cm["store-idc"] = CP{
		Secure:   false,
		HttpOnly: true,
	}
	cm["ssid_ucp_v1"] = CP{
		Secure:   true,
		HttpOnly: true,
	}
	cm["sessionid_ss"] = CP{
		Secure:   true,
		HttpOnly: true,
	}
	cm["odin_tt"] = CP{
		Secure:   false,
		HttpOnly: true,
	}
	cm["sid_ucp_v1"] = CP{
		Secure:   true,
		HttpOnly: true,
	}
	cm["tt_csrf_token"] = CP{
		Secure:   true,
		HttpOnly: false,
	}
	cm["bm_sv"] = CP{
		Secure:   true,
		HttpOnly: false,
	}
	cm["store-country-code"] = CP{
		Secure:   false,
		HttpOnly: true,
	}
	cm["_abck"] = CP{
		Secure:   true,
		HttpOnly: false,
	}
	cm["ak_bmsc"] = CP{
		Secure:   false,
		HttpOnly: true,
	}
	cm["passport_csrf_token"] = CP{
		Secure:   true,
		HttpOnly: false,
	}
	cm["s_v_web_id"] = CP{
		Secure:   true,
		HttpOnly: false,
	}
	cm["bm_sz"] = CP{
		Secure:   false,
		HttpOnly: false,
	}
	cm["sid_tt"] = CP{
		Secure:   false,
		HttpOnly: true,
	}
	cm["passport_csrf_token_default"] = CP{
		Secure:   false,
		HttpOnly: false,
	}
	cm["cmpl_token"] = CP{
		Secure:   true,
		HttpOnly: true,
	}
	cm["ttwid"] = CP{
		Secure:   true,
		HttpOnly: true,
	}
	cm["sid_guard"] = CP{
		Secure:   false,
		HttpOnly: true,
	}
	cm["uid_tt"] = CP{
		Secure:   false,
		HttpOnly: true,
	}
	cm["uid_tt_ss"] = CP{
		Secure:   true,
		HttpOnly: true,
	}
	cm["bm_mi"] = CP{
		Secure:   false,
		HttpOnly: true,
	}
	cm["sessionid"] = CP{
		Secure:   false,
		HttpOnly: true,
	}
	return cm
}

func SetCookiesTask(host string, documentCookie string) chromedp.Tasks {
	cookies := SplitCookies(documentCookie)

	props := Props()

	return chromedp.Tasks{
		// navigate to site
		chromedp.Navigate(host),

		chromedp.ActionFunc(func(ctx context.Context) error {
			// create cookie expiration
			expr := cdp.TimeSinceEpoch(time.Now().Add(10 * time.Minute))
			// add cookies to chrome
			for i := 0; i < len(cookies); i += 2 {
				err := network.SetCookie(cookies[i], cookies[i+1]).
					WithExpires(&expr).
					WithDomain("tiktok.com").
					WithPath("/").
					WithSecure(props[cookies[i]].Secure).
					WithHTTPOnly(props[cookies[i]].HttpOnly).
					Do(ctx)
				if err != nil {
					log.Println("[ERR]", err, i, cookies[i], cookies[i+1])
					//return err
				}
			}

			log.Println("Cookies set")
			return nil
		}),
		// read network values
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("reading network cookies")
			cookies, err := network.GetAllCookies().Do(ctx)
			if err != nil {
				return err
			}

			for i, cookie := range cookies {
				log.Printf("chrome cookie %d: %+v", i, cookie)
			}

			return nil
		}),
	}
}
