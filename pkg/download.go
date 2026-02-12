package pkg

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Superredstone/spotiflac-cli/app"
)

type MetadataSong struct {
	Track MetadataTrack `json:"track"`
}

type MetadataTrack struct {
	SpotifyID    string `json:"spotify_id"`
	Artists      string `json:"artists"`
	Name         string `json:"name"`
	AlbumName    string `json:"album_name"`
	AlbumArtist  string `json:"album_artist"`
	DurationMS   int    `json:"duration_ms"`
	Images       string `json:"images"`
	ReleaseDate  string `json:"release_date"`
	TrackNumber  int    `json:"track_number"`
	TotalTracks  int    `json:"total_tracks"`
	DiscNumber   int    `json:"disc_number"`
	TotalDiscs   int    `json:"total_discs"`
	ExternalURLs string `json:"external_urls"`
	Copyright    string `json:"copyright"`
	Publisher    string `json:"publisher"`
	Plays        string `json:"plays"`
	IsExplicit   bool   `json:"is_explicit"`
}

type MetadataPlaylist struct {
	TrackList []MetadataTrack `json:"track_list"`
}

func Download(application *app.App, url string) {
	metadata, err := GetMetadata[MetadataPlaylist](application, url)
	if err != nil {
		fmt.Println("Unable to fetch metadata for song " + url)
	}

	trackListSize := strconv.Itoa(len(metadata.TrackList))
	for idx, track := range metadata.TrackList {
		fmt.Println("[" + strconv.Itoa(idx+1) + "/" + trackListSize + "] " + track.Name + " - " + track.Artists)

		downloadRequest := app.DownloadRequest{
			Service:     "tidal",
			TrackName:   track.Name,
			ArtistName:  track.Artists,
			AlbumName:   track.AlbumName,
			AlbumArtist: track.AlbumArtist,
			ReleaseDate: track.ReleaseDate,
			OutputDir:   "downloads/",
			SpotifyID:   track.SpotifyID,
		}
		application.DownloadTrack(downloadRequest)
	}
}

func GetMetadata[T MetadataPlaylist | MetadataSong](application *app.App, url string) (T, error) {
	var result T

	metadata, err := GetGenericMetadata(application, url)
	if err != nil {
		return result, nil
	}

	err = json.Unmarshal([]byte(metadata), &result)
	if err != nil {
		return result, nil
	}

	return result, nil
}

func GetGenericMetadata(application *app.App, url string) (string, error) {
	metadataRequest := app.SpotifyMetadataRequest{
		URL:     url,
		Delay:   0,
		Timeout: 5,
	}

	metadata, err := application.GetSpotifyMetadata(metadataRequest)
	if err != nil {
		return metadata, err
	}

	return metadata, nil
}
