package lib

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	SONGLINK_API_BASE_URL    = "https://api.song.link/v1-alpha.1/links?url="
	RATE_LIMITED_RETURN_CODE = 429
)

type SongLinkResponse struct {
	EntityUniqueId  string          `json:"entityUniqueId"`
	UserCountry     string          `json:"userCountry"`
	PageUrl         string          `json:"pageUrl"`
	LinksByPlatform LinksByPlatform `json:"linksByPlatform"`
}

type LinksByPlatform struct {
	Deezer LinkByPlatform `json:"deezer"`
	Tidal  LinkByPlatform `json:"tidal"`
}

type LinkByPlatform struct {
	Country        string `json:"country"`
	Url            string `json:"url"`
	EntityUniqueId string `json:"entityUniqueId"`
}

func (app *App) ConvertSongUrl(url string) (SongLinkResponse, error) {
	var result SongLinkResponse

	app.log("Searching " + url)

	rawResponse, err := http.Get(SONGLINK_API_BASE_URL + url)
	if err != nil {
		return result, err
	}

	if rawResponse.StatusCode == RATE_LIMITED_RETURN_CODE {
		return result, errors.New("You have been rate limited by song.link, try again later.")
	}

	defer rawResponse.Body.Close()

	response, err := io.ReadAll(rawResponse.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(response, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
