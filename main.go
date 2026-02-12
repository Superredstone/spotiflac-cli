package main

import (
	"context"
	"log"
	"os"

	"github.com/Superredstone/spotiflac-cli/app"
	"github.com/Superredstone/spotiflac-cli/pkg"
	"github.com/urfave/cli/v3"
)

func main() {
	var song_url string
	application := app.NewApp()

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
			err := pkg.Download(application, song_url)
			return err
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
