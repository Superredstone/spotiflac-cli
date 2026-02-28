package lib

type App struct {
	SelectedTidalApiUrl string
	Verbose             bool
	SpotifyClient       *SpotifyClient
	ApiInterval         int // How many ms to wait between one call to apis and the other
	NoFallback          bool
	StopOnFail          bool
}

func NewApp() App {
	return App{
		Verbose:     false,
		ApiInterval: 800,
		NoFallback:  false,
	}
}

func (app *App) Init() error {
	app.log("Initializing Tidal")
	err := app.LoadTidalApis()
	if err != nil {
		return err
	}

	app.log("Initializing Spotify")
	if err := app.InitSpotifyClient(); err != nil {
		return err
	}

	return nil
}
