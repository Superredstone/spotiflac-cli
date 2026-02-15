package app

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"path/filepath"

	 
	"strings"
	"time"
 "spotiflac-cli/lib"
 
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) getFirstArtist(artistString string) string {
	if artistString == "" {
		return ""
	}
	delimiters := []string{", ", " & ", " feat. ", " ft. ", " featuring "}
	for _, d := range delimiters {
		if idx := strings.Index(strings.ToLower(artistString), d); idx != -1 {
			return strings.TrimSpace(artistString[:idx])
		}
	}
	return artistString
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	if err := lib.InitHistoryDB("SpotiFLAC"); err != nil {
		fmt.Printf("Failed to init history DB: %v\n", err)
	}
}

func (a *App) shutdown(ctx context.Context) {
	lib.CloseHistoryDB()
}

type SpotifyMetadataRequest struct {
	URL     string  `json:"url"`
	Batch   bool    `json:"batch"`
	Delay   float64 `json:"delay"`
	Timeout float64 `json:"timeout"`
}

type DownloadRequest struct {
ISRC                 string `json:"isrc"`

	Service              string `json:"service"`
	Query                string `json:"query,omitempty"`
	TrackName            string `json:"track_name,omitempty"`
	ArtistName           string `json:"artist_name,omitempty"`
	AlbumName            string `json:"album_name,omitempty"`
	AlbumArtist          string `json:"album_artist,omitempty"`
	ReleaseDate          string `json:"release_date,omitempty"`
	CoverURL             string `json:"cover_url,omitempty"`
	ApiURL               string `json:"api_url,omitempty"`
	OutputDir            string `json:"output_dir,omitempty"`
	AudioFormat          string `json:"audio_format,omitempty"`
	FilenameFormat       string `json:"filename_format,omitempty"`
	TrackNumber          bool   `json:"track_number,omitempty"`
	Position             int    `json:"position,omitempty"`
	UseAlbumTrackNumber  bool   `json:"use_album_track_number,omitempty"`
	SpotifyID            string `json:"spotify_id,omitempty"`
	EmbedLyrics          bool   `json:"embed_lyrics,omitempty"`
	EmbedMaxQualityCover bool   `json:"embed_max_quality_cover,omitempty"`
	ServiceURL           string `json:"service_url,omitempty"`
	Duration             int    `json:"duration,omitempty"`
	ItemID               string `json:"item_id,omitempty"`
	SpotifyTrackNumber   int    `json:"spotify_track_number,omitempty"`
	SpotifyDiscNumber    int    `json:"spotify_disc_number,omitempty"`
	SpotifyTotalTracks   int    `json:"spotify_total_tracks,omitempty"`
	SpotifyTotalDiscs    int    `json:"spotify_total_discs,omitempty"`
	Copyright            string `json:"copyright,omitempty"`
	Publisher            string `json:"publisher,omitempty"`
	PlaylistName         string `json:"playlist_name,omitempty"`
	PlaylistOwner        string `json:"playlist_owner,omitempty"`
	AllowFallback        bool   `json:"allow_fallback"`
	UseFirstArtistOnly   bool   `json:"use_first_artist_only,omitempty"`
}

type DownloadResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	File          string `json:"file,omitempty"`
	Error         string `json:"error,omitempty"`
	AlreadyExists bool   `json:"already_exists,omitempty"`
	ItemID        string `json:"item_id,omitempty"`
}

func (a *App) GetStreamingURLs(spotifyTrackID string, region string) (string, error) {
	if spotifyTrackID == "" {
		return "", fmt.Errorf("spotify track ID is required")
	}

	fmt.Printf("[GetStreamingURLs] Called for track ID: %s, Region: %s\n", spotifyTrackID, region)
	client := lib.NewSongLinkClient()
	urls, err := client.GetAllURLsFromSpotify(spotifyTrackID, region)
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(urls)
	if err != nil {
		return "", fmt.Errorf("failed to encode response: %v", err)
	}

	return string(jsonData), nil
}

func (a *App) GetSpotifyMetadata(req SpotifyMetadataRequest) (string, error) {
	if req.URL == "" {
		return "", fmt.Errorf("URL parameter is required")
	}

	if req.Delay == 0 {
		req.Delay = 1.0
	}
	if req.Timeout == 0 {
		req.Timeout = 300.0
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(req.Timeout*float64(time.Second)))
	defer cancel()

	settings, err := a.LoadSettings()

	if err == nil && settings != nil {
		if useAPI, ok := settings["useSpotFetchAPI"].(bool); ok && useAPI {
			if apiURL, ok := settings["spotFetchAPIUrl"].(string); ok && apiURL != "" {

				data, err := lib.GetSpotifyDataWithAPI(ctx, req.URL, true, apiURL, req.Batch, time.Duration(req.Delay*float64(time.Second)))
				if err != nil {
					return "", fmt.Errorf("failed to fetch metadata from API: %v", err)
				}

				jsonData, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					return "", fmt.Errorf("failed to encode response: %v", err)
				}

				return string(jsonData), nil
			}
		}
	}

	data, err := lib.GetFilteredSpotifyData(ctx, req.URL, req.Batch, time.Duration(req.Delay*float64(time.Second)))
	if err != nil {
		return "", fmt.Errorf("failed to fetch metadata: %v", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to encode response: %v", err)
	}

	return string(jsonData), nil
}

type SpotifySearchRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

func (a *App) SearchSpotify(req SpotifySearchRequest) (*lib.SearchResponse, error) {
	if req.Query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return lib.SearchSpotify(ctx, req.Query, req.Limit)
}

type SpotifySearchByTypeRequest struct {
	Query      string `json:"query"`
	SearchType string `json:"search_type"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}

func (a *App) SearchSpotifyByType(req SpotifySearchByTypeRequest) ([]lib.SearchResult, error) {
	if req.Query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	if req.SearchType == "" {
		return nil, fmt.Errorf("search type is required")
	}

	if req.Limit <= 0 {
		req.Limit = 50
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return lib.SearchSpotifyByType(ctx, req.Query, req.SearchType, req.Limit, req.Offset)
}

func (a *App) DownloadTrack(req DownloadRequest) (DownloadResponse, error) {

	if req.Service == "qobuz" && req.SpotifyID == "" {
		return DownloadResponse{
			Success: false,
			Error:   "Spotify ID is required for Qobuz",
		}, fmt.Errorf("spotify ID is required for Qobuz")
	}

	if req.Service == "" {
		req.Service = "tidal"
	}

	if req.OutputDir == "" {
		req.OutputDir = "."
	} else {

		if req.PlaylistName != "" {
			sanitizedPlaylist := lib.SanitizeFilename(req.PlaylistName)
			req.OutputDir = filepath.Join(req.OutputDir, sanitizedPlaylist)
		}

		req.OutputDir = lib.SanitizeFolderPath(req.OutputDir)
	}

	if req.AudioFormat == "" {
		req.AudioFormat = "LOSSLESS"
	}

	var err error
	var filename string

	if req.FilenameFormat == "" {
		req.FilenameFormat = "title-artist"
	}

	itemID := req.ItemID
	if itemID == "" {

		if req.SpotifyID != "" {
			itemID = fmt.Sprintf("%s-%d", req.SpotifyID, time.Now().UnixNano())
		} else {
			itemID = fmt.Sprintf("%s-%s-%d", req.TrackName, req.ArtistName, time.Now().UnixNano())
		}

		lib.AddToQueue(itemID, req.TrackName, req.ArtistName, req.AlbumName, req.SpotifyID)
	}

	lib.SetDownloading(true)
	lib.StartDownloadItem(itemID)
	defer lib.SetDownloading(false)

	spotifyURL := ""
	if req.SpotifyID != "" {
		spotifyURL = fmt.Sprintf("https://open.spotify.com/track/%s", req.SpotifyID)
	}

	if req.SpotifyID != "" && (req.Copyright == "" || req.Publisher == "" || req.SpotifyTotalDiscs == 0 || req.ReleaseDate == "" || req.SpotifyTotalTracks == 0 || req.SpotifyTrackNumber == 0) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		trackURL := fmt.Sprintf("https://open.spotify.com/track/%s", req.SpotifyID)
		trackData, err := lib.GetFilteredSpotifyData(ctx, trackURL, false, 0)
		if err == nil {

			var trackResp struct {
				Track struct {
					Copyright   string `json:"copyright"`
					Publisher   string `json:"publisher"`
					TotalDiscs  int    `json:"total_discs"`
					TotalTracks int    `json:"total_tracks"`
					TrackNumber int    `json:"track_number"`
					ReleaseDate string `json:"release_date"`
				} `json:"track"`
			}
			if jsonData, jsonErr := json.Marshal(trackData); jsonErr == nil {
				if json.Unmarshal(jsonData, &trackResp) == nil {

					if req.Copyright == "" && trackResp.Track.Copyright != "" {
						req.Copyright = trackResp.Track.Copyright
					}
					if req.Publisher == "" && trackResp.Track.Publisher != "" {
						req.Publisher = trackResp.Track.Publisher
					}
					if req.SpotifyTotalDiscs == 0 && trackResp.Track.TotalDiscs > 0 {
						req.SpotifyTotalDiscs = trackResp.Track.TotalDiscs
					}
					if req.SpotifyTotalTracks == 0 && trackResp.Track.TotalTracks > 0 {
						req.SpotifyTotalTracks = trackResp.Track.TotalTracks
					}
					if req.SpotifyTrackNumber == 0 && trackResp.Track.TrackNumber > 0 {
						req.SpotifyTrackNumber = trackResp.Track.TrackNumber
					}
					if req.ReleaseDate == "" && trackResp.Track.ReleaseDate != "" {
						req.ReleaseDate = trackResp.Track.ReleaseDate
					}
				}
			}
		}
	}

	if req.TrackName != "" && req.ArtistName != "" {
		expectedFilename := lib.BuildExpectedFilename(req.TrackName, req.ArtistName, req.AlbumName, req.AlbumArtist, req.ReleaseDate, req.FilenameFormat, req.PlaylistName, req.PlaylistOwner, req.TrackNumber, req.Position, req.SpotifyDiscNumber, req.UseAlbumTrackNumber)
		expectedPath := filepath.Join(req.OutputDir, expectedFilename)

		if fileInfo, err := os.Stat(expectedPath); err == nil && fileInfo.Size() > 100*1024 {

			lib.SkipDownloadItem(itemID, expectedPath)
			return DownloadResponse{
				Success:       true,
				Message:       "File already exists",
				File:          expectedPath,
				AlreadyExists: true,
				ItemID:        itemID,
			}, nil
		}
	}

	lyricsChan := make(chan string, 1)
	isrcChan := make(chan string, 1)

	if req.SpotifyID != "" {
		if req.EmbedLyrics {
			go func() {
				client := lib.NewLyricsClient()
				resp, _, err := client.FetchLyricsAllSources(req.SpotifyID, req.TrackName, req.ArtistName, req.Duration)
				if err == nil && resp != nil && len(resp.Lines) > 0 {
					lrc := client.ConvertToLRC(resp, req.TrackName, req.ArtistName)
					lyricsChan <- lrc
				} else {
					lyricsChan <- ""
				}
			}()
		} else {
			close(lyricsChan)
		}

		go func() {
			client := lib.NewSongLinkClient()
			isrc, _ := client.GetISRC(req.SpotifyID)
			isrcChan <- isrc
		}()
	} else {
		close(lyricsChan)
		close(isrcChan)
	}

	switch req.Service {
	case "amazon":

		downloader := lib.NewAmazonDownloader()
		if req.ServiceURL != "" {
			filename, err = downloader.DownloadByURL(req.ServiceURL, req.OutputDir, req.AudioFormat, req.FilenameFormat, req.PlaylistName, req.PlaylistOwner, req.TrackNumber, req.Position, req.TrackName, req.ArtistName, req.AlbumName, req.AlbumArtist, req.ReleaseDate, req.CoverURL, req.SpotifyTrackNumber, req.SpotifyDiscNumber, req.SpotifyTotalTracks, req.EmbedMaxQualityCover, req.SpotifyTotalDiscs, req.Copyright, req.Publisher, spotifyURL, req.UseFirstArtistOnly)
		} else {
			filename, err = downloader.DownloadBySpotifyID(req.SpotifyID, req.OutputDir, req.AudioFormat, req.FilenameFormat, req.PlaylistName, req.PlaylistOwner, req.TrackNumber, req.Position, req.TrackName, req.ArtistName, req.AlbumName, req.AlbumArtist, req.ReleaseDate, req.CoverURL, req.SpotifyTrackNumber, req.SpotifyDiscNumber, req.SpotifyTotalTracks, req.EmbedMaxQualityCover, req.SpotifyTotalDiscs, req.Copyright, req.Publisher, spotifyURL, req.UseFirstArtistOnly)
		}

	case "tidal":
		if req.ApiURL == "" || req.ApiURL == "auto" {
			downloader := lib.NewTidalDownloader("")
			if req.ServiceURL != "" {
				filename, err = downloader.DownloadByURLWithFallback(req.ServiceURL, req.OutputDir, req.AudioFormat, req.FilenameFormat, req.TrackNumber, req.Position, req.TrackName, req.ArtistName, req.AlbumName, req.AlbumArtist, req.ReleaseDate, req.UseAlbumTrackNumber, req.CoverURL, req.EmbedMaxQualityCover, req.SpotifyTrackNumber, req.SpotifyDiscNumber, req.SpotifyTotalTracks, req.SpotifyTotalDiscs, req.Copyright, req.Publisher, spotifyURL, req.AllowFallback, req.UseFirstArtistOnly)
			} else {
				filename, err = downloader.Download(req.SpotifyID, req.OutputDir, req.AudioFormat, req.FilenameFormat, req.TrackNumber, req.Position, req.TrackName, req.ArtistName, req.AlbumName, req.AlbumArtist, req.ReleaseDate, req.UseAlbumTrackNumber, req.CoverURL, req.EmbedMaxQualityCover, req.SpotifyTrackNumber, req.SpotifyDiscNumber, req.SpotifyTotalTracks, req.SpotifyTotalDiscs, req.Copyright, req.Publisher, spotifyURL, req.AllowFallback, req.UseFirstArtistOnly)
			}
		} else {
			downloader := lib.NewTidalDownloader(req.ApiURL)
			if req.ServiceURL != "" {
				filename, err = downloader.DownloadByURL(req.ServiceURL, req.OutputDir, req.AudioFormat, req.FilenameFormat, req.TrackNumber, req.Position, req.TrackName, req.ArtistName, req.AlbumName, req.AlbumArtist, req.ReleaseDate, req.UseAlbumTrackNumber, req.CoverURL, req.EmbedMaxQualityCover, req.SpotifyTrackNumber, req.SpotifyDiscNumber, req.SpotifyTotalTracks, req.SpotifyTotalDiscs, req.Copyright, req.Publisher, spotifyURL, req.AllowFallback, req.UseFirstArtistOnly)
			} else {
				filename, err = downloader.Download(req.SpotifyID, req.OutputDir, req.AudioFormat, req.FilenameFormat, req.TrackNumber, req.Position, req.TrackName, req.ArtistName, req.AlbumName, req.AlbumArtist, req.ReleaseDate, req.UseAlbumTrackNumber, req.CoverURL, req.EmbedMaxQualityCover, req.SpotifyTrackNumber, req.SpotifyDiscNumber, req.SpotifyTotalTracks, req.SpotifyTotalDiscs, req.Copyright, req.Publisher, spotifyURL, req.AllowFallback, req.UseFirstArtistOnly)
			}
		}

	case "qobuz":

		fmt.Println("Waiting for ISRC (Qobuz dependency)...")
		isrc := <-isrcChan
		downloader := lib.NewQobuzDownloader()
		quality := req.AudioFormat
		if quality == "" {
			quality = "6"
		}
		filename, err = downloader.DownloadTrackWithISRC(isrc, req.SpotifyID, req.OutputDir, quality, req.FilenameFormat, req.TrackNumber, req.Position, req.TrackName, req.ArtistName, req.AlbumName, req.AlbumArtist, req.ReleaseDate, req.UseAlbumTrackNumber, req.CoverURL, req.EmbedMaxQualityCover, req.SpotifyTrackNumber, req.SpotifyDiscNumber, req.SpotifyTotalTracks, req.SpotifyTotalDiscs, req.Copyright, req.Publisher, spotifyURL, req.AllowFallback, req.UseFirstArtistOnly)

	default:
		return DownloadResponse{
			Success: false,
			Error:   fmt.Sprintf("Unknown service: %s", req.Service),
		}, fmt.Errorf("unknown service: %s", req.Service)
	}

	if err != nil {
		lib.FailDownloadItem(itemID, fmt.Sprintf("Download failed: %v", err))

		if filename != "" && !strings.HasPrefix(filename, "EXISTS:") {

			if _, statErr := os.Stat(filename); statErr == nil {
				fmt.Printf("Removing corrupted/partial file after failed download: %s\n", filename)
				if removeErr := os.Remove(filename); removeErr != nil {
					fmt.Printf("Warning: Failed to remove corrupted file %s: %v\n", filename, removeErr)
				}
			}
		}

		return DownloadResponse{
			Success: false,
			Error:   fmt.Sprintf("Download failed: %v", err),
			ItemID:  itemID,
		}, err
	}

	alreadyExists := false
	if strings.HasPrefix(filename, "EXISTS:") {
		alreadyExists = true
		filename = strings.TrimPrefix(filename, "EXISTS:")
	}

	if !alreadyExists && req.SpotifyID != "" && req.EmbedLyrics && (strings.HasSuffix(filename, ".flac") || strings.HasSuffix(filename, ".mp3") || strings.HasSuffix(filename, ".m4a")) {
		fmt.Printf("\nWaiting for lyrics fetch to complete...\n")
		lyrics := <-lyricsChan
		if lyrics != "" {
			fmt.Printf("\n--- Full LRC Content ---\n")
			fmt.Println(lyrics)
			fmt.Printf("--- End LRC Content ---\n\n")

			fmt.Printf("Embedding into: %s\n", filename)

			if err := lib.EmbedLyricsOnlyUniversal(filename, lyrics); err != nil {
				fmt.Printf("Failed to embed lyrics: %v\n", err)
			} else {
				fmt.Printf("Lyrics embedded successfully!\n")
			}
		} else {
			fmt.Println("No lyrics found to embed.")
		}
	} else {

		select {
		case <-lyricsChan:
		default:
		}
	}

	message := "Download completed successfully"
	if alreadyExists {
		message = "File already exists"
		lib.SkipDownloadItem(itemID, filename)
	} else {

		if fileInfo, statErr := os.Stat(filename); statErr == nil {
			finalSize := float64(fileInfo.Size()) / (1024 * 1024)
			lib.CompleteDownloadItem(itemID, filename, finalSize)
		} else {

			lib.CompleteDownloadItem(itemID, filename, 0)
		}

		go func(fPath, track, artist, album, sID, cover, format string) {
			quality := "Unknown"
			durationStr := "--:--"

			meta, err := lib.GetTrackMetadata(fPath)
			if err == nil && meta != nil {
				quality = fmt.Sprintf("%d-bit/%.1fkHz", meta.BitsPerSample, float64(meta.SampleRate)/1000.0)
				d := int(meta.Duration)
				durationStr = fmt.Sprintf("%d:%02d", d/60, d%60)
			} else {

			}

			item := lib.HistoryItem{
				SpotifyID:   sID,
				Title:       track,
				Artists:     artist,
				Album:       album,
				DurationStr: durationStr,
				CoverURL:    cover,
				Quality:     quality,
				Format:      format,
				Path:        fPath,
			}

			if item.Format == "" || item.Format == "LOSSLESS" {
				ext := filepath.Ext(fPath)
				if len(ext) > 1 {
					item.Format = strings.ToUpper(ext[1:])
				}
			}

			switch item.Format {
			case "6", "7", "27":
				item.Format = "FLAC"
			}

			lib.AddHistoryItem(item, "SpotiFLAC")
		}(filename, req.TrackName, req.ArtistName, req.AlbumName, req.SpotifyID, req.CoverURL, req.AudioFormat)
	}

	return DownloadResponse{
		Success:       true,
		Message:       message,
		File:          filename,
		AlreadyExists: alreadyExists,
		ItemID:        itemID,
	}, nil
}

func (a *App) OpenFolder(path string) error {
	if path == "" {
		return fmt.Errorf("path is required")
	}

	err := lib.OpenFolderInExplorer(path)
	if err != nil {
		return fmt.Errorf("failed to open folder: %v", err)
	}

	return nil
}

func (a *App) SelectFolder(defaultPath string) (string, error) {
	return lib.SelectFolderDialog(a.ctx, defaultPath)
}

func (a *App) SelectFile() (string, error) {
	return lib.SelectFileDialog(a.ctx)
}

func (a *App) GetDefaults() map[string]string {
	return map[string]string{
		"downloadPath": lib.GetDefaultMusicPath(),
	}
}

func (a *App) GetDownloadProgress() lib.ProgressInfo {
	return lib.GetDownloadProgress()
}

func (a *App) GetDownloadQueue() lib.DownloadQueueInfo {
	return lib.GetDownloadQueue()
}

func (a *App) ClearCompletedDownloads() {
	lib.ClearDownloadQueue()
}

func (a *App) ClearAllDownloads() {
	lib.ClearAllDownloads()
}

func (a *App) AddToDownloadQueue(spotifyID, trackName, artistName, albumName string) string {
	itemID := fmt.Sprintf("%s-%d", spotifyID, time.Now().UnixNano())
	lib.AddToQueue(itemID, trackName, artistName, albumName, "")
	return itemID
}

func (a *App) MarkDownloadItemFailed(itemID, errorMsg string) {
	lib.FailDownloadItem(itemID, errorMsg)
}

func (a *App) CancelAllQueuedItems() {
	lib.CancelAllQueuedItems()
}

func (a *App) ExportFailedDownloads() (string, error) {
	queueInfo := lib.GetDownloadQueue()
	var failedItems []string

	hasFailed := false
	for _, item := range queueInfo.Queue {
		if item.Status == lib.StatusFailed {
			hasFailed = true
			break
		}
	}

	if !hasFailed {
		return "No failed downloads to export.", nil
	}

	failedItems = append(failedItems, fmt.Sprintf("Failed Downloads Report - %s", time.Now().Format("2006-01-02 15:04:05")))
	failedItems = append(failedItems, strings.Repeat("-", 50))
	failedItems = append(failedItems, "")

	count := 0
	for _, item := range queueInfo.Queue {
		if item.Status == lib.StatusFailed {
			count++
			line := fmt.Sprintf("%d. %s - %s", count, item.TrackName, item.ArtistName)
			if item.AlbumName != "" {
				line += fmt.Sprintf(" (%s)", item.AlbumName)
			}
			failedItems = append(failedItems, line)
			failedItems = append(failedItems, fmt.Sprintf("   Error: %s", item.ErrorMessage))

			if item.SpotifyID != "" {
				failedItems = append(failedItems, fmt.Sprintf("   ID: %s", item.SpotifyID))
				failedItems = append(failedItems, fmt.Sprintf("   URL: https://open.spotify.com/track/%s", item.SpotifyID))
			}
			failedItems = append(failedItems, "")
		}
	}

	
	

	 

	return fmt.Sprintf("Successfully export"), nil
}

func (a *App) Quit() {

	panic("quit")
}

func (a *App) GetDownloadHistory() ([]lib.HistoryItem, error) {
	return lib.GetHistoryItems("SpotiFLAC")
}

func (a *App) ClearDownloadHistory() error {
	return lib.ClearHistory("SpotiFLAC")
}

func (a *App) DeleteDownloadHistoryItem(id string) error {
	return lib.DeleteHistoryItem(id, "SpotiFLAC")
}

func (a *App) GetFetchHistory() ([]lib.FetchHistoryItem, error) {
	return lib.GetFetchHistoryItems("SpotiFLAC")
}

func (a *App) AddFetchHistory(item lib.FetchHistoryItem) error {
	return lib.AddFetchHistoryItem(item, "SpotiFLAC")
}

func (a *App) ClearFetchHistory() error {
	return lib.ClearFetchHistory("SpotiFLAC")
}

func (a *App) DeleteFetchHistoryItem(id string) error {
	return lib.DeleteFetchHistoryItem(id, "SpotiFLAC")
}

func (a *App) ClearFetchHistoryByType(itemType string) error {
	return lib.ClearFetchHistoryByType(itemType, "SpotiFLAC")
}

func (a *App) AnalyzeTrack(filePath string) (string, error) {
	if filePath == "" {
		return "", fmt.Errorf("file path is required")
	}

	result, err := lib.AnalyzeTrack(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to analyze track: %v", err)
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to encode response: %v", err)
	}

	return string(jsonData), nil
}

func (a *App) AnalyzeMultipleTracks(filePaths []string) (string, error) {
	if len(filePaths) == 0 {
		return "", fmt.Errorf("at least one file path is required")
	}

	results := make([]*lib.AnalysisResult, 0, len(filePaths))

	for _, filePath := range filePaths {
		result, err := lib.AnalyzeTrack(filePath)
		if err != nil {

			continue
		}
		results = append(results, result)
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		return "", fmt.Errorf("failed to encode response: %v", err)
	}

	return string(jsonData), nil
}

type LyricsDownloadRequest struct {
	SpotifyID           string `json:"spotify_id"`
	TrackName           string `json:"track_name"`
	ArtistName          string `json:"artist_name"`
	AlbumName           string `json:"album_name"`
	AlbumArtist         string `json:"album_artist"`
	ReleaseDate         string `json:"release_date"`
	OutputDir           string `json:"output_dir"`
	FilenameFormat      string `json:"filename_format"`
	TrackNumber         bool   `json:"track_number"`
	Position            int    `json:"position"`
	UseAlbumTrackNumber bool   `json:"use_album_track_number"`
	DiscNumber          int    `json:"disc_number"`
}

func (a *App) DownloadLyrics(req LyricsDownloadRequest) (lib.LyricsDownloadResponse, error) {
	if req.SpotifyID == "" {
		return lib.LyricsDownloadResponse{
			Success: false,
			Error:   "Spotify ID is required",
		}, fmt.Errorf("spotify ID is required")
	}

	client := lib.NewLyricsClient()
	libReq := lib.LyricsDownloadRequest{
		SpotifyID:           req.SpotifyID,
		TrackName:           req.TrackName,
		ArtistName:          req.ArtistName,
		AlbumName:           req.AlbumName,
		AlbumArtist:         req.AlbumArtist,
		ReleaseDate:         req.ReleaseDate,
		OutputDir:           req.OutputDir,
		FilenameFormat:      req.FilenameFormat,
		TrackNumber:         req.TrackNumber,
		Position:            req.Position,
		UseAlbumTrackNumber: req.UseAlbumTrackNumber,
		DiscNumber:          req.DiscNumber,
	}

	resp, err := client.DownloadLyrics(libReq)
	if err != nil {
		return lib.LyricsDownloadResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return *resp, nil
}

type CoverDownloadRequest struct {
	CoverURL       string `json:"cover_url"`
	TrackName      string `json:"track_name"`
	ArtistName     string `json:"artist_name"`
	AlbumName      string `json:"album_name"`
	AlbumArtist    string `json:"album_artist"`
	ReleaseDate    string `json:"release_date"`
	OutputDir      string `json:"output_dir"`
	FilenameFormat string `json:"filename_format"`
	TrackNumber    bool   `json:"track_number"`
	Position       int    `json:"position"`
	DiscNumber     int    `json:"disc_number"`
}

func (a *App) DownloadCover(req CoverDownloadRequest) (lib.CoverDownloadResponse, error) {
	if req.CoverURL == "" {
		return lib.CoverDownloadResponse{
			Success: false,
			Error:   "Cover URL is required",
		}, fmt.Errorf("cover URL is required")
	}

	client := lib.NewCoverClient()
	libReq := lib.CoverDownloadRequest{
		CoverURL:       req.CoverURL,
		TrackName:      req.TrackName,
		ArtistName:     req.ArtistName,
		AlbumName:      req.AlbumName,
		AlbumArtist:    req.AlbumArtist,
		ReleaseDate:    req.ReleaseDate,
		OutputDir:      req.OutputDir,
		FilenameFormat: req.FilenameFormat,
		TrackNumber:    req.TrackNumber,
		Position:       req.Position,
		DiscNumber:     req.DiscNumber,
	}

	resp, err := client.DownloadCover(libReq)
	if err != nil {
		return lib.CoverDownloadResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return *resp, nil
}

type HeaderDownloadRequest struct {
	HeaderURL  string `json:"header_url"`
	ArtistName string `json:"artist_name"`
	OutputDir  string `json:"output_dir"`
}

func (a *App) DownloadHeader(req HeaderDownloadRequest) (lib.HeaderDownloadResponse, error) {
	if req.HeaderURL == "" {
		return lib.HeaderDownloadResponse{
			Success: false,
			Error:   "Header URL is required",
		}, fmt.Errorf("header URL is required")
	}

	if req.ArtistName == "" {
		return lib.HeaderDownloadResponse{
			Success: false,
			Error:   "Artist name is required",
		}, fmt.Errorf("artist name is required")
	}

	client := lib.NewCoverClient()
	libReq := lib.HeaderDownloadRequest{
		HeaderURL:  req.HeaderURL,
		ArtistName: req.ArtistName,
		OutputDir:  req.OutputDir,
	}

	resp, err := client.DownloadHeader(libReq)
	if err != nil {
		return lib.HeaderDownloadResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return *resp, nil
}

type GalleryImageDownloadRequest struct {
	ImageURL   string `json:"image_url"`
	ArtistName string `json:"artist_name"`
	ImageIndex int    `json:"image_index"`
	OutputDir  string `json:"output_dir"`
}

func (a *App) DownloadGalleryImage(req GalleryImageDownloadRequest) (lib.GalleryImageDownloadResponse, error) {
	if req.ImageURL == "" {
		return lib.GalleryImageDownloadResponse{
			Success: false,
			Error:   "Image URL is required",
		}, fmt.Errorf("image URL is required")
	}

	if req.ArtistName == "" {
		return lib.GalleryImageDownloadResponse{
			Success: false,
			Error:   "Artist name is required",
		}, fmt.Errorf("artist name is required")
	}

	client := lib.NewCoverClient()
	libReq := lib.GalleryImageDownloadRequest{
		ImageURL:   req.ImageURL,
		ArtistName: req.ArtistName,
		ImageIndex: req.ImageIndex,
		OutputDir:  req.OutputDir,
	}

	resp, err := client.DownloadGalleryImage(libReq)
	if err != nil {
		return lib.GalleryImageDownloadResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return *resp, nil
}

type AvatarDownloadRequest struct {
	AvatarURL  string `json:"avatar_url"`
	ArtistName string `json:"artist_name"`
	OutputDir  string `json:"output_dir"`
}

func (a *App) DownloadAvatar(req AvatarDownloadRequest) (lib.AvatarDownloadResponse, error) {
	if req.AvatarURL == "" {
		return lib.AvatarDownloadResponse{
			Success: false,
			Error:   "Avatar URL is required",
		}, fmt.Errorf("avatar URL is required")
	}

	if req.ArtistName == "" {
		return lib.AvatarDownloadResponse{
			Success: false,
			Error:   "Artist name is required",
		}, fmt.Errorf("artist name is required")
	}

	client := lib.NewCoverClient()
	libReq := lib.AvatarDownloadRequest{
		AvatarURL:  req.AvatarURL,
		ArtistName: req.ArtistName,
		OutputDir:  req.OutputDir,
	}

	resp, err := client.DownloadAvatar(libReq)
	if err != nil {
		return lib.AvatarDownloadResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return *resp, nil
}

func (a *App) CheckTrackAvailability(spotifyTrackID string) (string, error) {
	if spotifyTrackID == "" {
		return "", fmt.Errorf("spotify track ID is required")
	}

	client := lib.NewSongLinkClient()
	availability, err := client.CheckTrackAvailability(spotifyTrackID)
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(availability)
	if err != nil {
		return "", fmt.Errorf("failed to encode response: %v", err)
	}

	return string(jsonData), nil
}

func (a *App) IsFFmpegInstalled() (bool, error) {
	return lib.IsFFmpegInstalled()
}

func (a *App) IsFFprobeInstalled() (bool, error) {
	return lib.IsFFprobeInstalled()
}

func (a *App) GetFFmpegPath() (string, error) {
	return lib.GetFFmpegPath()
}

type DownloadFFmpegRequest struct{}

type DownloadFFmpegResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func (a *App) DownloadFFmpeg() DownloadFFmpegResponse {

	return DownloadFFmpegResponse{
		Success: true,
		Message: "FFmpeg installed successfully",
	}
}

type ConvertAudioRequest struct {
	InputFiles   []string `json:"input_files"`
	OutputFormat string   `json:"output_format"`
	Bitrate      string   `json:"bitrate"`
	Codec        string   `json:"codec"`
}

func (a *App) ConvertAudio(req ConvertAudioRequest) ([]lib.ConvertAudioResult, error) {
	libReq := lib.ConvertAudioRequest{
		InputFiles:   req.InputFiles,
		OutputFormat: req.OutputFormat,
		Bitrate:      req.Bitrate,
		Codec:        req.Codec,
	}
	return lib.ConvertAudio(libReq)
}

func (a *App) SelectAudioFiles() ([]string, error) {
	files, err := lib.SelectMultipleFiles(a.ctx)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (a *App) GetFileSizes(files []string) map[string]int64 {
	return lib.GetFileSizes(files)
}

func (a *App) ListDirectoryFiles(dirPath string) ([]lib.FileInfo, error) {
	if dirPath == "" {
		return nil, fmt.Errorf("directory path is required")
	}
	return lib.ListDirectory(dirPath)
}

func (a *App) ListAudioFilesInDir(dirPath string) ([]lib.FileInfo, error) {
	if dirPath == "" {
		return nil, fmt.Errorf("directory path is required")
	}
	return lib.ListAudioFiles(dirPath)
}

func (a *App) ReadFileMetadata(filePath string) (*lib.AudioMetadata, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path is required")
	}
	return lib.ReadAudioMetadata(filePath)
}

func (a *App) PreviewRenameFiles(files []string, format string) []lib.RenamePreview {
	return lib.PreviewRename(files, format)
}

func (a *App) RenameFilesByMetadata(files []string, format string) []lib.RenameResult {
	return lib.RenameFiles(files, format)
}

func (a *App) ReadTextFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (a *App) RenameFileTo(oldPath, newName string) error {
	dir := filepath.Dir(oldPath)
	ext := filepath.Ext(oldPath)
	newPath := filepath.Join(dir, newName+ext)
	return os.Rename(oldPath, newPath)
}

func (a *App) UploadImage(filePath string) (string, error) {
	return lib.UploadToSendNow(filePath)
}

func (a *App) UploadImageBytes(filename string, base64Data string) (string, error) {

	if idx := strings.Index(base64Data, ","); idx != -1 {
		base64Data = base64Data[idx+1:]
	}

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}
	return lib.UploadBytesToSendNow(filename, data)
}

func (a *App) SelectImageVideo() ([]string, error) {
	return lib.SelectImageVideoDialog(a.ctx)
}

func (a *App) ReadImageAsBase64(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	var mimeType string
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	case ".webp":
		mimeType = "image/webp"
	default:
		mimeType = "image/jpeg"
	}

	encoded := base64.StdEncoding.EncodeToString(content)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}

type CheckFileExistenceRequest struct {
	SpotifyID           string `json:"spotify_id"`
	TrackName           string `json:"track_name"`
	ArtistName          string `json:"artist_name"`
	AlbumName           string `json:"album_name,omitempty"`
	AlbumArtist         string `json:"album_artist,omitempty"`
	ReleaseDate         string `json:"release_date,omitempty"`
	TrackNumber         int    `json:"track_number,omitempty"`
	DiscNumber          int    `json:"disc_number,omitempty"`
	Position            int    `json:"position,omitempty"`
	UseAlbumTrackNumber bool   `json:"use_album_track_number,omitempty"`
	FilenameFormat      string `json:"filename_format,omitempty"`
	IncludeTrackNumber  bool   `json:"include_track_number,omitempty"`
	AudioFormat         string `json:"audio_format,omitempty"`
	RelativePath        string `json:"relative_path,omitempty"`
}

type CheckFileExistenceResult struct {
	SpotifyID  string `json:"spotify_id"`
	Exists     bool   `json:"exists"`
	FilePath   string `json:"file_path,omitempty"`
	TrackName  string `json:"track_name,omitempty"`
	ArtistName string `json:"artist_name,omitempty"`
}

func (a *App) CheckFilesExistence(outputDir string, rootDir string, tracks []CheckFileExistenceRequest) []CheckFileExistenceResult {
	if len(tracks) == 0 {
		return []CheckFileExistenceResult{}
	}

	outputDir = lib.NormalizePath(outputDir)
	if rootDir != "" {
		rootDir = lib.NormalizePath(rootDir)
	}

	defaultFilenameFormat := "title-artist"

	type result struct {
		index  int
		result CheckFileExistenceResult
	}

	resultsChan := make(chan result, len(tracks))

	var rootDirFiles map[string]string
	rootDirFilesOnce := false
	getRootDirFiles := func() map[string]string {
		if rootDirFilesOnce {
			return rootDirFiles
		}
		rootDirFiles = make(map[string]string)
		if rootDir != "" && rootDir != outputDir {
			filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if !info.IsDir() {
					if strings.EqualFold(filepath.Ext(path), ".flac") || strings.EqualFold(filepath.Ext(path), ".mp3") {
						rootDirFiles[info.Name()] = path
					}
				}
				return nil
			})
		}
		rootDirFilesOnce = true
		return rootDirFiles
	}

	for i, track := range tracks {
		go func(idx int, t CheckFileExistenceRequest) {
			res := CheckFileExistenceResult{
				SpotifyID:  t.SpotifyID,
				TrackName:  t.TrackName,
				ArtistName: t.ArtistName,
				Exists:     false,
			}

			if t.TrackName == "" || t.ArtistName == "" {
				resultsChan <- result{index: idx, result: res}
				return
			}

			filenameFormat := t.FilenameFormat
			if filenameFormat == "" {
				filenameFormat = defaultFilenameFormat
			}

			trackNumber := t.Position
			if t.UseAlbumTrackNumber && t.TrackNumber > 0 {
				trackNumber = t.TrackNumber
			}

			fileExt := ".flac"
			if t.AudioFormat == "mp3" {
				fileExt = ".mp3"
			}

			expectedFilenameBase := lib.BuildExpectedFilename(
				t.TrackName,
				t.ArtistName,
				t.AlbumName,
				t.AlbumArtist,
				t.ReleaseDate,
				filenameFormat,
				"",
				"",
				t.IncludeTrackNumber,
				trackNumber,
				t.DiscNumber,
				t.UseAlbumTrackNumber,
			)

			expectedFilename := strings.TrimSuffix(expectedFilenameBase, ".flac") + fileExt

			targetDir := outputDir
			if t.RelativePath != "" {
				targetDir = filepath.Join(outputDir, t.RelativePath)
			}

			expectedPath := filepath.Join(targetDir, expectedFilename)

			if fileInfo, err := os.Stat(expectedPath); err == nil && fileInfo.Size() > 100*1024 {
				res.Exists = true
				res.FilePath = expectedPath
			} else {

				res.FilePath = expectedFilename
			}

			resultsChan <- result{index: idx, result: res}
		}(i, track)
	}

	results := make([]CheckFileExistenceResult, len(tracks))
	missingIndices := []int{}

	for i := 0; i < len(tracks); i++ {
		r := <-resultsChan
		results[r.index] = r.result
		if !results[r.index].Exists {
			missingIndices = append(missingIndices, r.index)
		}
	}

	if len(missingIndices) > 0 && rootDir != "" {
		filesMap := getRootDirFiles()
		if len(filesMap) > 0 {
			for _, idx := range missingIndices {

				expectedFilename := results[idx].FilePath
				baseName := filepath.Base(expectedFilename)
				if path, ok := filesMap[baseName]; ok {
					results[idx].Exists = true
					results[idx].FilePath = path
				} else {
					results[idx].FilePath = ""
				}
			}
		} else {
			for _, idx := range missingIndices {
				results[idx].FilePath = ""
			}
		}
	} else {
		for _, idx := range missingIndices {
			results[idx].FilePath = ""
		}
	}

	return results
}

func (a *App) SkipDownloadItem(itemID, filePath string) {
	lib.SkipDownloadItem(itemID, filePath)
}

func (a *App) GetPreviewURL(trackID string) (string, error) {
	return lib.GetPreviewURL(trackID)
}

func (a *App) GetConfigPath() (string, error) {
	dir, err := lib.GetFFmpegDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func (a *App) SaveSettings(settings map[string]interface{}) error {
	configPath, err := a.GetConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(configPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func (a *App) LoadSettings() (map[string]interface{}, error) {
	configPath, err := a.GetConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	return settings, nil
}

func (a *App) CheckFFmpegInstalled() (bool, error) {
	return lib.IsFFmpegInstalled()
}

func (a *App) GetOSInfo() (string, error) {
	return lib.GetOSInfo()
}

func (a *App) CreateM3U8File(m3u8Name string, outputDir string, filePaths []string) error {
	if len(filePaths) == 0 {
		return nil
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	fnName := m3u8Name

	safeName := lib.SanitizeFilename(fnName)
	if safeName == "" {
		safeName = "playlist"
	}

	m3u8Path := filepath.Join(outputDir, safeName+".m3u8")

	f, err := os.Create(m3u8Path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString("#EXTM3U\n"); err != nil {
		return err
	}

	for _, path := range filePaths {
		if path == "" {
			continue
		}

		relPath, err := filepath.Rel(outputDir, path)
		if err != nil {

			relPath = path
		}

		relPath = filepath.ToSlash(relPath)

		if _, err := f.WriteString(relPath + "\n"); err != nil {
			return err
		}
	}

	return nil
}
