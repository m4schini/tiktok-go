package tiktok_go

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const (
	scraperContextTimeout = 20 * time.Second
)

type Scraper interface {
	Text(url string, selector interface{}) (string, error)
	InnerHTML(url string, selector interface{}) (string, error)
	Number(url string, selector interface{}) (int, error)
	Contains(url string, selector interface{}, text string) (bool, error)
	HTML(url string) (string, error)
	ScrollDown(url string) error
	Close()
}

type ChromeScraper interface {
	Scraper
	Nodes(url string, sel interface{}) ([]*cdp.Node, error)
}

//GoqueryScraper DEPRECATED
type GoqueryScraper interface {
	Scraper
	Nodes(url string, sel interface{}) ([]*goquery.Selection, error)
}

func GetRemoteDebugURL(host string) string {
	resp, err := http.Get(fmt.Sprintf("http://%s:9222/json/version", host))
	if err != nil {
		log.Println(err)
		return ""
	}

	var result map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
		return ""
	}

	return result["webSocketDebuggerUrl"].(string)
}

type chromedpScraper struct {
	ctx             context.Context
	close           func()
	renewFunc       func() (ctx context.Context, cancelFunc context.CancelFunc)
	currentLocation string
}

func NewScraper() (*chromedpScraper, error) {
	allocCtx, _ := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	ctx, cancel := context.WithTimeout(allocCtx, scraperContextTimeout)

	return &chromedpScraper{
		ctx:             ctx,
		close:           cancel,
		currentLocation: "",
	}, nil
}

func NewRemoteScraper(chromedpHost string) (*chromedpScraper, error) {
	allocCtx, _ := chromedp.NewRemoteAllocator(
		context.Background(),
		GetRemoteDebugURL(chromedpHost),
	)

	ctx, cancel := chromedp.NewContext(allocCtx)

	return &chromedpScraper{
		ctx:             ctx,
		close:           cancel,
		currentLocation: "",
	}, nil
}

func (c *chromedpScraper) setLocation(location string) (bool, error) {
	if c.currentLocation != location {

		err := chromedp.Run(c.ctx,
			chromedp.Navigate(location),
		)
		if err != nil {
			return false, err
		}

		c.currentLocation = location
		return true, nil
	} else {
		return false, nil
	}
}

func (c *chromedpScraper) Text(url string, selector interface{}) (string, error) {
	_, err := c.setLocation(url)
	if err != nil {
		return "", err
	}

	var out string
	err = chromedp.Run(c.ctx,
		chromedp.Text(selector, &out),
	)
	if err != nil {
		return "nil", err
	}

	return out, nil
}

func (c *chromedpScraper) InnerHTML(url string, selector interface{}) (string, error) {
	_, err := c.setLocation(url)
	if err != nil {
		return "", err
	}

	var out string
	err = chromedp.Run(c.ctx,
		chromedp.InnerHTML(selector, &out),
	)
	if err != nil {
		return "nil", err
	}

	return out, nil
}

func (c *chromedpScraper) Number(url string, selector interface{}) (int, error) {
	_, err := c.setLocation(url)
	if err != nil {
		return 0, err
	}

	text, err := c.Text(url, selector)
	if err != nil {
		return -1, err
	}

	i, err := strconv.Atoi(text)
	if err != nil {
		return -1, err
	}

	return i, nil
}

func (c *chromedpScraper) Contains(url string, selector interface{}, expect string) (bool, error) {
	actual, err := c.Text(url, selector)
	if err != nil {
		return false, err
	}

	if actual == expect {
		return true, nil
	} else {
		return false, nil
	}
}

func (c *chromedpScraper) HTML(url string) (string, error) {
	_, err := c.setLocation(url)
	if err != nil {
		return "", err
	}

	var out string
	err = chromedp.Run(c.ctx,
		chromedp.OuterHTML("body", &out, chromedp.ByQuery),
	)

	if err != nil {
		return "<html></html>", err
	}

	return out, nil
}

func (c *chromedpScraper) ScrollDown(url string) error {
	_, err := c.setLocation(url)
	if err != nil {
		return err
	}

	var x string
	err = chromedp.Run(c.ctx,
		chromedp.EvaluateAsDevTools("window.scrollBy(0,5000)", &x),
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *chromedpScraper) Nodes(url string, sel interface{}) ([]*goquery.Selection, error) {
	selector, ok := sel.(string)
	if !ok {
		return nil, errors.New("sel must be of type string")
	}

	html, err := c.HTML(url)
	if err != nil {
		return nil, err
	}

	htmlReader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(htmlReader)

	nodes := make([]*goquery.Selection, 0)

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		nodes = append(nodes, s)

		// For each item found, get the title
		num := s.Text()
		fmt.Printf("%d: %s\n", i, num)
	})

	return nodes, nil
}

func (c *chromedpScraper) Close() {
	c.close()
}
