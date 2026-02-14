package pkg

import (
	"errors"
	"fmt"
	"strconv"

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

	url_type := GetUrlType(url)

	switch url_type {
	case UrlTypeTrack:
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
	case UrlTypePlaylist:
		metadata, err := GetMetadata[MetadataPlaylist](application, url)
		if err != nil {
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

			_, err = application.DownloadTrack(downloadRequest)
			if err != nil { 
				fmt.Println("Unable to download " + track.Name + " - " + track.Artists)
			}
		}

		return nil
	}

	return errors.New("Invalid URL.")
}
