package scraper

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestChromedpScraper_Text(t *testing.T) {
	s, err := NewChromedpScraper()
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

func TestChromedpScraper_Screenshot(t *testing.T) {
	s, err := NewChromedpScraper()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	screenshot, err := s.Screenshot("https://www.tiktok.com/@krawallklara")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if err := ioutil.WriteFile("fullScreenshot.png", screenshot, 0o644); err != nil {
		log.Fatal(err)
	}
}

func TestChromedpScraper_ScreenshotElement(t *testing.T) {
	s, err := NewChromedpScraper()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	screenshot, err := s.ScreenshotElement("https://www.tiktok.com/@krawallklara",
		"[data-e2e=\"user-avatar\"]")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if err := ioutil.WriteFile("elementScreenshot.png", screenshot, 0o644); err != nil {
		log.Fatal(err)
	}
}

func TestChromedpScraper_Login(t *testing.T) {
	s, err := NewChromedpScraper()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	s.Login("", "")
	screenshot, err := s.Screenshot("https://www.tiktok.com/")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if err := ioutil.WriteFile("login.png", screenshot, 0o644); err != nil {
		log.Fatal(err)
	}
}
