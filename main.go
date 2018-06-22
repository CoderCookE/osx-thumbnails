package main

import (
	"fmt"
	"log"
)

func main() {
	println("starting")
	dal := NewDAL()
	defer dal.Shutdown()

	rows, err := dal.Db.Query("SELECT width, height, bitspercomponent, bitsperpixel, bytesperrow, bitmapdata_location, bitmapdata_length FROM thumbnails")
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
			dal.DataFile,
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
