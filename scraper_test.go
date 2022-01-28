package tiktok_go

import "testing"

func TestChromedpScraper_Text(t *testing.T) {
	s, err := NewScraper()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	actual, err := s.Text("https://ankiweb.net/account/login", "main.container > h1:nth-child(1)")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	expected := "Log in"
	if actual != expected {
		t.Logf("Expected: %s| Got: %s\n", expected, actual)
		t.Fail()
	}
}
