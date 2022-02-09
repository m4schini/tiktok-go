package scraper

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestChromedpScraper_Text(t *testing.T) {
	//http://gabenewell.org/
	scr, err := NewChromedpScraper()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	text, err := scr.Text("http://gabenewell.org/", ".run")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Log(text)

	text, err = scr.Text("https://www.tiktok.com/@krawallklara/video/7062067366166891781",
		"[data-e2e=\"share-count\"]")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Log(text)
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
