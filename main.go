package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Superredstone/spotiflac-cli/lib"
	"github.com/urfave/cli/v3"
)

func main() {
	outputFolder := ""
	service := ""

	app := lib.NewApp()
	err := app.Init()

	// Ignore this check for nix builds
	if err != nil && !strings.Contains(os.Args[0], "/nix/store/") {
		log.Fatal(err)
	}

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
						Usage:       "set output folder/file",
						DefaultText: outputFolder,
						Destination: &outputFolder,
					},
					&cli.StringFlag{
						Name:        "service",
						Aliases:     []string{"s"},
						Usage:       "set default service (only tidal is supported at the moment)",
						Destination: &service,
					},
					&cli.IntFlag{
						Name:        "interval",
						Aliases:     []string{"i"},
						Usage:       "interval between api requests in milliseconds",
						DefaultText: strconv.Itoa(app.ApiInterval),
						Destination: &app.ApiInterval,
					},
					&cli.BoolFlag{
						Name:        "no-fallback",
						Usage:       "do not fallback in case a source is not found",
						Destination: &app.NoFallback,
					},
					&cli.BoolFlag{
						Name:        "stop-on-fail",
						Usage:       "continue on download failure",
						Destination: &app.StopOnFail,
					},
					&cli.BoolFlag{
						Name:        "override",
						Usage:       "override already downloaded songs",
						Destination: &app.OverrideDownload,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					song_url := cmd.Args().First()
					quality := "LOSSLESS"
					err := app.Download(song_url, outputFolder, service, quality)
					return err
				},
			},
			&cli.Command{
				Name:    "metadata",
				Aliases: []string{"m"},
				Usage:   "view song metadata",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					url := cmd.Args().First()
					return app.PrintMetadata(url)
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Usage:       "verbose output",
				Destination: &app.Verbose,
			},
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
