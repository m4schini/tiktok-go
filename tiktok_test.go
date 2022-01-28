package tiktok_go

import "testing"

func TestGetAccountByUsername(t *testing.T) {
	scr, err := NewScraper()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	_, err = GetAccountByUsername(scr, "fabiola.baglieri")
	if err != nil {
		t.Log("err wasn't expected")
		t.Fail()
	}

	_, err = GetAccountByUsername(scr, "sdmlpfsdpgmpm")
	if err == nil {
		t.Log("err was expected")
		t.Fail()
	}
}

func TestGetVideoByUrl(t *testing.T) {
	scr, err := NewScraper()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	v, err := GetVideoByUrl(scr, "https://www.tiktok.com/@krawallklara/video/7021545676710432006")
	t.Log(v)
}
