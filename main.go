package main

import (
	"context"
	"log"
	"os"

	"spotiflac-cli/app"
	"spotiflac-cli/pkg"
	"github.com/urfave/cli/v3"
)

func main() {
	var output_folder, service string

	application := app.NewApp()

	cmd := &cli.Command{
		Name:                  "spotiflac-cli",
		EnableShellCompletion: true,
		DefaultCommand:        "help",
		Usage:                 "Spotify downloader with playlist sync in mind.",
		Commands: []*cli.Command{
			&cli.Command{
				Name:    "download",
				Aliases: []string{"d"},
				Usage:   "download a song/playlist",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "output",
						Aliases:     []string{"o"},
						Usage:       "set output folder",
						Destination: &output_folder,
					},
					&cli.StringFlag{
						Name:        "service",
						Aliases:     []string{"s"},
						Usage:       "set service to tidal/amazon/qobuz (FFmpeg is required for amazon and qobuz)",
						Destination: &service,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					song_url := cmd.Args().First()
					err := pkg.Download(application, song_url, output_folder, service)
					return err
				},
			},
			&cli.Command{
				Name:    "metadata",
				Aliases: []string{"m"},
				Usage:   "view song metadata",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					url := cmd.Args().First()
					return pkg.PrintMetadata(application, url)
				},
			},
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
