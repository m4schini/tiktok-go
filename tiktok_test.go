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

func TestTiktok_GetVideoByUrl(t *testing.T) {
	tt := NewTikTok()

	vid, err := tt.GetVideoByUrl("https://www.tiktok.com/@krawallklara/video/7021545676710432006")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(vid)
}
