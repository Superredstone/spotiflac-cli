package lib

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	DEFAULT_DOWNLOAD_SERVICE       = "tidal"
	DEFAULT_DOWNLOAD_OUTPUT_FOLDER = "."
)

type AvailableServices int

const (
	AvailableServicesTidal = iota
)

func ParseAvailableServices(service string) (AvailableServices, error) {
	switch service {
	case "tidal":
		return AvailableServicesTidal, nil
		break
	}
	return AvailableServicesTidal, errors.New("Invalid service.")
}

type DownloadRequest struct {
	Service AvailableServices
}

func (app *App) Download(url string, outputFolder string, serviceString string) error {
	if outputFolder == "" {
		outputFolder = DEFAULT_DOWNLOAD_OUTPUT_FOLDER
	}

	if serviceString == "" {
		serviceString = DEFAULT_DOWNLOAD_SERVICE
	}

	service, err := ParseAvailableServices(serviceString)
	if err != nil {
		return err
	}

	url_type := GetUrlType(url)

	switch url_type {
	case UrlTypeTrack:
		metadata, err := GetMetadata[MetadataSong](app, url)
		if err != nil {
			return err
		}

		track := metadata.Track
		downloadRequest := DownloadRequest{
			Service:     service,
			TrackName:   track.Name,
			ArtistName:  track.Artists,
			AlbumName:   track.AlbumName,
			AlbumArtist: track.AlbumArtist,
			ReleaseDate: track.ReleaseDate,
			CoverURL:    track.Images,
			OutputDir:   outputFolder,
			SpotifyID:   track.SpotifyID,
		}

		_, err = app.DownloadTrack(downloadRequest)
		return err
	case UrlTypePlaylist:
		metadata, err := GetMetadata[MetadataPlaylist](app, url)
		if err != nil {
			return err
		}

		trackListSize := strconv.Itoa(len(metadata.TrackList))
		for idx, track := range metadata.TrackList {
			fmt.Println("[" + strconv.Itoa(idx+1) + "/" + trackListSize + "] " + track.Name + " - " + track.Artists)

			downloadRequest := DownloadRequest{
				Service:      service,
				TrackName:    track.Name,
				ArtistName:   track.Artists,
				AlbumName:    track.AlbumName,
				AlbumArtist:  track.AlbumArtist,
				ReleaseDate:  track.ReleaseDate,
				CoverURL:     track.Images,
				OutputDir:    outputFolder,
				SpotifyID:    track.SpotifyID,
				PlaylistName: metadata.Info.Owner.Name,
			}

			_, err = app.DownloadTrack(downloadRequest)
			if err != nil {
				fmt.Println("Unable to download " + track.Name + " - " + track.Artists)
			}
		}

		return nil
	}

	return errors.New("Invalid URL.")
}

func (app *App) DownloadTrack(dr DownloadRequest) {

}
