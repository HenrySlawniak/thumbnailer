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
	"bufio"
	"fmt"
	"github.com/go-playground/log"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gomono"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

const (
	GutterSize   = 20
	FramesPerRow = 3

	FontSize    = 40
	FontSpacing = 0.9
	FontDPI     = 72
	HeaderSize  = 200
)

var (
	text = image.Black
	bg   = image.NewUniform(color.RGBA{0xE0, 0xEB, 0xF5, 0xff})
	font *truetype.Font
)

func init() {
	var err error
	font, err = freetype.ParseFont(gomono.TTF)
	if err != nil {
		log.Panic(err)
	}
}

func generateContactSheet(vid *Video, numFrames int) {
	var FrameWidth int
	var FrameHeight int
	frames := map[int]image.Image{}
	for i := 0; i < numFrames; i++ {
		frameLoc := fmt.Sprintf("tmp/%s-%d.png", vid.SHA1.Hex(), i)
		if !FileExists(frameLoc) {
			log.Errorf("Frame missing from disk `%s`\n", frameLoc)
			return
		}
		f, err := os.Open(frameLoc)
		if err != nil {
			log.Error(err)
			return
		}

		img, _, err := image.Decode(f)
		if err != nil {
			log.Error(err)
			return
		}
		f.Close()
		os.Remove(f.Name())

		if i == 0 {
			FrameWidth = img.Bounds().Dx()
			FrameHeight = img.Bounds().Dy()
		}

		frames[i] = img
	}

	log.Infof("Loaded %d frames for %s", len(frames), vid.Filename)

	rowCount := numFrames / FramesPerRow

	sheetWidth := (FramesPerRow * FrameWidth) + ((FramesPerRow + 1) * GutterSize)
	sheetHeight := (HeaderSize) + (rowCount * FrameHeight) + ((rowCount + 1) * GutterSize)
	log.Infof("Sheet Dimmensions: %dx%d\n", sheetWidth, sheetHeight)

	sheet := image.NewRGBA(image.Rect(0, 0, sheetWidth, sheetHeight))

	draw.Draw(sheet, sheet.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(FontDPI)
	c.SetFont(font)
	c.SetFontSize(FontSize)
	c.SetClip(sheet.Bounds())
	c.SetDst(sheet)
	c.SetSrc(text)

	pt := freetype.Pt(10, 10+int(c.PointToFixed(FontSize)>>6))
	for _, s := range vid.Filename {
		_, err := c.DrawString(string(s), pt)
		if err != nil {
			log.Error(err)
			return
		}
		pt.X += c.PointToFixed(FontSize * FontSpacing)
	}

	pt = freetype.Pt(10, 20+FontSize+int(c.PointToFixed((FontSize))>>6))
	for _, s := range "SHA1: " + vid.SHA1.Hex() {
		_, err := c.DrawString(string(s), pt)
		if err != nil {
			log.Error(err)
			return
		}
		pt.X += c.PointToFixed((FontSize * .7) * FontSpacing)
	}

	pt = freetype.Pt(10, 60+FontSize+int(c.PointToFixed((FontSize))>>6))
	for _, s := range fmt.Sprintf("Duration: %s, Dimmensions: %dx%d", stampToString(vid.Duration), vid.Width, vid.Height) {
		_, err := c.DrawString(string(s), pt)
		if err != nil {
			log.Error(err)
			return
		}
		pt.X += c.PointToFixed((FontSize * .7) * FontSpacing)
	}

	for i := 0; i < len(frames); i++ {
		frame := frames[i]
		row := i / FramesPerRow
		yOff := (row * FrameHeight) + GutterSize + HeaderSize + (GutterSize * row)
		col := i % FramesPerRow
		xOff := col*FrameWidth + GutterSize + (GutterSize * col)
		// frameTime := stampToString(((float64(vid.Duration)) / float64(numFrames)) * float64(i))
		rect := image.Rect(xOff, yOff, xOff+FrameWidth, yOff+FrameHeight)
		draw.Draw(sheet, rect, frame, frame.Bounds().Min, draw.Src)
	}

	outFile, err := os.Create(vid.Filename + ".png")
	if err != nil {
		log.Error(err)
		return
	}
	defer outFile.Close()

	b := bufio.NewWriter(outFile)
	err = png.Encode(b, sheet)
	if err != nil {
		log.Error(err)
		return
	}

	err = b.Flush()
	if err != nil {
		log.Error(err)
		return
	}

}

func stampToString(stamp float64) string {
	ts := int(stamp) % (24 * 3600)
	h := ts / 3600

	ts = ts % 3600
	m := ts / 60

	ts = ts % 60
	s := ts

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
