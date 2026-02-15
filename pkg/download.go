package pkg

import (
	"errors"
	"fmt"
	"strconv"
"strings"
"flag"
	"spotiflac-cli/app"
	 "regexp"
	 "spotiflac-cli/lib"
	 "net/http"
	 "time"
	 "encoding/json"
)

const (
	DEFAULT_DOWNLOAD_SERVICE       = "qobuz"
	DEFAULT_DOWNLOAD_OUTPUT_FOLDER = "."
)
var isrcRegex = regexp.MustCompile(`^[A-Z]{2}[A-Z0-9]{3}\d{2}\d{5}$`)
func isValidISRC(isrc string) bool {
	return isrcRegex.MatchString(isrc)
}
func GetDeezerISRC(deezerURL string) (string, error) {

	var trackID string
	if strings.Contains(deezerURL, "/track/") {
		parts := strings.Split(deezerURL, "/track/")
		if len(parts) > 1 {
			trackID = strings.Split(parts[1], "?")[0]
			trackID = strings.TrimSpace(trackID)
		}
	}

	if trackID == "" {
		return "", fmt.Errorf("could not extract track ID from Deezer URL: %s", deezerURL)
	}

	apiURL := fmt.Sprintf("https://api.deezer.com/track/%s", trackID)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to call Deezer API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Deezer API returned status %d", resp.StatusCode)
	}

	var deezerTrack struct {
		ID    int64  `json:"id"`
		ISRC  string `json:"isrc"`
		Title string `json:"title"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&deezerTrack); err != nil {
		return "", fmt.Errorf("failed to decode Deezer API response: %w", err)
	}

	if deezerTrack.ISRC == "" {
		return "", fmt.Errorf("ISRC not found in Deezer API response for track %s", trackID)
	}

	fmt.Printf("Found ISRC from Deezer: %s (track: %s)\n", deezerTrack.ISRC, deezerTrack.Title)
	return deezerTrack.ISRC, nil
}
func Download(application *app.App, url string, output_folder string, service string) error {
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
			  ISRC:         track.ISRC,
		}
quality := flag.String("quality", "LOSSLESS", "Audio quality (LOSSLESS, HI_RES)")
	if output_folder == "" {
		output_folder = DEFAULT_DOWNLOAD_OUTPUT_FOLDER
	}

	if service == "" {
		service = DEFAULT_DOWNLOAD_SERVICE
	}

	if service == "amazon" {
		isInstalled, err := application.CheckFFmpegInstalled()
		if err != nil {
			return err
		}

		if !isInstalled {
			return errors.New("FFmpeg is not installed.")
		}
	}
if service == "qobuz" {
 // Map quality: 6 = 16-bit, 7 = 24-bit
        qCode := "6"
        if strings.ToUpper(*quality) == "HI_RES" {
            qCode = "7"
        }
deezerISRC := track.ISRC
			// If ISRC is not valid (looks like a Spotify ID), try to fetch from Deezer
			if len(deezerISRC) != 12 || !isValidISRC(deezerISRC) {
				fmt.Printf("ISRC is invalid (%s), fetching from Deezer...\n", deezerISRC)
				songlinkClient := lib.NewSongLinkClient()
				deezerURL, err := songlinkClient.GetDeezerURLFromSpotify(track.SpotifyID)
				if err == nil {
					deezerISRC, err = GetDeezerISRC(deezerURL)
				}
			}

			if deezerISRC == "" || !isValidISRC(deezerISRC) {
				fmt.Println("❌ Could not obtain a valid ISRC for Qobuz. Skipping.")
				
			}

			fmt.Printf("Using ISRC: %s\n", deezerISRC)

        fmt.Printf("🔍 Searching Qobuz by ISRC: %s\n", track.ISRC)
        qDownloader := lib.NewQobuzDownloader()
		var finalPath string
 
        finalPath, err = qDownloader.DownloadTrackWithISRC(
            deezerISRC,               // 1. ISRC
            DEFAULT_DOWNLOAD_OUTPUT_FOLDER,                   // 2. Output Dir
            DEFAULT_DOWNLOAD_OUTPUT_FOLDER,                        // 3. Quality
            qCode, // 4. Format
			 "{artist} - {title}", 
            false,                         // 5. includeTrackNumber
            1,                          // 6. position
            track.Name,               // 7. trackTitle
            track.Artists,            // 8. artists
            track.AlbumName,          // 9. albumTitle
            track.AlbumArtist,        // 10. albumArtist
            track.ReleaseDate,        // 11. releaseDate
            true,                         // 12. useAlbumTrackNumber
            "",                   // 13. coverPath
            true,                         // 14. embedMaxQualityCover
            track.TrackNumber,        // 15. spotifyTrackNumber
            track.DiscNumber,         // 16. spotifyDiscNumber
            track.TotalTracks,        // 17. spotifyTotalTracks
            track.TotalDiscs,         // 18. spotifyTotalDiscs
            track.Copyright,          // 19. spotifyCopyright
            track.Publisher,          // 20. spotifyPublisher
            "",        // 21. spotifyURL
			true,
			false,
        )
		report, err := lib.AnalyzeTrack(finalPath)
    if err == nil {
        fmt.Printf("Success: %s (%d-bit / %dHz)\n", finalPath, report.BitsPerSample, report.SampleRate)
    }
		_, err = application.DownloadTrack(downloadRequest)
		return nil
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
	}

	return errors.New("Invalid URL.")
}
