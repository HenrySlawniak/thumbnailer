# Thumbnailer

_An open-source contact sheet tool_

## Install

1. Download the latest release from the [Releases](https://github.com/HenrySlawniak/thumbnailer/releases) page
2. Download ffmpeg and ffprobe from the [official page](https://www.ffmpeg.org/download.html)
3. Place the appropriate ffmpeg and ffprobe binaries for your platform either:
  - Next to the thumbnailer binary
  - On your PATH

## Running

To generate contact sheets run `thumbnailer $FILENAME`.

By default, contact sheets will be written next to the video file. This can be disabled via the `in-place` flag.

See `thumbnailer -h` for a complete list of options.
