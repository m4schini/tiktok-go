package scraper

import "testing"

func TestHttpScraper_Text(t *testing.T) {
	//http://gabenewell.org/
	scr, err := NewHttpScraper()
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
