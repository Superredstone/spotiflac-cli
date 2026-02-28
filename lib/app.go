package lib

type App struct {
	UserAgent           string // User agent used for scraping requests
	SelectedTidalApiUrl string
	Verbose             bool
	SpotifyClient       *SpotifyClient
	ApiInterval           int // How many ms to wait between one call to apis and the other
}

func NewApp() App {
	return App{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36",
		Verbose:   false,
		ApiInterval: 800,
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
