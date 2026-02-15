package pkg

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	DEFAULT_DOWNLOAD_SERVICE       = "tidal"
	DEFAULT_DOWNLOAD_OUTPUT_FOLDER = "."
)

func Download(app *App, url string, output_folder string, service string) error {
	if output_folder == "" {
		output_folder = DEFAULT_DOWNLOAD_OUTPUT_FOLDER
	}

	if service == "" {
		service = DEFAULT_DOWNLOAD_SERVICE
	}

	if service == "amazon" || service == "qobuz" {
		isInstalled, err := app.CheckFFmpegInstalled()
		if err != nil {
			return err
		}

		if !isInstalled {
			return errors.New("FFmpeg is not installed.")
		}
	}

	url_type := GetUrlType(url)

	switch url_type {
	case UrlTypeTrack:
		metadata, err := GetMetadata[MetadataSong](app, url)
		if err != nil {
			return err
		}

		track := metadata.Track
		downloadRequest := app.DownloadRequest{
			Service:     service,
			TrackName:   track.Name,
			ArtistName:  track.Artists,
			AlbumName:   track.AlbumName,
			AlbumArtist: track.AlbumArtist,
			ReleaseDate: track.ReleaseDate,
			CoverURL:    track.Images,
			OutputDir:   output_folder,
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

			downloadRequest := app.DownloadRequest{
				Service:      service,
				TrackName:    track.Name,
				ArtistName:   track.Artists,
				AlbumName:    track.AlbumName,
				AlbumArtist:  track.AlbumArtist,
				ReleaseDate:  track.ReleaseDate,
				CoverURL:     track.Images,
				OutputDir:    output_folder,
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
