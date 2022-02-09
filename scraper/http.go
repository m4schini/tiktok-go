package scraper

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
)

const (
	userAgent  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"
	accept     = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"
	acceptLang = "en-US,en;q=0.5"
)

type httpScraper struct {
	client   *http.Client
	location struct {
		url string
		doc *goquery.Document
	}
}

func NewHttpScraper() (*httpScraper, error) {
	newGetter := httpScraper{
		client: &http.Client{},
	}

	return &newGetter, nil
}

func (g *httpScraper) getDoc(url string) (*goquery.Document, error) {
	if g.location.url == url {
		return g.location.doc, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", accept)
	req.Header.Set("Accept-Language", acceptLang)
	log.Println(req.Header)

	resp, err := g.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//if resp.StatusCode != 200 {
	//	return nil, errors.New("request unsuccessful")
	//}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	g.location.url = url
	g.location.doc = doc
	return doc, nil
}

func (g *httpScraper) Text(url string, selector interface{}) (string, error) {
	doc, err := g.getDoc(url)
	if err != nil {
		return "", err
	}

	log.Println(doc.Html())

	str, ok := selector.(string)
	if !ok {
		return "", errors.New("selector must have type string")
	}

	text := doc.Find(str).First().Text()
	return text, nil
}

func (g *httpScraper) InnerHTML(url string, selector interface{}) (string, error) {
	doc, err := g.getDoc(url)
	if err != nil {
		return "", err
	}

	str, ok := selector.(string)
	if !ok {
		return "", errors.New("selector must have type string")
	}

	html, err := doc.Find(str).First().Html()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (g *httpScraper) Number(url string, selector interface{}) (int, error) {
	doc, err := g.getDoc(url)
	if err != nil {
		return -1, err
	}

	str, ok := selector.(string)
	if !ok {
		return -1, errors.New("selector must have type string")
	}

	txt := doc.Find(str).First().Text()
	// avoid to number from tiktok?
	return strconv.Atoi(txt)
}

func (g *httpScraper) HTML(url string) (string, error) {
	doc, err := g.getDoc(url)
	if err != nil {
		return "", err
	}

	html, err := doc.Html()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (g *httpScraper) Close() {

}
