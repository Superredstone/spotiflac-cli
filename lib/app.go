package lib

type App struct {
	UserAgent           string // User agent used for scraping requests
	SelectedTidalApiUrl string
	Verbose             bool
	SpotifyClient       *SpotifyClient
}

func NewApp() App {
	return App{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36",
		Verbose:   false,
	}
}

func (app *App) Init() error {
	err := app.LoadTidalApis()
	if err != nil {
		return err
	}

	if err := app.InitSpotifyClient(); err != nil {
		return err
	}

	return nil
}
