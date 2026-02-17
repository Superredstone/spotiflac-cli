package lib

import (
	"encoding/json"
	"errors"
)

func (app *App) GetTrackMetadata(url string) (TrackMetadata, error) {
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
