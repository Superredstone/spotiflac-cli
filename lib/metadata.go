package lib

import (
	"encoding/json"
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
	urlType := ParseUrlType(url)

	switch urlType {
	case UrlTypeTrack:
		app.GetTrackMetadata(url)
	}

	return Metadata{}, errors.New("Invalid URL.")
}

func (app *App) GetTrackMetadata(url string) error {
	client := NewSpotifyClient()

	err := client.Initialize()
	if err != nil {
		return errors.New("Unable to fetch Spotify metadata.")
	}

	trackId, err := ParseTrackId(url)
	if err != nil {
		return err
	}

	payload := BuildSpotifyReqPayloadTrack(trackId)

	rawMetadata, err := client.Query(payload)
	if err != nil {
		return err
	}
	a, err := json.Marshal(rawMetadata)
	println(string(a))

	return nil
}

func (app *App) PrintMetadata(url string) error {
	return errors.New("Unimplemented.")
}
