package tiktok_go

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"
	referer   = "https://www.tiktok.com/"
)

type getter struct {
	client   *http.Client
	location struct {
		url string
		doc *goquery.Document
	}
}

func NewGetter() (*getter, error) {
	newGetter := getter{
		client: http.DefaultClient,
	}

	return &newGetter, nil
}

func NewTorGetter(port string) (*getter, error) {
	// Create a transport that uses Tor Browser's SocksPort.  If
	// talking to a system tor, this may be an AF_UNIX socket, or
	// 127.0.0.1:9050 instead.
	tbProxyURL, err := url.Parse("socks5://127.0.0.1:" + port)
	if err != nil {
		return nil, err
	}

	// Get a proxy Dialer that will create the connection on our
	// behalf via the SOCKS5 proxy.  Specify the authentication
	// and re-create the dialer/transport/client if tor's
	// IsolateSOCKSAuth is needed.
	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil {
		return nil, err
	}

	// Make a http.Transport that uses the proxy dialer, and a
	// http.Client that uses the transport.
	tbTransport := &http.Transport{Dial: tbDialer.Dial}
	client := &http.Client{Transport: tbTransport}

	// Example: Fetch something.  Real code will probably want to use
	// client.Do() so they can change the User-Agent.
	//resp, err := client.Get("http://check.torproject.org")
	//if err != nil {
	//	return nil, err
	//}
	//defer resp.Body.Close()

	newGetter := getter{
		client: client,
	}

	return &newGetter, nil
}

func (g *getter) NewCircuit() {
	conn, _ := net.Dial("tcp", "127.0.0.1:9151")
	fmt.Fprintf(conn, ""+"\n")
}

func (g *getter) getDoc(url string) (*goquery.Document, error) {
	if g.location.url == url {
		return g.location.doc, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("user-agent", userAgent)
	req.Header.Set("referer", referer)

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

func (g *getter) Text(url string, selector interface{}) (string, error) {
	doc, err := g.getDoc(url)
	if err != nil {
		return "", err
	}

	str, ok := selector.(string)
	if !ok {
		return "", errors.New("selector must have type string")
	}

	text := doc.Find(str).First().Text()
	return text, nil
}

func (g *getter) InnerHTML(url string, selector interface{}) (string, error) {
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

func (g *getter) Number(url string, selector interface{}) (int, error) {
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
	return toNumber(txt), err
}

func (g *getter) Contains(url string, selector interface{}, text string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (g *getter) HTML(url string) (string, error) {
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

func (g *getter) ScrollDown(url string) error {
	//TODO implement me
	panic("implement me")
}

func (g *getter) Close() {
	//TODO implement me
	panic("implement me")
}
