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
)

var (
	numFrames  = flag.Int("frames", 12, "The number of frames to generate")
	writeInfo  = flag.Bool("write-info", true, "Write info JSON to file")
	fontFile   = flag.String("font-file", "IBMPlexMono-Regular.ttf", "The ttf font file to use")
	frameWidth = flag.Int("frame-width", 854, "The width to generate thumbnails at")
	fontBytes  []byte
)

func init() {
	flag.Parse()
	cLog := console.New(false)
	log.AddHandler(cLog, log.AllLevels...)
}

func main() {
	err := os.MkdirAll("tmp", 0755)
	if err != nil {
		log.Error("Cannot create tmp directory")
		log.Error(err)
		return
	}

	fontBytes, err = ioutil.ReadFile(*fontFile)
	if err != nil {
		log.Error("Error loading font")
		log.Error(err)
		return
	}

	for _, a := range os.Args[1:] {
		if FileExists(a) {
			f, err := os.Open(a)
			if err != nil {
				log.Error(err)
				continue
			}

			h := sha1.New()
			io.Copy(h, f)
			sum := h.Sum(nil)
			f.Close()

			video := Video{
				Filename: filepath.Base(a),
				Location: a,
				SHA1:     sum,
			}

			meta, err := getFFProbeMetadata(video.Location)
			if err != nil {
				log.Error(err)
				continue
			}
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
	}
}

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}
