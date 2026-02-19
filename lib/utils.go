package lib

import (
	"errors"
	"fmt"
	"os"
	"path"
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

func BuildFileName(metadata TrackMetadata, extension string) (string, error) {
	var result string

	artists, err := GetArtists(metadata)
	if err != nil {
		return result, err
	}

	result = fmt.Sprintf("%s - %s.%s", metadata.Data.TrackUnion.Name, artists, extension)

	return result, nil
}

func GetArtists(metadata TrackMetadata) (string, error) {
	var result string

	firstArtistLen := len(metadata.Data.TrackUnion.FirstArtist.Items)
	if firstArtistLen == 0 {
		return result, errors.New("What? This should never happen.")
	}
	result = metadata.Data.TrackUnion.FirstArtist.Items[firstArtistLen-1].Profile.Name

	for _, artist := range metadata.Data.TrackUnion.OtherArtists.Items {
		result += ", " + artist.Profile.Name
	}

	return result, nil
}

func BuildFileOutput(outputFile string, extension string, metadata TrackMetadata) (string, error) {
	var result string

	fileName, err := BuildFileName(metadata, extension)
	if err != nil {
		return result, err
	}

	if outputFile == "" {
		result = path.Join(DEFAULT_DOWNLOAD_OUTPUT_FOLDER, fileName)
	} else {
		result = outputFile
	}

	return result, nil
}

func (app *App) log(message string) {
	if app.Verbose {
		fmt.Println(message)
	}
}

func GetFormatFromQuality(quality string) (string, error) {
	switch quality {
	case "LOW":
		return "aac", nil
	case "HIGH":
		return "aac", nil
	case "LOSSLESS":
		return "flac", nil
	case "HI_RES_LOSSLESS":
		return "flac", nil
	default:
		return "", errors.New("Invalid quality.")
	}
}

func FileExists(file string) (bool, error) {
	_, err := os.Stat(file)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}
