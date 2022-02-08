package tiktok_go

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
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

type Account struct {
	Username    string
	DisplayName string
	Bio         string
	Following   int
	Followers   int
	Likes       int
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
func (a *Account) GetLatestVideoURLs(scr Scraper) ([]string, []int, error) {
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
func GetAccountByUsername(scr Scraper, username string) (*Account, error) {

	account := Account{
		Username: username,
	}

	var err error
	url := account.URL()

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 404 {
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
	URL       string
	ID        string
	Username  string
	Timestamp string

	ViewCount       int
	LikeCount       int
	CommentCount    int
	ShareCount      int
	Audio           string
	VideoLength     int
	Description     string
	DescriptionHTML string
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

func GetVideo(scr Scraper, username, id string) (*Video, error) {
	return GetVideoByUrl(scr, fmt.Sprintf(urlFormatPost, username, id))
}

func GetVideoByUrl(scr Scraper, url string) (*Video, error) {
	var err error

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

	post := Video{
		URL:      url,
		ID:       id,
		Username: username,
	}

	likeCountText, err := scr.Text(url, selPostLikeCount)
	post.LikeCount = toNumber(likeCountText)
	if err != nil {
		return nil, err
	}

	commentCountText, err := scr.Text(url, selPostCommentCount)
	post.CommentCount = toNumber(commentCountText)
	if err != nil {
		return nil, err
	}

	shareCountText, err := scr.Text(url, selPostShareCount)
	post.ShareCount = toNumber(shareCountText)
	if err != nil {
		return nil, err
	}

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
