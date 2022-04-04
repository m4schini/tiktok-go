package model

import (
	"fmt"
	"regexp"
)

const (
	// parameter: 1. Username, 2. Video Id
	urlFormatVideo = `https://www.tiktok.com/@%s/video/%s`
)

type VideoPreview struct {
	URL   string `json:"URL,omitempty"`
	Views int    `json:"Views"`
}

func ToVideoURL(username, videoId string) string {
	return fmt.Sprintf(urlFormatVideo, username, videoId)
}

type Video struct {
	Available bool `json:"available"`

	URL          string `json:"URL"`
	ID           string `json:"ID"`
	Username     string `json:"username"`
	VideoURL     string `json:"video"`
	Timestamp    string `json:"timestamp"`
	ThumbnailURL string `json:"thumbnail"`

	Views           int    `json:"views"`
	Likes           int    `json:"likes"`
	Comments        int    `json:"comments"`
	Shares          int    `json:"shares"`
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
		p.Likes,
		p.Comments,
		p.Shares,
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
