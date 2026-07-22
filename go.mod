module github.com/Superredstone/spotiflac-cli

go 1.25.0

replace github.com/Superredstone/spotiflac-cli/app => ./app

require (
	github.com/Eyevinn/mp4ff v0.54.0
	github.com/bogem/id3v2/v2 v2.1.4
	github.com/go-flac/flacpicture v0.3.0
	github.com/go-flac/flacvorbis v0.2.0
	github.com/go-flac/go-flac v1.0.0
	github.com/pquerna/otp v1.5.0
	github.com/ulikunitz/xz v0.5.15
	github.com/urfave/cli/v3 v3.6.2
	github.com/wailsapp/wails/v2 v2.11.0
	go.etcd.io/bbolt v1.4.3
	go.senan.xyz/taglib v0.13.0
	golang.org/x/image v0.23.0
	golang.org/x/text v0.22.0
)

require (
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/leaanthony/slicer v1.6.0 // indirect
	github.com/leaanthony/u v1.1.1 // indirect
	github.com/tetratelabs/wazero v1.11.1-0.20260428013916-2bbd517b7633 // indirect
	golang.org/x/sys v0.43.0 // indirect
)
