package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"os"
)

type DAL struct {
	Db       *sql.DB
	DataFile *os.File
}

func NewDAL() *DAL {
	createTmpFile()
	db, err := sql.Open("sqlite3", tmpDBFile())
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	dataFile := fmt.Sprintf("%s/thumbnails.data", getInputDir())
	f, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}

	return &DAL{
		db,
		f,
	}
}

func (d *DAL) Shutdown() {
	d.Db.Close()
	d.DataFile.Close()
}

func createTmpFile() (err error) {
	dbLocation := fmt.Sprintf("%s/index.sqlite", getInputDir())
	data, err := ioutil.ReadFile(dbLocation)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(tmpDBFile(), data, 0644)
	if err != nil {
		return
	}

	return
}

func tmpDBFile() string {
	return fmt.Sprintf("%s/tmposxthumbnails.sqlite", getTmpDir())
}

func getTmpDir() (tmpDir string) {
	tmpDir = os.TempDir()
	if len(tmpDir) == 0 {
		tmpDir = os.Getenv("TMPDIR")
	}
	return
}

func getInputDir() string {
	return fmt.Sprintf("%s../C/com.apple.QuickLook.thumbnailcache", getTmpDir())
}
