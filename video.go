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
	"fmt"
	"strconv"
	"strings"
)

type Video struct {
	Filename string
	Location string
	Duration float64
	SHA1     sha1sum
	Width    int
	Height   int
}

type sha1sum []byte

func (s *sha1sum) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Bytes []byte
		Hex   string
	}{
		Bytes: []byte(*s),
		Hex:   s.Hex(),
	})
}

func (s *sha1sum) Hex() string {
	return strings.TrimLeft(fmt.Sprintf("%x", s), "&")
}

type ffprobeOutput struct {
	Streams []ffprobeStreamInfo
	Format  struct {
		Duration string
	}
}

type ffprobeStreamInfo struct {
	Index     int
	CodecType string `json:"codec_type"`
	Width     int
	Height    int
}

func (o ffprobeOutput) DurationSeconds() float64 {
	f, _ := strconv.ParseFloat(o.Format.Duration, 64)
	return f
}
