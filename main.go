package main

import (
	"context"
	"log"
	"os"

	"github.com/Superredstone/spotiflac-cli/app"
	"github.com/Superredstone/spotiflac-cli/lib"
	"github.com/Superredstone/spotiflac-cli/pkg"
	"github.com/urfave/cli/v3"
)

func main() {
	var song_url string
	application := app.NewApp()
	startup()

	cmd := &cli.Command{
		Name: "spotiflac-cli",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "download",
				Usage:       "Download a song/playlist",
				Destination: &song_url,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			pkg.Download(application, song_url)

			return nil
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
		shutdown()
	}

	shutdown()
}

func startup() {
	if err := lib.InitHistoryDB("SpotiFLAC"); err != nil {
		log.Fatal("Failed to init history DB: %v\n", err)
	}
}

func shutdown() {
	lib.CloseHistoryDB()
}
