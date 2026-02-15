package pkg

import (
	"strings"
)

type UrlType int

const (
	UrlTypeTrack UrlType = iota
	UrlTypePlaylist
	UrlTypeInvalid
)

func GetUrlType(url string) UrlType {
	if strings.Contains(url, "https://open.spotify.com/track") {
		return UrlTypeTrack
	}

	if strings.Contains(url, "https://open.spotify.com/playlist") {
		return UrlTypePlaylist
	}

	return UrlTypeInvalid
}
