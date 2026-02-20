package lib

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (app *App) LoadTidalApis() error {
	var found bool

	for _, url := range app.GetAvailableApis() {
		res, err := http.Get(url)
		if err != nil {
			continue
		}

		if res.StatusCode == http.StatusOK {
			app.SelectedTidalApiUrl = url
			found = true
			break
		}
	}

	if !found {
		return errors.New("No available Tidal APIs found.")
	}

	return nil
}

func (app *App) GetAvailableApis() []string {
	// TODO: Make this load from a JSON file inside of $HOME/.config/spotiflac-cli/apis.json
	return []string{
		"https://triton.squid.wtf",
		"https://hifi-one.spotisaver.net",
		"https://hifi-two.spotisaver.net",
		"https://tidal.kinoplus.online",
		"https://tidal-api.binimum.org",
	}
}

type TidalAPIResponseV2 struct {
	Version string `json:"version"`
	Data    struct {
		TrackID           int64  `json:"trackId"`
		AssetPresentation string `json:"assetPresentation"`
		AudioMode         string `json:"audioMode"`
		AudioQuality      string `json:"audioQuality"`
		ManifestMimeType  string `json:"manifestMimeType"`
		ManifestHash      string `json:"manifestHash"`
		Manifest          string `json:"manifest"`
		BitDepth          int    `json:"bitDepth"`
		SampleRate        int    `json:"sampleRate"`
	} `json:"data"`
}

func (app *App) GetTidalDownloadUrl(tidalId string, quality string) (string, error) {
	url := fmt.Sprintf("%s/track/?id=%s&quality=%s", app.SelectedTidalApiUrl, tidalId, quality)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", app.UserAgent)

	rawResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer rawResponse.Body.Close()

	body, err := io.ReadAll(rawResponse.Body)
	if err != nil {
		return "", err
	}

	var response TidalAPIResponseV2
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if response.Data.Manifest != "" {
		manifest, err := app.ParseTidalManifestFromBase64(response.Data.Manifest)
		if err != nil {
			return "", err
		}

		if len(manifest.Urls) == 0 {
			return "", errors.New("No download URL found inside of Tidal APIs manifest.")
		}

		return manifest.Urls[0], nil
	}

	return "", errors.New("Unimplemented download from API v1.")
}

type TidalManifest struct {
	MimeType string   `json:"mimeType"`
	Codecs   string   `json:"codecs"`
	Urls     []string `json:"urls"`
}

func (app *App) ParseTidalManifestFromBase64(manifestBase64 string) (TidalManifest, error) {
	var result TidalManifest

	manifestDecoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(manifestBase64))
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(manifestDecoded, &result)
	if err != nil {
		return result, err
	} 

	return result, nil
}

func (app *App) GetTidalIdFromSonglink(songlink SongLinkResponse) (string, error) {
	return ParseTrackId(songlink.LinksByPlatform.Tidal.Url)
}
