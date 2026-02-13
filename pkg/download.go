package pkg

import (
	"encoding/json"
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

type MetadataSong struct {
	Track MetadataTrack `json:"track"`
}

type MetadataTrack struct {
	SpotifyID    string `json:"spotify_id"`
	Artists      string `json:"artists"`
	Name         string `json:"name"`
	AlbumName    string `json:"album_name"`
	AlbumArtist  string `json:"album_artist"`
	DurationMS   int    `json:"duration_ms"`
	Images       string `json:"images"`
	ReleaseDate  string `json:"release_date"`
	TrackNumber  int    `json:"track_number"`
	TotalTracks  int    `json:"total_tracks"`
	DiscNumber   int    `json:"disc_number"`
	TotalDiscs   int    `json:"total_discs"`
	ExternalURLs string `json:"external_urls"`
	Copyright    string `json:"copyright"`
	Publisher    string `json:"publisher"`
	Plays        string `json:"plays"`
	IsExplicit   bool   `json:"is_explicit"`
}

type MetadataPlaylist struct {
	TrackList []MetadataTrack `json:"track_list"`
}

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

			application.DownloadTrack(downloadRequest)
		}

		return nil
	}

	return errors.New("Invalid Spotify URL.")
}

func GetMetadata[T MetadataPlaylist | MetadataSong](application *app.App, url string) (T, error) {
	var result T

	metadataRequest := app.SpotifyMetadataRequest{
		URL:     url,
		Delay:   0,
		Timeout: 5,
	}

	metadata, err := application.GetSpotifyMetadata(metadataRequest)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(metadata), &result)
	if err != nil {
		return result, nil
	}

	return result, nil
}
