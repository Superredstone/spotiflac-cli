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
	var song_url, output_folder string

	application := app.NewApp()

	cmd := &cli.Command{
		Name: "spotiflac-cli",
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "download",
				Aliases:     []string{"d"},
				Usage:       "download a song/playlist",
				Destination: &song_url,
			},
			&cli.StringFlag{
				Name: "output",
				Aliases: []string{"o"},
				Usage: "set output folder",
				Destination: &output_folder,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			err := pkg.Download(application, song_url, output_folder)
			return err
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
