package lib

import (
	"errors"
)

const (
	DEFAULT_DOWNLOAD_SERVICE       = "tidal"
	DEFAULT_DOWNLOAD_OUTPUT_FOLDER = "."
)

type DownloadRequest struct {
	Service     AvailableServices
	Track       string
	Artist      string
	Album       string
	Artists     string
	ReleaseDate string
	Cover       string
	OutputDir   string
	SpotifyID   string
}

func (app *App) Download(url string, outputFolder string, serviceString string) error {
	if outputFolder == "" {
		outputFolder = DEFAULT_DOWNLOAD_OUTPUT_FOLDER
	}

	if serviceString == "" {
		serviceString = DEFAULT_DOWNLOAD_SERVICE
	}

	_, err := app.GetMetadata(url)
	if err != nil {
		return err
	}

	_ = ParseUrlType(url)

	return errors.New("Invalid URL.")
}

func (app *App) DownloadTrack(dr DownloadRequest) (bool, error) {
	return false, nil
}
