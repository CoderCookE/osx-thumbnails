package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

type Thumbnail struct {
	width               int
	height              int
	bitspercomponent    int
	bitsperpixel        int
	bytesperrow         int
	bitmapdata_location int
	bitmapdata_length   int
}

const TMPLOCATION = "./tmp/dat1"

func main() {
	println("starting")
	//echo $TMPDIR
	input := fmt.Sprintf("%s../C/com.apple.QuickLook.thumbnailcache", os.Getenv("TMPDIR"))
	println(input)

	dbLocation := fmt.Sprintf("%s/index.sqlite", input)

	println(dbLocation)

	data, err := ioutil.ReadFile(dbLocation)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(TMPLOCATION, data, 0644)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", TMPLOCATION)
	defer db.Close()

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT width, height, bitspercomponent, bitsperpixel, bytesperrow, bitmapdata_location, bitmapdata_length FROM thumbnails")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	var width int
	var height int
	var bitspercomponent int
	var bitsperpixel int
	var bytesperrow int
	var bitmapdata_location int
	var bitmapdata_length int

	thumbnails := make([]*Thumbnail, 0)
	for rows.Next() {
		err = rows.Scan(&width, &height, &bitspercomponent, &bitsperpixel, &bytesperrow, &bitmapdata_location, &bitmapdata_length)
		if err != nil {
			log.Fatal(err)
		}

		thumbnail := &Thumbnail{width, height, bitspercomponent, bitsperpixel, bytesperrow, bitmapdata_location, bitmapdata_length}
		thumbnails = append(thumbnails, thumbnail)
	}

	dataFile := fmt.Sprintf("%s/thumbnails.data", input)
	f, err := os.Open(dataFile)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	for i, thumb := range thumbnails {
		dataLocation := int64(thumb.bitmapdata_location)
		dataLength := int64(thumb.bitmapdata_length)

		buf := make([]byte, dataLength)
		_, err = f.ReadAt(buf, dataLocation)
		if err != nil {
			log.Fatal(err)
		}

		img := &image.RGBA{Pix: buf, Stride: thumb.bytesperrow, Rect: image.Rect(0, 0, thumb.width, thumb.height)}

		filename := fmt.Sprintf("./output/out-%d.png", i)
		out, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
		defer out.Close()

		png.Encode(out, img)
	}
}
