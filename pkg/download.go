package pkg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Superredstone/spotiflac-cli/app"
)

const (
	DEFAULT_DOWNLOAD_SERVICE       = "tidal"
	DEFAULT_DOWNLOAD_OUTPUT_FOLDER = "."
)

func Download(application *app.App, url string, output_folder string) error {
	if output_folder == "" {
		output_folder = DEFAULT_DOWNLOAD_OUTPUT_FOLDER
	}

	if strings.Contains(url, "https://open.spotify.com/track") {
		metadata, err := GetMetadata[MetadataSong](application, url)
		if err != nil {
			return err
		}

		track := metadata.Track
		downloadRequest := app.DownloadRequest{
			Service:     DEFAULT_DOWNLOAD_SERVICE,
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
	} else if strings.Contains(url, "https://open.spotify.com/playlist") {
		metadata, err := GetMetadata[MetadataPlaylist](application, url)
		if err != nil {
			fmt.Println("Unable to fetch metadata for song " + url)
			return err
		}

		trackListSize := strconv.Itoa(len(metadata.TrackList))
		for idx, track := range metadata.TrackList {
			fmt.Println("[" + strconv.Itoa(idx+1) + "/" + trackListSize + "] " + track.Name + " - " + track.Artists)

			downloadRequest := app.DownloadRequest{
				Service:      DEFAULT_DOWNLOAD_SERVICE,
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

			application.DownloadTrack(downloadRequest)
		}

		return nil
	}

	return errors.New("Invalid Spotify URL.")
}
