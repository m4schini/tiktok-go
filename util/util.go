package util

import (
	"github.com/m4schini/tiktok-go/scraper"
	"os"
	"strconv"
	"strings"
)

func ToNumber(in string) int {
	if in == "" {
		return -1
	}
	i, err := strconv.Atoi(in)
	if err != nil {
		unit := rune(in[len(in)-1])

		unitMul := 1.0
		switch unit {
		case 'K':
			unitMul = 1000
		case 'M':
			unitMul = 1000000
		}

		i, err := strconv.ParseFloat(in[0:len(in)-1], 32)
		if err != nil {
			return -1
		}

		return int(i * unitMul)
	}

	return i
}

func ExtractUsernameAndId(url string) (string, string) {
	var username string
	var id string

	parts := strings.Split(url, "/")
	if len(parts) == 4 {
		username = parts[1][1:len(parts[1])]
		id = parts[3]
	} else {
		username = parts[3][1:len(parts[3])]
		id = parts[5]
	}

	return username, id
}

func GetScraper() (scraper.Scraper, error) {
	addr := os.Getenv("TIKTOK_GO_TOR_PROXY_ADDR")
	if addr == "" {
		return scraper.NewChromedpScraper()
	} else {
		return scraper.NewProxyChromedpScraper(addr)
	}
}
