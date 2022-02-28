package tiktok_go

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/m4schini/tiktok-go/scraper"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func TestGetAccountByUsername_chromedp(t *testing.T) {
	scr, err := scraper.NewChromedpScraper()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	acc, err := GetAccountByUsername(scr, "fabiola.baglieri")
	if err != nil {
		t.Log("err wasn't expected")
		t.Fail()
	}
	t.Log(acc)

	acc, err = GetAccountByUsername(scr, "sdmlpfsdpgmpm")
	if err == nil {
		t.Log("err was expected")
		t.Fail()
	}
	t.Log(acc)

}

func TestGetVideoByUrl_chromedp(t *testing.T) {
	scr, err := scraper.NewChromedpScraper()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	v, err := GetVideoByUrl(scr, "https://www.tiktok.com/@krawallklara/video/7021545676710432006")
	t.Log(v)
}

func TestGetCookies(t *testing.T) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := chromedp.Run(ctx,
		QrCode(),
	)

	t.Log(err)
}

func TestPrintData(t *testing.T) {
	fmt.Println("data:image/png;base64,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")
}

func TestSetCookies(t *testing.T) {
	cookieData := `[
		{"name":"tt_csrf_token","value":"SSbWYirajGiGuk1mhEFaiIaC","domain":".tiktok.com","path":"/","expires":1.646232056920446e+09,"size":37,"httpOnly":false,"secure":true,"session":false,"sameSite":"None","priority":"Medium","sameParty":false,"sourceScheme":"Secure","sourcePort":443},
{"name":"_abck","value":"B08AC327C26291C9CBCE2AECB2FA1B23~-1~YAAQC80QAhM3zAR/AQAAOyCUFwdXvswDeyGTgD6FGGmIsNe5pagxpsR37Vp5fsI1wCb5nUYIYPa0Vd/NY2FaV397H9ZyRRXxMuCe0U/biJ1F1mjFJ90hmNF4Su4gahgN8bIkXujyj3buAjzcSqL/R4Vi05ScHyTsQwyp2XBHtspvrGeCmYJrrfRwPV4S3z+smfRGf+9j8KSoxmvg45souuul3wjsMNgG4/GWtV24etz5wUh1Y/pywoaN7QF9MQIjeDpc45C0g9PGjWxCwmHA+PkNOxIVI22dK6+43udAuLxJuQelZ2Lu4vxR6Wxu4D07IN5mj4C9dYjPIFmRMgcf/NOgRtkS1vJja7sMBdPjbCBYeP31EjfofEQ=~-1~-1~-1","domain":".tiktok.com","path":"/","expires":1.676904057920972e+09,"size":426,"httpOnly":false,"secure":true,"session":false,"priority":"Medium","sameParty":false,"sourceScheme":"Secure","sourcePort":443},
{"name":"bm_sz","value":"78CC50AA88E5CDEDC69B3D3A533D4938~YAAQC80QAhQ3zAR/AQAAOyCUFw4qVxDzGOT7dwCIyhKC7tUn5mlWEk72xlgRcu5L5faSZi4ltAJUFRFU7D7ECMuOs8tEcsIWd2eUDAD/7tBPdiDm3UNFmf9vGtIDuovz1B2wi8ZolPawuON9JaxeJu2jizlXO8uttK3KWLVWYCtg2QpqtQN0qSgE+tE1ztaJqIJUL3FolCEC7FfDHQExaQCavqv+BymKLlX+OIVguvVfFU9pkrPLxE+qeVQC14gDRIQtymuca2w7z+TOzRkmK0SxHLGM8o58rmf7MyoPRPuqU+Y=~3749688~3750212","domain":".tiktok.com","path":"/","expires":1.64538245692104e+09,"size":358,"httpOnly":false,"secure":false,"session":false,"priority":"Medium","sameParty":false,"sourceScheme":"Secure","sourcePort":443},
{"name":"s_v_web_id","value":"verify_e0564587a0e1cc481f701d69742f150b","domain":".tiktok.com","path":"/","expires":-1,"size":49,"httpOnly":false,"secure":true,"session":true,"sameSite":"None","priority":"Medium","sameParty":false,"sourceScheme":"Secure","sourcePort":443},
{"name":"MONITOR_WEB_ID","value":"32d4c886-2870-4425-869d-57874ae5515b","domain":".mon-va.byteoversea.com","path":"/","expires":1.653144061758717e+09,"size":50,"httpOnly":false,"secure":true,"session":false,"sameSite":"None","priority":"Medium","sameParty":false,"sourceScheme":"Secure","sourcePort":443},
{"name":"ttwid","value":"1%7C5jfTLSX_zNOTbZsSHpVtwo-rn7aSolRaxxKUyXGt0js%7C1645368062%7C279502cd2617d36763b06d3dc2014d422380e8464db1d3094ec02f2b3408fff2","domain":".tiktok.com","path":"/","expires":1.67690406210659e+09,"size":132,"httpOnly":true,"secure":true,"session":false,"priority":"Medium","sameParty":false,"sourceScheme":"Secure","sourcePort":443},
{"name":"passport_csrf_token","value":"b8bb3cb8a44d29a0fe1bd145669f1610","domain":".tiktok.com","path":"/","expires":1.65055206254677e+09,"size":51,"httpOnly":false,"secure":true,"session":false,"sameSite":"None","priority":"Medium","sameParty":false,"sourceScheme":"Secure","sourcePort":443},
{"name":"passport_csrf_token_default","value":"b8bb3cb8a44d29a0fe1bd145669f1610","domain":".tiktok.com","path":"/","expires":1.65055206254688e+09,"size":59,"httpOnly":false,"secure":false,"session":false,"priority":"Medium","sameParty":false,"sourceScheme":"Secure","sourcePort":443},
{"name":"msToken","value":"qmXw2N8-i01gCSx6m6N0PV1q4XGkGu8AFSIAsEs7yJNGDWHvinkxWtbQp9DGFk6eUShwihcZsN8ILInKOB4vxJKHc65yRSGWg5RVpiXMbNUiBrM4SmdU_CTE6LNbpk_m","domain":".tiktok.com","path":"/","expires":1.646232063860338e+09,"size":135,"httpOnly":false,"secure":true,"session":false,"sameSite":"None","priority":"Medium","sameParty":false,"sourceScheme":"Secure","sourcePort":443},
{"name":"msToken","value":"fh6jkn-Giv0L-bG8mcvirdsu0P_3IHGWBq2_Wv_IusWJFyWigwCtDE2O6qQlZZ_bIzmXTjaW3F5WDTpYFjUbbEZ5bqFdwRNEqmsb44BoChxGg77DVN7mAPzzBmnAXMvj6Mb5bVVCtfThOj29IoOt1hw=","domain":".tiktokv.com","path":"/","expires":1.646232065259609e+09,"size":159,"httpOnly":false,"secure":true,"session":false,"sameSite":"None","priority":"Medium","sameParty":false,"sourceScheme":"Secure","sourcePort":443},
{"name":"msToken","value":"fh6jkn-Giv0L-bG8mcvirdsu0P_3IHGWBq2_Wv_IusWJFyWigwCtDE2O6qQlZZ_bIzmXTjaW3F5WDTpYFjUbbEZ5bqFdwRNEqmsb44BoChxGg77DVN7mAPzzBmnAXMvj6Mb5bVVCtfThOj29IoOt1hw=","domain":"www.tiktok.com","path":"/","expires":1.653144065e+09,"size":159,"httpOnly":false,"secure":false,"session":false,"priority":"Medium","sameParty":false,"sourceScheme":"Secure","sourcePort":443}
	]`

	var cookies []network.Cookie
	err := json.Unmarshal([]byte(cookieData), &cookies)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(len(cookies))
	for _, cookie := range cookies {
		jsonStr, err := cookie.MarshalJSON()
		if err != nil {
			t.Log(err)
		} else {
			fmt.Printf("%s,\n", string(jsonStr))
		}
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	i := 0
	capture := func() chromedp.Tasks {
		var buffer []byte

		return chromedp.Tasks{
			chromedp.FullScreenshot(&buffer, 100),
			chromedp.ActionFunc(func(ctx context.Context) error {
				if err := ioutil.WriteFile(fmt.Sprintf("captures/c%d.png", i), buffer, 0o644); err != nil {
					log.Println(err)
				}

				log.Println(i)
				i = i + 1
				return nil
			}),
		}
	}

	err = chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {

			// add cookies to chrome
			for _, cookie := range cookies {
				// create cookie expiration
				expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))

				err := network.SetCookie(cookie.Name, cookie.Value).
					WithDomain(cookie.Domain).
					WithPath(cookie.Path).
					WithExpires(&expr).
					WithHTTPOnly(cookie.HTTPOnly).
					WithSecure(cookie.Secure).
					WithSameSite(cookie.SameSite).
					WithPriority(cookie.Priority).
					WithSameParty(cookie.SameParty).
					WithSourceScheme(cookie.SourceScheme).
					WithSourcePort(cookie.SourcePort).
					Do(ctx)
				if err != nil {
					log.Println(err)
				}
			}

			return nil
		}),
		chromedp.Navigate("https://www.tiktok.com/"),
		capture(),
		chromedp.Sleep(3*time.Second),
		capture(),
		chromedp.Sleep(3*time.Second),
		capture(),
	)

	if err != nil {
		t.Fatal(err)
	}
}
