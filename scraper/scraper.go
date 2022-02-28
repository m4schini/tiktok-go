package scraper

import (
	"time"
)

const (
	ScraperContextTimeout = 30 * time.Second
)

type Scraper interface {
	Text(url string, selector interface{}) (string, error)
	InnerHTML(url string, selector interface{}) (string, error)
	Number(url string, selector interface{}) (int, error)
	HTML(url string) (string, error)
	Attr(url string, selector interface{}, attrName string) (string, error)
	Close()
}
