package main

import (
	"image"
	"image/png"
	"log"
	"os"
)

type Thumbnail struct {
	width               int
	height              int
	bitspercomponent    int
	bitsperpixel        int
	bytesperrow         int
	bitmapdata_location int64
	bitmapdata_length   int64
	data                []byte
}

func NewThumbnail(width, height, bitspercomponent, bitsperpixel, bytesperrow int, bitmapdata_location, bitmapdata_length int64, f *os.File) *Thumbnail {
	thumbnail := &Thumbnail{width, height, bitspercomponent, bitsperpixel, bytesperrow, bitmapdata_location, bitmapdata_length, nil}
	thumbnail.fetchData(f)

	return thumbnail
}

func (thumb *Thumbnail) fetchData(f *os.File) {
	buf := make([]byte, thumb.bitmapdata_length)
	_, err := f.ReadAt(buf, thumb.bitmapdata_location)
	if err != nil {
		log.Print(err)
		return
	}

	thumb.data = buf
}

func (thumb *Thumbnail) CreateImage(output string) {
	img := &image.RGBA{Pix: thumb.data, Stride: thumb.bytesperrow, Rect: image.Rect(0, 0, thumb.width, thumb.height)}

	out, _ := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0600)
	defer out.Close()

	png.Encode(out, img)
}
