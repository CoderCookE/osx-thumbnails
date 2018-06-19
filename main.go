package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"image"
	"image/png"
	"os"
)

type Thumnails struct {
	width               int
	height              int
	bitspercomponent    int
	bitsperpixel        int
	bytesperrow         int
	bitmapdata_location int
	bitmapdata_length   int
}

func main() {
	println("starting")
	//echo $TMPDIR
	input := fmt.Sprintf("%s../C/com.apple.QuickLook.thumbnailcache", os.Args[1])
	println(input)

	dbLocation := fmt.Sprintf("%s/index.sqlite", input)
	println(dbLocation)
	db, err := sql.Open("sqlite3", dbLocation)
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("select * from thumnails limit 1;")
	if err != nil {
		panic(err)
	}

	thumbnailData := &Thumnails{}
	for rows.Next() {
		err = rows.Scan(thumbnailData)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(thumbnailData.width))
	}
	defer rows.Close()

	dataFile := fmt.Sprintf("%s/thumbnails.data", input)

	println(dataFile)
	f, err := os.Open(dataFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	dataLocation := int64(6583848)
	//ret, err := f.Seek(dataLocation, 0)
	//if err != nil {
	//	panic(err)
	//}
	//println(ret)

	const dataLength = 4096
	buf := make([]byte, dataLength)
	fileContents, err := f.ReadAt(buf, dataLocation)
	if err != nil {
		panic(err)
	}

	println(fileContents)

	println(len(buf))

	println(string(buf[:dataLength]))

	bytesperrow := 128
	bitsperpixel := 32
	bitspercomponent := 8

	width := bytesperrow / (bitsperpixel / bitspercomponent)
	height := 32

	img := &image.RGBA{Pix: buf, Stride: bytesperrow, Rect: image.Rect(0, 0, width, height)}

	out, _ := os.OpenFile("out.png", os.O_WRONLY|os.O_CREATE, 0600)
	defer out.Close()
	png.Encode(out, img)
}
