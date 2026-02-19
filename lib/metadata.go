package lib

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-flac/flacpicture/v2"
	"github.com/go-flac/flacvorbis/v2"
	"github.com/go-flac/go-flac/v2"
)

func (app *App) GetPlaylistMetadata(url string) (PlaylistMetadata, error) {
	app.log("Fetching playlist metadata")

	var result PlaylistMetadata
	playlistId, err := ParseTrackId(url)
	if err != nil {
		return result, err
	}

	payload := BuildSpotifyReqPayloadPlaylist(playlistId)

	rawMetadata, err := app.SpotifyClient.Query(payload)
	if err != nil {
		return result, err
	}

	byteMetadata, err := json.Marshal(rawMetadata)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(byteMetadata, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (app *App) GetTrackMetadata(url string) (TrackMetadata, error) {
	app.log("Fetching metadata for " + url)

	var result TrackMetadata

	trackId, err := ParseTrackId(url)
	if err != nil {
		return result, err
	}

	payload := BuildSpotifyReqPayloadTrack(trackId)

	rawMetadata, err := app.SpotifyClient.Query(payload)
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

func (app *App) EmbedMetadata(fileName string, metadata TrackMetadata) error {
	app.log("Embedding metadata")

	file, err := flac.ParseFile(fileName)
	if err != nil {
		return err
	}

	artists, err := GetArtists(metadata)
	if err != nil {
		return err
	}

	cmt := flacvorbis.New()
	cmt.Add(flacvorbis.FIELD_ALBUM, metadata.Data.TrackUnion.AlbumOfTrack.Name)
	cmt.Add(flacvorbis.FIELD_DATE, string(metadata.Data.TrackUnion.AlbumOfTrack.Date.IsoString.Year()))
	cmt.Add(flacvorbis.FIELD_ARTIST, artists)
	cmt.Add(flacvorbis.FIELD_TITLE, metadata.Data.TrackUnion.Name)
	cmtBlock := cmt.Marshal()
	file.Meta = append(file.Meta, &cmtBlock)

	cover, err := app.GetAlbumCover(metadata)
	if err != nil {
		return err
	}

	picture, err := flacpicture.NewFromImageData(
		flacpicture.PictureTypeFrontCover, "Front cover", cover, "image/jpeg")

	pictureMeta := picture.Marshal()
	file.Meta = append(file.Meta, &pictureMeta)
	file.Save(fileName)

	return nil
}

func (app *App) GetAlbumCover(metadata TrackMetadata) ([]byte, error) {
	app.log("Embedding cover")

	for _, source := range metadata.Data.TrackUnion.AlbumOfTrack.CoverArt.Sources {
		rawResponse, err := http.Get(source.Url)
		if err != nil {
			continue
		}
		defer rawResponse.Body.Close()

		response, err := io.ReadAll(rawResponse.Body)
		if err != nil {
			continue
		}

		return response, nil
	}

	return []byte{}, errors.New("Unable to download album cover for " + metadata.Data.TrackUnion.Name + ".")
}
