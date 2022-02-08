package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
)

type torScraper struct {
	client   *http.Client
	location struct {
		url string
		doc *goquery.Document
	}
}

func NewTorScraper(host string, port int) (*torScraper, error) {
	// Create a transport that uses Tor Browser's SocksPort.  If
	// talking to a system tor, this may be an AF_UNIX socket, or
	// 127.0.0.1:9050 instead.
	tbProxyURL, err := url.Parse(fmt.Sprintf("socks5://%s:%d", host, port))
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

	newGetter := torScraper{
		client: client,
	}

	return &newGetter, nil
}

func (g *torScraper) NewCircuit() {
	conn, _ := net.Dial("tcp", "127.0.0.1:9151")
	fmt.Fprintf(conn, ""+"\n")
}
