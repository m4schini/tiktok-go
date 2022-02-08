package tiktok_go

import (
	"github.com/m4schini/tiktok-go/scraper"
	"testing"
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

func TestGetAccountByUsername_getter(t *testing.T) {
	scr, err := scraper.NewTorScraper()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	acc, err := GetAccountByUsername(scr, "fabiola.baglieri")
	if err != nil {
		t.Log("err wasn't expected:", err)
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
