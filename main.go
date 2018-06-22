package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"os"
)

const TMPLOCATION = "./tmp/dat1"

func main() {
	println("starting")
	input := fmt.Sprintf("%s../C/com.apple.QuickLook.thumbnailcache", os.Getenv("TMPDIR"))
	dbLocation := fmt.Sprintf("%s/index.sqlite", input)

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
	var bitmapdata_location int64
	var bitmapdata_length int64

	dataFile := fmt.Sprintf("%s/thumbnails.data", input)
	f, err := os.Open(dataFile)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	thumbnails := make([]*Thumbnail, 0)
	for rows.Next() {
		err = rows.Scan(&width, &height, &bitspercomponent, &bitsperpixel, &bytesperrow, &bitmapdata_location, &bitmapdata_length)
		if err != nil {
			log.Fatal(err)
		}

		thumbnail := NewThumbnail(
			width,
			height,
			bitspercomponent,
			bitsperpixel,
			bytesperrow,
			bitmapdata_location,
			bitmapdata_length,
			f,
		)

		thumbnails = append(thumbnails, thumbnail)
	}

	total := len(thumbnails)
	var successCount int
	for i, thumb := range thumbnails {
		if thumb.data != nil {
			filename := fmt.Sprintf("./output/out-%d.png", i)
			thumb.CreateImage(filename)
			successCount += 1
			log.Printf("Successfully created %d of %d thumbnails", successCount, total)
		} else {
			log.Print("Image creation failed, missing data")
		}
	}
}
