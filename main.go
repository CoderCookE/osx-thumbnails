package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	println("starting")
	dal := NewDAL()
	defer dal.Shutdown()

	makeOutputDirectory()

	rows := dal.FindThumnails()
	defer rows.Close()

	var width int
	var height int
	var bitspercomponent int
	var bitsperpixel int
	var bytesperrow int
	var bitmapdata_location int64
	var bitmapdata_length int64

	thumbnails := make([]*Thumbnail, 0)
	for rows.Next() {
		err := rows.Scan(&width, &height, &bitspercomponent, &bitsperpixel, &bytesperrow, &bitmapdata_location, &bitmapdata_length)

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
			dal.DataFile,
		)

		thumbnails = append(thumbnails, thumbnail)
	}

	total := len(thumbnails)
	var successCount int

	for i, thumb := range thumbnails {
		filename := fmt.Sprintf("%s/out-%d.png", getOutputDir(), i)
		if err := thumb.CreateImage(filename); err != nil {
			log.Printf("Image creation failed: %s", err.Error())
		} else {
			successCount += 1
			log.Printf("Successfully created %d of %d thumbnails", successCount, total)
		}
	}
}

func makeOutputDirectory() {
	log.Printf("Outputting to %s", getOutputDir())

	if err := os.Mkdir(getOutputDir(), os.ModeDir|os.ModePerm); err != nil {
		log.Fatal(err.Error())
	}
}

func getOutputDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
	}

	return fmt.Sprintf("%s/osx-thumbnails-output", cwd)
}
