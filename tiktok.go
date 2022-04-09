package tiktok

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/m4schini/tiktok-go/model"
	"github.com/m4schini/tiktok-go/util"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
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

type TikTok interface {
	GetAccount(username string) (*model.Account, error)
	GetAccountByUrl(url string) (*model.Account, error)

	GetLatestVideos(username string) ([]*model.VideoPreview, error)
	GetLatestVideosByUrl(url string) ([]*model.VideoPreview, error)

	GetVideo(username, videoId string) (*model.Video, error)
	GetVideoByUrl(url string) (*model.Video, error)
}

type tiktok struct {
	mu sync.Mutex
}

func NewTikTok() *tiktok {
	t := new(tiktok)
	return t
}

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

//GetAccount
func (t *tiktok) GetAccount(username string) (*model.Account, error) {
	scr, err := util.GetScraper()
	if err != nil {
		return nil, err
	}
	defer scr.Close()

	acc := model.Account{
		Username: username,
	}
	url := acc.URL()

	if CheckUrl(url) == 404 {
		return nil, errors.New("account doesn't exist")
	}

	acc.DisplayName, err = scr.Text(url, selUserDisplayName)
	if err != nil {
		return nil, err
	}

	acc.Bio, err = scr.Text(url, selUserBio)
	if err != nil {
		return nil, err
	}

	followingText, err := scr.Text(url, selUserFollowingCount)
	if err != nil {
		return nil, err
	}
	acc.Following = util.ToNumber(followingText)

	followersText, err := scr.Text(url, selUserFollowersCount)
	if err != nil {
		return nil, err
	}
	acc.Followers = util.ToNumber(followersText)

	likesText, err := scr.Text(url, selUserLikesCount)
	if err != nil {
		return nil, err
	}
	acc.Likes = util.ToNumber(likesText)

	acc.AvatarURL, err = scr.Attr(url, `div[data-e2e="user-page"] span img`, "src")
	if err != nil {
		return nil, err
	}

	return &acc, nil
}

func (t *tiktok) GetAccountByUrl(url string) (*model.Account, error) {
	parts := strings.Split(url, "@")
	if len(parts) != 2 {
		return nil, errors.New("malformed url. Expected URL format: https://www.tiktok.com/@<username>")
	}
	return t.GetAccount(parts[1])
}

func (t *tiktok) GetLatestVideos(username string) ([]*model.VideoPreview, error) {
	scr, err := util.GetScraper()
	if err != nil {
		return nil, err
	}
	defer scr.Close()

	// extract rendered html from chromedp
	html, err := scr.HTML(model.ToAccountURL(username))
	if err != nil {
		return nil, err
	}

	// Create goquery html doc from html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	// extract video vids
	vids := make([]*model.VideoPreview, 0)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		t, o := s.Attr("href")

		// if attribute exists -> check if url contains "video"
		// update: also check for "tiktok", to be sure we got an absolute and not a relative url
		// (also to avoid duplication)
		if o && strings.Contains(t, "video") && strings.Contains(t, "http") {
			vids = append(vids, &model.VideoPreview{
				URL: t,
			})
		}
	})

	// extract video views
	l := len(vids)
	doc.Find("[data-e2e=\"video-views\"]").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		t := s.Text()
		//fmt.Printf("%d: %d %s\n", i, o, t)

		if i >= 0 && i < l {
			vids[i].Views = util.ToNumber(t)
		}
	})

	return vids, nil
}

func (t *tiktok) GetLatestVideosByUrl(url string) ([]*model.VideoPreview, error) {
	parts := strings.Split(url, "@")
	if len(parts) != 2 {
		return nil, errors.New("malformed url. Expected URL format: https://www.tiktok.com/@<username>")
	}
	return t.GetLatestVideos("")
}

func (t *tiktok) GetVideo(username, videoId string) (*model.Video, error) {
	return t.GetVideoByUrl(fmt.Sprintf(urlFormatPost, username, videoId))
}

func (t *tiktok) GetVideoByUrl(url string) (*model.Video, error) {
	scr, err := util.GetScraper()
	if err != nil {
		return nil, err
	}
	defer scr.Close()

	username, id := util.ExtractUsernameAndId(url)

	if CheckUrl(url) != 200 {
		return &model.Video{
			URL:       url,
			ID:        id,
			Username:  username,
			Available: false,
		}, nil
	}

	post := model.Video{
		URL:       url,
		ID:        id,
		Username:  username,
		Available: true,
	}

	likeCountText, err := scr.Text(url, selPostLikeCount)
	if err != nil {
		return nil, err
	}
	post.Likes = util.ToNumber(likeCountText)

	commentCountText, err := scr.Text(url, selPostCommentCount)
	if err != nil {
		return nil, err
	}
	post.Comments = util.ToNumber(commentCountText)

	shareCountText, err := scr.Text(url, selPostShareCount)
	if err != nil {
		return nil, err
	}
	post.Shares = util.ToNumber(shareCountText)

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

	post.VideoURL, err = scr.Attr(url, "video", "src")
	if err != nil {
		return nil, err
	}

	post.ThumbnailURL, err = scr.Attr(url, `div[data-e2e="feed-video"] img`, "src")
	if err != nil {
		return nil, err
	}

	return &post, nil
}
