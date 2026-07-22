package pkg

import (
	"strings"
)

type UrlType int

const (
	UrlTypeTrack UrlType = iota
	UrlTypePlaylist
	UrlTypeAlbum
	UrlTypeInvalid
)

func GetUrlType(url string) UrlType {
	if strings.Contains(url, "https://open.spotify.com/track") {
		return UrlTypeTrack
	}

	if strings.Contains(url, "https://open.spotify.com/playlist") {
		return UrlTypePlaylist
	}

	if strings.Contains(url, "https://open.spotify.com/album") {
		return UrlTypeAlbum
	}

	return UrlTypeInvalid
}
