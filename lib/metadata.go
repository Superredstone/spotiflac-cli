package lib

import (
	"errors"
)

type Metadata struct {
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

func (app *App) GetMetadata(url string) (Metadata, error) {
	return Metadata{}, nil
}

func (app *App) PrintMetadata(url string) error {
	return errors.New("Invalid URL.")
}
