package tiktok

import (
	"testing"
)

func TestTiktok_GetAccountByUrl(t *testing.T) {
	tt := NewTikTok()

	acc, err := tt.GetAccountByUrl("https://www.tiktok.com/@krawallklara2.0")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(acc)
}

func TestTiktok_GetAccount(t *testing.T) {
	tt := NewTikTok()

	acc, err := tt.GetAccount("krawallklara2.0")
	if err != nil {
		t.Fatal(err)
	}

	acc, err = tt.GetAccount("sdmlpfsdpgmpm")
	if err == nil {
		t.Log("Expected error")
		t.Fail()
	}

	t.Log(acc)
}

func TestTiktok_GetVideoByUrl1(t *testing.T) {
	tt := NewTikTok()

	vid, err := tt.GetVideoByUrl("https://www.tiktok.com/@krawallklara/video/7021545676710432006")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(vid)
}

func TestTiktok_GetVideoByUrl2(t *testing.T) {
	tt := NewTikTok()

	vid, err := tt.GetVideoByUrl("https://www.tiktok.com/@krawallklara2.0/video/7005642467865398534")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(vid)
}

func TestTiktok_GetVideoByUrl3(t *testing.T) {
	tt := NewTikTok()

	vid, err := tt.GetVideoByUrl("https://www.tiktok.com/@tagesschau/video/7084270134969273605")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(vid)
}
