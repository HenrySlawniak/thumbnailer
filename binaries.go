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
	"path/filepath"
	"runtime"
)

// GetFFMpegBinary returns the location of the correct ffmpeg binary for the runtime OS
func GetFFMpegBinary() string {
	var location string
	if runtime.GOOS == "linux" {
		location = filepath.Join("bin", "ffmpeg")
	} else if runtime.GOOS == "windows" {
		location = filepath.Join("bin", "ffmpeg.exe")
	}

	if !FileExists(location) {
		return "ffmpeg"
	}

	return location
}

// GetFFProbeBinary returns the location of the correct ffprobe binary for the runtime OS
func GetFFProbeBinary() string {
	var location string
	if runtime.GOOS == "linux" {
		location = filepath.Join("bin", "ffprobe")
	} else if runtime.GOOS == "windows" {
		location = filepath.Join("bin", "ffprobe.exe")
	}

	if !FileExists(location) {
		return "ffprobe"
	}

	return location
}
