package lib

import (
	"encoding/json"
	"errors"

	id3v2 "github.com/bogem/id3v2/v2"
)

func (app *App) GetTrackMetadata(url string) (TrackMetadata, error) {
	app.log("Getting metadata for " + url)

	client := NewSpotifyClient()
	var result TrackMetadata

	err := client.Initialize()
	if err != nil {
		return result, errors.New("Unable to fetch Spotify metadata.")
	}

	trackId, err := ParseTrackId(url)
	if err != nil {
		return result, err
	}

	payload := BuildSpotifyReqPayloadTrack(trackId)

	rawMetadata, err := client.Query(payload)
	if err != nil {
		return result, err
	}

	byteMetadata, err := json.Marshal(rawMetadata)
	err = json.Unmarshal(byteMetadata, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (app *App) PrintMetadata(url string) error {
	return errors.New("Unimplemented.")
}

func (app *App) EmbedMetadata(file string, metadata TrackMetadata) error {
	tag, err := id3v2.Open(file, id3v2.Options{Parse: true})
	if err != nil {
		return err
	}

	artists, err := GetArtists(metadata)
	if err != nil {
		return err
	}

	tag.SetArtist(artists)
	tag.SetTitle(metadata.Data.TrackUnion.Name)
	tag.SetYear(string(metadata.Data.TrackUnion.AlbumOfTrack.Date.Year))
	tag.SetAlbum(metadata.Data.TrackUnion.AlbumOfTrack.Name)

	if err = tag.Save(); err != nil {
		return err
	}

	return nil
}
