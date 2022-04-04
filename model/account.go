package model

import (
	"fmt"
	"regexp"
)

const (
	urlTikTok = `https://www.tiktok.com/`
	// parameter: 1. Username
	urlFormatUser = `https://www.tiktok.com/@%s`
)

func ToAccountURL(username string) string {
	return fmt.Sprintf(urlFormatUser, username)
}

type Account struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Bio         string `json:"bio"`
	Following   int    `json:"following"`
	Followers   int    `json:"followers"`
	Likes       int    `json:"likes"`
	AvatarURL   string `json:"avatar"`
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
