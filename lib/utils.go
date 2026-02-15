package lib

import (
	"errors"
	"strings"
)

type UrlType int

const (
	UrlTypeTrack UrlType = iota
	UrlTypePlaylist
	UrlTypeInvalid
)

func ParseUrlType(url string) UrlType {
	if strings.Contains(url, "https://open.spotify.com/track") {
		return UrlTypeTrack
	}

	if strings.Contains(url, "https://open.spotify.com/playlist") {
		return UrlTypePlaylist
	}

	return UrlTypeInvalid
}

func ParseTrackId(url string) (string, error) {
	tmp := strings.Split(url, "/")

	if len(tmp) == 0 {
		return "", errors.New("Invalid URL.")
	}

	return tmp[len(tmp)-1], nil
}
