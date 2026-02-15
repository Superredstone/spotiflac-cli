package lib

import (
	"encoding/json"
	"errors"
	"fmt"
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
	Owner  MetadataPlaylistOwner  `json:"owner"`
	Tracks MetadataPlaylistTracks `json:"tracks"`
	Cover  string                 `json:"cover"`
}

type MetadataPlaylistTracks struct {
	Total int `json:"total"`
}

type MetadataPlaylistOwner struct {
	Name   string `json:"name"`         // Playlist name, this makes no sense
	Owner  string `json:"display_name"` // Playlist owner
	Images string `json:"images"`
}

func GetMetadata[T MetadataPlaylist | MetadataSong](app *App, url string) (T, error) {
	var result T

	metadataRequest := app.SpotifyMetadataRequest{
		URL:     url,
		Delay:   0,
		Timeout: 5,
	}

	metadata, err := app.GetSpotifyMetadata(metadataRequest)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(metadata), &result)
	if err != nil {
		return result, nil
	}

	return result, nil
}

func (app *App) PrintMetadata(url string) error {
	switch GetUrlType(url) {
	case UrlTypeTrack:
		metadata, err := GetMetadata[MetadataSong](app, url)
		if err != nil {
			return err
		}

		unformatted := `Name: %s
Artist: %s
Album: %s
Release date: %s
Images: %s`
		msg := fmt.Sprintf(unformatted,
			metadata.Track.Name,
			metadata.Track.Artists,
			metadata.Track.AlbumName,
			metadata.Track.ReleaseDate,
			metadata.Track.Images)
		fmt.Println(msg)

		return nil
	case UrlTypePlaylist:
		metadata, err := GetMetadata[MetadataPlaylist](app, url)
		if err != nil {
			return err
		}

		unformatted := `Name: %s
Owner: %s
Tracks: %d
Cover: %s`
		msg := fmt.Sprintf(unformatted,
			metadata.Info.Owner.Name,
			metadata.Info.Owner.Owner,
			metadata.Info.Tracks.Total,
			metadata.Info.Cover)
		fmt.Println(msg)

		return nil
	}

	return errors.New("Invalid URL.")
}
