package tiktok_go

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/m4schini/tiktok-go/scraper"
	"log"
	"net/http"
	"strconv"
	"strings"
)

import (
	"regexp"
)

const (
	urlTikTok = `https://www.tiktok.com/`
	// parameter: 1. Username
	urlFormatUser = `https://www.tiktok.com/@%s`
	// parameter: 1. Username, 2. Video Id
	urlFormatPost = `https://www.tiktok.com/@%s/video/%s`

	selUserDisplayName    = "[data-e2e=\"user-subtitle\"]"
	selUserBio            = "[data-e2e=\"user-bio\"]"
	selUserFollowingCount = "[data-e2e=\"following-count\"]"
	selUserFollowersCount = "[data-e2e=\"followers-count\"]"
	selUserLikesCount     = "[data-e2e=\"likes-count\"]"
	selUserPostItem       = "[data-e2e=\"user-post-item\"]"

	selPostLikeCount    = "[data-e2e=\"like-count\"]"
	selPostCommentCount = "[data-e2e=\"comment-count\"]"
	selPostShareCount   = "[data-e2e=\"share-count\"]"
	selPostDescription  = "[data-e2e=\"video-desc\"]"
	selPostAudio        = "[data-e2e=\"video-music\"]"
	selPostTimestamp    = "[data-e2e=\"browser-nickname\"]"
)

func CheckUrl(url string) int {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	hc := http.Client{}
	response, err := hc.Do(req)
	if err != nil {
		return -1
	}

	return response.StatusCode
}

func toNumber(in string) int {
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

func ToAccountURL(username string) string {
	return fmt.Sprintf(urlFormatUser, username)
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

type Account struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Bio         string `json:"bio"`
	Following   int    `json:"following"`
	Followers   int    `json:"followers"`
	Likes       int    `json:"likes"`
}

func (a *Account) URL() string {
	if a.Username == "" {
		return urlTikTok
	}
	return ToAccountURL(a.Username)
}

func (a *Account) String() string {
	rx := regexp.MustCompile("\\n")
	olBio := rx.ReplaceAllString(a.Bio, " \\\\ ")

	s := fmt.Sprintf("[%s|%s] Following: %d Followers: %d Likes: %d | Bio: \"%s\"",
		a.Username,
		a.DisplayName,
		a.Following,
		a.Followers,
		a.Likes,
		olBio,
	)

	return s
}

// GetLatestVideoURLs Get the latest video urls and views. Returned as to separate lists.
// Assume that the same index returns the value for the same video in bots slices.
func (a *Account) GetLatestVideoURLs(scr scraper.Scraper) ([]string, []int, error) {

	// extract rendered html from chromedp
	html, err := scr.HTML(a.URL())
	if err != nil {
		return nil, nil, err
	}

	// Create goquery html doc from html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	// extract video urls
	urls := make([]string, 0)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		t, o := s.Attr("href")

		// if attribute exists -> check if url contains "video"
		// update: also check for "tiktok", to be sure we got an absolute and not a relative url
		// (also to avoid duplication)
		if o && strings.Contains(t, "video") && strings.Contains(t, "http") {
			urls = append(urls, t)
		}
	})

	// extract video views
	views := make([]int, 0)
	doc.Find("[data-e2e=\"video-views\"]").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		t := s.Text()
		//fmt.Printf("%d: %d %s\n", i, o, t)

		views = append(views, toNumber(t))
	})

	return urls, views, nil
}

// GetAccountByUsername TODO test what happens for user that doesn't exist
func GetAccountByUsername(scr scraper.Scraper, username string) (*Account, error) {

	account := Account{
		Username: username,
	}

	var err error
	url := account.URL()

	if CheckUrl(url) == 404 {
		return nil, errors.New("account doesn't exist")
	}

	account.DisplayName, err = scr.Text(url, selUserDisplayName)
	if err != nil {
		return nil, err
	}

	account.Bio, err = scr.Text(url, selUserBio)
	if err != nil {
		return nil, err
	}

	followingText, err := scr.Text(url, selUserFollowingCount)
	account.Following = toNumber(followingText)
	if err != nil {
		return nil, err
	}

	followersText, err := scr.Text(url, selUserFollowersCount)
	account.Followers = toNumber(followersText)
	if err != nil {
		return nil, err
	}

	likesText, err := scr.Text(url, selUserLikesCount)
	account.Likes = toNumber(likesText)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

type Video struct {
	Available bool `json:"available"`

	URL       string `json:"URL"`
	VideoURL  string `json:"videoURL"`
	ID        string `json:"ID"`
	Username  string `json:"username"`
	Timestamp string `json:"timestamp"`

	ViewCount       int    `json:"viewCount"`
	LikeCount       int    `json:"likeCount"`
	CommentCount    int    `json:"commentCount"`
	ShareCount      int    `json:"shareCount"`
	Audio           string `json:"audio"`
	VideoLength     int    `json:"videoLength"`
	Description     string `json:"description"`
	DescriptionHTML string `json:"descriptionHTML"`
}

func (p Video) String() string {
	rx := regexp.MustCompile("\\n")
	olDesc := rx.ReplaceAllString(p.Description, " \\\\ ")

	s := fmt.Sprintf("[%s|%s] Likes: %d Comments: %d Shares: %d | Desc: \"%s\"",
		p.Username,
		p.ID,
		p.LikeCount,
		p.CommentCount,
		p.ShareCount,
		olDesc,
	)

	return s
}

func (p Video) GetMentions() []string {
	r := regexp.MustCompile(`/@[^/"]+`)
	matches := r.FindAllString(p.DescriptionHTML, -1)

	return matches
}

func (p Video) GetTags() []string {
	r := regexp.MustCompile(`/tag/\w+`)
	matches := r.FindAllString(p.DescriptionHTML, -1)

	return matches
}

func GetVideo(scr scraper.Scraper, username, id string) (*Video, error) {
	return GetVideoByUrl(scr, fmt.Sprintf(urlFormatPost, username, id))
}

func GetVideoByUrl(scr scraper.Scraper, url string) (*Video, error) {
	var err error

	username, id := ExtractUsernameAndId(url)

	if CheckUrl(url) != 200 {
		return &Video{
			URL:       url,
			ID:        id,
			Username:  username,
			Available: false,
		}, nil
	}

	post := Video{
		URL:       url,
		ID:        id,
		Username:  username,
		Available: true,
	}

	likeCountText, err := scr.Text(url, selPostLikeCount)
	if err != nil {
		return nil, err
	}
	post.LikeCount = toNumber(likeCountText)

	commentCountText, err := scr.Text(url, selPostCommentCount)
	if err != nil {
		return nil, err
	}
	post.CommentCount = toNumber(commentCountText)

	shareCountText, err := scr.Text(url, selPostShareCount)
	if err != nil {
		return nil, err
	}
	post.ShareCount = toNumber(shareCountText)

	post.DescriptionHTML, err = scr.InnerHTML(url, selPostDescription)
	if err != nil {
		return nil, err
	}

	post.Description, err = scr.Text(url, selPostDescription)
	if err != nil {
		return nil, err
	}

	post.Audio, err = scr.Text(url, selPostAudio)
	if err != nil {
		return nil, err
	}

	return &post, nil
}
