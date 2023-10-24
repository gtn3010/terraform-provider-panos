package main

import (
	"archive/zip"
	"log"
	"os"
)

func main() {

	fileName := os.Args[1]
	zipName := os.Args[2]

	zipFile, err := os.Create(zipName)

	w := zip.NewWriter(zipFile)

	srcFile, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	dstZip, err := w.Create(zipName)
	if err != nil {
		log.Fatal(err)
	}
	_, err = dstZip.Write(srcFile)
	if err != nil {
		log.Fatal(err)
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
}
