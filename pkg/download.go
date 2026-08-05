package pkg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Superredstone/spotiflac-cli/app"
	backend "github.com/Superredstone/spotiflac-cli/lib"
)

const (
	DEFAULT_DOWNLOAD_SERVICE       = "tidal"
	DEFAULT_DOWNLOAD_OUTPUT_FOLDER = "."
)

func Download(application *app.App, url string, output_folder string, service string) error {
	if output_folder == "" {
		output_folder = DEFAULT_DOWNLOAD_OUTPUT_FOLDER
	}

	if service == "" {
		service = DEFAULT_DOWNLOAD_SERVICE
	}

	if service == "amazon" || service == "qobuz" {
		isInstalled, err := application.CheckFFmpegInstalled()
		if err != nil {
			return err
		}

		if !isInstalled {
			return errors.New("FFmpeg is not installed.")
		}
	}

	// Tidal alone can bypass the community proxy via a persisted custom API URL; the others
	// always go through it, so only check the shared cooldown when that bypass isn't set.
	if service != "tidal" || strings.TrimSpace(backend.GetCustomTidalAPISetting()) == "" {
		if status, ok := application.GetCommunityBreakStatuses()[service]; ok && status.Available && status.IsBreak {
			return fmt.Errorf("%s community service is on a scheduled break, try again in ~%d minute(s)", service, status.RemainingMinutes)
		}
	}

	url_type := GetUrlType(url)

	switch url_type {
	case UrlTypeTrack:
		metadata, err := GetMetadata[MetadataSong](application, url)
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

		_, err = application.DownloadTrack(downloadRequest)
		return err
	case UrlTypePlaylist:
		metadata, err := GetMetadata[MetadataPlaylist](application, url)
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

			_, err = application.DownloadTrack(downloadRequest)
			if err != nil {
				fmt.Println("Unable to download " + track.Name + " - " + track.Artists)
			}
		}

		return nil
	case UrlTypeAlbum:
		metadata, err := GetMetadata[MetadataAlbum](application, url)
		if err != nil {
			return err
		}

		trackListSize := strconv.Itoa(len(metadata.TrackList))
		for idx, track := range metadata.TrackList {
			fmt.Println("[" + strconv.Itoa(idx+1) + "/" + trackListSize + "] " + track.Name + " - " + track.Artists)

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

			_, err = application.DownloadTrack(downloadRequest)
			if err != nil {
				fmt.Println("Unable to download " + track.Name + " - " + track.Artists)
			}
		}

		return nil
	}

	return errors.New("Invalid URL.")
}
