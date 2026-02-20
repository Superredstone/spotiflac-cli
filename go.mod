module github.com/Superredstone/spotiflac-cli

go 1.24.4

replace github.com/Superredstone/spotiflac-cli/lib => ./lib

require (
	github.com/go-flac/flacpicture/v2 v2.0.2
	github.com/go-flac/flacvorbis/v2 v2.0.2
	github.com/go-flac/go-flac/v2 v2.0.4
	github.com/pquerna/otp v1.5.0
	github.com/urfave/cli/v3 v3.6.2
)

require github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
