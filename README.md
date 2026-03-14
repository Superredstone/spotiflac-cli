# spotiflac-cli

<p align="center">
  <img src="./assets/logo.png" width=300 height=300/>
</p>

> [!IMPORTANT]
> This project was originally based on [SpotiFLAC](github.com/afkarxyz/SpotiFLAC) which did not integrate a CLI and [the developer is not willing to add support for it](https://github.com/afkarxyz/SpotiFLAC/pull/381#issuecomment-3888433673).
>
> The code was a mess and not modular at all, which made adding new features painful. It was full of bad design choices. One of those bad design choices was not informing users how songs were actually downloaded. Tracks from Amazon were downloaded using what I later discovered to be closed-source APIs as M4A files, and then converted to FLAC with FFmpeg, which significantly reduced the quality. Qobuz as a source never existed in the first place, everything was downloaded either from Tidal or Amazon.
>
> For these reasons, I'm not going to support Amazon and Qobuz downloads until someone finds a better way to handle this.

Spotify downloader with playlist sync in mind.

## Usage 
```bash
spotiflac-cli download [URL] -v -o ~/Music/song.flac
```

## How to build 
1) Clone the repo
```bash
git clone https://github.com/Superredstone/spotiflac-cli && cd spotiflac-cli/
```
2) Build 
```bash
go build .
```

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
