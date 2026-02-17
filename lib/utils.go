package lib

import (
	"errors"
	"strings"
)

type UrlType int

const (
	UrlTypeTrack UrlType = iota
	UrlTypePlaylist
)

func ParseUrlType(url string) (UrlType, error) {
	if strings.Contains(url, "https://open.spotify.com/track") {
		return UrlTypeTrack, nil
	}

	if strings.Contains(url, "https://open.spotify.com/playlist") {
		return UrlTypePlaylist, nil
	}

	return UrlTypeTrack, errors.New("Invalid URL, not a playlist nor a track.")
}

func ParseTrackId(url string) (string, error) {
	tmp := strings.Split(url, "/")

	if len(tmp) == 0 {
		return "", errors.New("Invalid URL.")
	}

	tmp2 := strings.Split(tmp[len(tmp)-1], "?")
	if len(tmp2) == 0 {
		return tmp[len(tmp)-1], nil
	}

	return tmp2[0], nil
}
