package lib

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
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

func (app *App) Download(url string, outputFile string, service string, quality string) error {
	if service == "" {
		service = DEFAULT_DOWNLOAD_SERVICE
	}

	urlType, err := ParseUrlType(url)
	if err != nil {
		return err
	}

	switch urlType {
	case UrlTypeTrack:
		outputFileRune := []rune(outputFile)
		lastCharacter := string(outputFileRune[len(outputFileRune)-1:])
		downloadInFolder := lastCharacter == "/"

		if err := app.DownloadTrack(url, outputFile, service, quality, downloadInFolder); err != nil {
			return err
		}

		return nil
	case UrlTypePlaylist:
		if err := app.DownloadPlaylist(url, outputFile, service, quality); err != nil {
			return err
		}

		return nil
	}

	return errors.New("Invalid URL type.")
}

func (app *App) DownloadPlaylist(url string, outputFile string, service string, quality string) error {
	playlist, err := app.GetPlaylistMetadata(url)
	if err != nil {
		return err
	}

	var urls []string
	for _, item := range playlist.Data.Playlist.Content.Items {
		url, err := SpotifyUriToLink(item.Item.Data.Uri)
		if err != nil {
			return err
		}

		urls = append(urls, url)
	}

	trackListSize := len(urls)
	for idx, url := range urls {
		metadata, err := app.GetTrackMetadata(url)
		if err != nil {
			return err
		}

		artists, err := GetArtists(metadata)
		if err != nil {
			return err
		}

		fmt.Println("[" + strconv.Itoa(idx+1) + "/" + strconv.Itoa(trackListSize) + "] " + metadata.Data.TrackUnion.Name + " - " + artists)

		if err := app.DownloadTrack(url, outputFile+"/", service, quality, true); err != nil {
			return err
		}

		// Avoid getting rate limited
		time.Sleep(800 * time.Millisecond)
	}

	return nil
}

func (app *App) GetDownloadUrlOrFallback(askedService string, quality string, songlink SongLinkResponse) (string, error) {
	servicesToTry := []string{}

	switch askedService {
	default:
	case "tidal":
		servicesToTry = []string{"tidal", "amazon", "qoboz"}
		break
	case "amazon":
		servicesToTry = []string{"amazon", "tidal", "qoboz"}
		break
	case "qoboz":
		servicesToTry = []string{"qoboz", "tidal", "amazon"}
		break
	}

	var downloadUrl string
	var lastError error
	for idx, service := range servicesToTry {
		if idx > 0 {
			app.log("Falling back to " + service)
		}

		switch service {
		case "tidal":
			if songlink.LinksByPlatform.Tidal == nil {
				continue
			}

			tidalId, err := app.GetTidalIdFromSonglink(songlink)
			if err != nil {
				lastError = err
				continue
			}

			downloadUrl, err = app.GetTidalDownloadUrl(tidalId, quality)
			if err != nil {
				lastError = err
				continue
			}

			break
		}
	}

	if lastError != nil || downloadUrl == "" {
		return "", errors.New("Unable to download from any source.")
	}

	return downloadUrl, nil
}

func (app *App) DownloadTrack(url string, outputFile string, service string, quality string, downloadInFolder bool) error {
	songlink, err := app.ConvertSongUrl(url)
	if err != nil {
		return err
	}

	downloadUrl, err := app.GetDownloadUrlOrFallback(service, quality, songlink)
	if err != nil {
		return err
	}

	metadata, err := app.GetTrackMetadata(url)
	if err != nil {
		return err
	}

	extension, err := GetFormatFromQuality(quality)
	if err != nil {
		return err
	}

	if downloadInFolder {
		fileName, err := BuildFileName(metadata, extension)
		if err != nil {
			return err
		}
		outputFile = path.Join(outputFile, fileName)
	} else {
		outputFile, err = BuildFileOutput(outputFile, extension, metadata)
		if err != nil {
			return err
		}
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
