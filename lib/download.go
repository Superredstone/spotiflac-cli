package lib

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

func (app *App) Download(url string, outputFolder string, serviceString string) error {
	if outputFolder == "" {
		outputFolder = DEFAULT_DOWNLOAD_OUTPUT_FOLDER
	}

	if serviceString == "" {
		serviceString = DEFAULT_DOWNLOAD_SERVICE
	}

	urlType, err := ParseUrlType(url)
	if err != nil {
		return err
	}

	switch urlType {
	case UrlTypeTrack:
		// metadata, err := app.GetTrackMetadata(url)
		// if err != nil {
			// return err
		// }

		// println(metadata.Data.TrackUnion.Id)
		songlink, err := app.ConvertSongUrl(url)
		if err != nil {
			return err 
		}

		tidalId, err := app.GetTidalIdFromSonglink(songlink)
		if err != nil {
			return err
		}

		// err = app.DownloadFromTidal(tidalId)
		url, err = app.GetTidalDownloadUrl(tidalId, "LOSSLESS")
		if err != nil {
			return err
		}
		println(url)
	}

	return nil
}

func (app *App) DownloadTrack(dr DownloadRequest) (bool, error) {
	return false, nil
}
