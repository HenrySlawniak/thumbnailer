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
	"fmt"
	"github.com/go-playground/log"
	"os"
	"os/exec"
	"runtime/debug"
)

func generateThumbnails(vid *Video, numFrames int) {
	binary := GetFFMpegBinary()

	for i := 0; i < numFrames; i++ {
		cmd := exec.Command(
			binary, "-n",
			"-i", vid.Location,
			"-vframes", "1",
			"-ss", fmt.Sprintf("%f", ((float64(vid.Duration))/float64(numFrames))*float64(i)),
			"-vf", fmt.Sprintf("scale=%d:-1:", *frameWidth),
			fmt.Sprintf("%s/%s-%d.png", os.TempDir(), vid.SHA1.Hex(), i),
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			debug.PrintStack()
			log.Error(err.Error())
		}
	}
}
