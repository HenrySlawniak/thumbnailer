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
	"encoding/json"
	"os/exec"
	"strconv"
)

type ffprobeOutput struct {
	Streams []ffprobeStreamInfo
	Format  struct {
		Duration   string
		FormatName string `json:"format_name"`
		BitRate    string `json:"bit_rate"`
		Size       string
	}
}

type ffprobeStreamInfo struct {
	Index            int
	CodecType        string `json:"codec_type"`
	CodecName        string `json:"codec_name"`
	AverageFrameRate string `json:"avg_frame_rate"`
	Width            int
	Height           int
}

func (o ffprobeOutput) DurationSeconds() float64 {
	f, _ := strconv.ParseFloat(o.Format.Duration, 64)
	return f
}

func getFFProbeMetadata(path string) (*ffprobeOutput, error) {
	binary := GetFFProbeBinary()

	cmd := exec.Command(
		binary,
		"-v", "error",
		"-show_streams",
		"-show_format",
		// "-show_entries", "format=width,height,duration_ts,duration,index,codec_type,codec_name,format_name,avg_frame_rate,bit_rate",
		"-print_format", "json",
		path,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	dat := ffprobeOutput{}
	err = json.Unmarshal(out, &dat)

	return &dat, err
}
