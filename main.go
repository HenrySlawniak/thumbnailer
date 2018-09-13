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
	"strings"
)

var (
	numFrames        = flag.Int("frames", 12, "The number of frames to generate")
	writeInfo        = flag.Bool("write-info", true, "Write info JSON to file")
	frameWidth       = flag.Int("frame-width", 854, "The width to generate thumbnails at")
	outputDir        = flag.String("o", ".", "The directory to output to, ignored when in-place is true")
	outputInPlace    = flag.Bool("in-place", true, "Write images next to videos")
	framesPerRow     = flag.Int("frames-per-row", 3, "The number of frames per each row in the final contact sheet")
	writeAttribution = flag.Bool("write-attribution", true, "Writed \"Generated by thumbnailer.net\" to contact sheet")
	walkDirectories  = flag.Bool("walk-directories", true, "Walk directories provided as arguments")
	workers          = flag.Int("workers", runtime.NumCPU()-1, "Number of contact sheet workers to run, dafaults to number of CPUs -1")

	videoQueue = make(chan Video)

	buildTime string
	commit    string
)

func init() {
	flag.Parse()
	cLog := console.New(false)
	log.AddHandler(cLog, log.AllLevels...)
}

func main() {
	log.Infof("Starting Thumbnailer with %d workers", *workers)
	if buildTime != "" {
		log.Info("Built: " + buildTime)
	}
	if commit != "" {
		log.Info("Revision: " + commit)
	}
	log.Info("Go: " + runtime.Version())

	createDirectories()

	if len(flag.Args()) < 1 {
		log.Warn("Please provide a file path to generate a contact sheet from")
		log.Info("Use thumbnailer -h for a full list of options")
		os.Exit(1)
	}

	for w := 1; w <= *workers; w++ {
		go sheetWorker(w, videoQueue)
	}

	for _, a := range flag.Args() {
		if FileExists(a) {
			if !IsDir(a) {
				ProcessFile(a)
			} else {
				if *walkDirectories {
					WalkDir(a)
				}
			}
		}
	}

	close(videoQueue)
}

func sheetWorker(id int, videos chan Video) {
	for vid := range videos {
		log.Infof("Worker %d processing %s", id, vid.Filename)
		if *writeInfo {
			j, _ := json.MarshalIndent(vid, "", "  ")
			if *outputInPlace {
				ioutil.WriteFile(filepath.Join(filepath.Dir(vid.Location), vid.Filename+".json"), j, 0644)
			} else {
				ioutil.WriteFile(filepath.Join(*outputDir, vid.Filename+".json"), j, 0644)
			}
		}

		generateThumbnails(&vid, *numFrames)
		generateContactSheet(&vid, *numFrames)
	}
}

func createDirectories() {
	if !FileExists(*outputDir) {
		err := os.MkdirAll(*outputDir, 0755)
		if err != nil {
			log.Error("Cannot create output directory")
			log.Panic(err)
		}
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
	if filepath.Ext(path) == ".json" {
		return
	}

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
	log.Infof("Adding %s to queue", video.Filename)

	meta, err := getFFProbeMetadata(video.Location)
	if err != nil {
		log.Errorf("Error getting metadata for %s", video.Filename)
		log.Error(err)
		return
	}
	video.Meta = meta
	video.Duration = meta.DurationSeconds()

	for _, stream := range meta.Streams {
		if stream.CodecType == "video" && stream.AverageFrameRate != "0/0" {
			video.Width = stream.Width
			video.Height = stream.Height
			video.Codec = stream.CodecName
			break
		}
	}

	if video.Width < 1 || video.Height < 1 {
		return
	}

	if strings.Contains(meta.Format.FormatName, "pipe") {
		return
	}

	videoQueue <- video
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
