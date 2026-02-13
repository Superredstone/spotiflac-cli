_default: 
	@just --list

[group("windows")]
windows-amd64:
	GOOS="windows" GOARCH="amd64" go build -o build/spotiflac-cli-windows-amd64.exe

[group("windows")]
windows-arm64:
	GOOS="windows" GOARCH="arm64" go build -o build/spotiflac-cli-windows-arm64.exe
	
[group("darwin")]
darwin-amd64:
	GOOS="darwin" GOARCH="amd64" go build -o build/spotiflac-cli-macos-amd64

[group("darwin")]
darwin-arm64:
	GOOS="darwin" GOARCH="arm64" go build -o build/spotiflac-cli-macos-arm64

[group("linux")]
linux-amd64:
	GOOS="linux" GOARCH="amd64" go build -o build/spotiflac-cli-linux-amd64

[group("linux")]
linux-arm64:
	GOOS="linux" GOARCH="arm64" go build -o build/spotiflac-cli-linux-arm64

[group("windows")]
windows: windows-amd64 windows-arm64

[group("darwin")]
darwin: darwin-amd64 darwin-arm64

[group("linux")]
linux: linux-amd64 linux-arm64

# Build for all platforms
build: windows darwin linux

# Push binaries to GitHub releases
publish tag: 
	gh release upload $tag build/spotiflac-cli-*

clean:
	rm -rf build/
