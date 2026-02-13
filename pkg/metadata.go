package pkg

import (
	"encoding/json"

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
	TrackList []MetadataTrack      `json:"track_list"`
	Info      MetadataPlaylistInfo `json:"playlist_info"`
}

type MetadataPlaylistInfo struct {
	Owner MetadataPlaylistOwner `json:"owner"`
}

type MetadataPlaylistOwner struct {
	Name string `json:"name"`
}

func GetMetadata[T MetadataPlaylist | MetadataSong](application *app.App, url string) (T, error) {
	var result T

	metadataRequest := app.SpotifyMetadataRequest{
		URL:     url,
		Delay:   0,
		Timeout: 5,
	}

	metadata, err := application.GetSpotifyMetadata(metadataRequest)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(metadata), &result)
	if err != nil {
		return result, nil
	}

	return result, nil
}

