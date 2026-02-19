package lib

import (
	"io"
	"net/http"
	"os"
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

func (app *App) Download(url string, outputFile string, serviceString string, quality string) error {
	var downloadUrl string
	var fileName string

	if serviceString == "" {
		serviceString = DEFAULT_DOWNLOAD_SERVICE
	}

	urlType, err := ParseUrlType(url)
	if err != nil {
		return err
	}

	switch urlType {
	case UrlTypeTrack:
		songlink, err := app.ConvertSongUrl(url)
		if err != nil {
			return err
		}

		tidalId, err := app.GetTidalIdFromSonglink(songlink)
		if err != nil {
			return err
		}

		downloadUrl, err = app.GetTidalDownloadUrl(tidalId, quality)
		if err != nil {
			return err
		}
	}

	metadata, err := app.GetTrackMetadata(url)
	if err != nil {
		return err
	}

	extension, err := GetFormatFromQuality(quality)
	if err != nil {
		return err
	}

	outputFile, err = BuildFileOutput(outputFile, fileName, extension, metadata)
	if err != nil {
		return err
	}

	fileExists, err := FileExists(outputFile)
	if err != nil {
		return err
	}

	if fileExists {
		app.log("File " + outputFile + " already exists")
		return nil
	}

	err = app.DownloadFromUrl(downloadUrl, outputFile)
	if err != nil {
		return err
	}

	err = app.EmbedMetadata(outputFile, metadata)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) DownloadFromUrl(url string, outputFilePath string) error {
	app.log("Downloading " + outputFilePath)

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(outputFile, res.Body)
	if err != nil {
		return err
	}

	app.log("Download completed")

	return nil
}
