// Copyright (c) 2018 Henry Slawniak <https://datacenterscumbags.com/>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"crypto/sha1"
	"encoding/json"
	"flag"
	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

var (
	numFrames  = flag.Int("frames", 12, "The number of frames to generate")
	writeInfo  = flag.Bool("write-info", true, "Write info JSON to file")
	frameWidth = flag.Int("frame-width", 854, "The width to generate thumbnails at")

	buildTime string
	commit    string
)

const binHelp = `This directory is used to store binaries for ffmpeg.

Place ffmpeg.exe, ffprobe.exe or the appropriate linux binaries here.

If no binaries are found, thumbnailer will use your PATH.
`

func init() {
	flag.Parse()
	cLog := console.New(false)
	log.AddHandler(cLog, log.AllLevels...)
}

func main() {
	log.Info("Starting Thumbnailer")
	if buildTime != "" {
		log.Info("Built: " + buildTime)
	}
	if commit != "" {
		log.Info("Revision: " + commit)
	}
	log.Info("Go: " + runtime.Version())

	createDirectories()

	for _, a := range os.Args[1:] {
		if FileExists(a) {
			if !IsDir(a) {
				ProcessFile(a)
			} else {
				WalkDir(a)
			}
		}
	}
}

func createDirectories() {
	err := os.MkdirAll("bin", 0755)
	if err != nil {
		log.Error("Cannot create bin directory")
		log.Panic(err)
	}

	err = ioutil.WriteFile(filepath.Join("bin", "readme.txt"), []byte(binHelp), 0644)
	if err != nil {
		log.Error("Error writing readme")
		log.Panic(err)
	}
}

func WalkDir(path string) {
	if !IsDir(path) {
		return
	}
	log.Infof("Walking %s", path)

	filepath.Walk(path, func(p string, stat os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if path == p {
			return nil
		}

		if IsDir(p) {
			WalkDir(p)
		} else {
			ProcessFile(p)
		}
		return nil
	})
}

func ProcessFile(path string) {
	log.Infof("processing %s", path)

	f, err := os.Open(path)
	if err != nil {
		log.Error(err)
		return
	}

	h := sha1.New()
	io.Copy(h, f)
	sum := h.Sum(nil)
	f.Close()

	video := Video{
		Filename: filepath.Base(path),
		Location: path,
		SHA1:     sum,
	}

	meta, err := getFFProbeMetadata(video.Location)
	if err != nil {
		log.Error(err)
		return
	}
	video.Meta = meta
	video.Duration = meta.DurationSeconds()

	for _, stream := range meta.Streams {
		if stream.CodecType == "video" {
			video.Width = stream.Width
			video.Height = stream.Height
			break
		}
	}

	if *writeInfo {
		j, _ := json.MarshalIndent(video, "", "  ")
		ioutil.WriteFile(video.Filename+".json", j, 0644)
	}

	generateThumbnails(&video, *numFrames)
	generateContactSheet(&video, *numFrames)
}

func IsDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		log.Panic(err)
	}
	return stat.IsDir()
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
