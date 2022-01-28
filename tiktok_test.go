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
