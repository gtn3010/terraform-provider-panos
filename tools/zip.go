package main

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"log"
	"os"
)

const (
	filePerm fs.FileMode = 0755
)

func main() {

	fileName := os.Args[1]
	zipName := os.Args[2]

	srcFile, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	zipFile, err := os.Create(zipName)
	if err != nil {
		log.Fatal(err)
	}
	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	fh := &zip.FileHeader{
		Name:   fileName,
		Method: 8, // Deflate compression
	}
	fh.SetMode(filePerm)
	// fmt.Println(fh.ExternalAttrs)
	// fmt.Println(filePerm)
	fmt.Println(uint32(0777))

	dstZip, err := w.CreateHeader(fh)
	if err != nil {
		log.Fatal(err)
	}
	_, err = dstZip.Write(srcFile)
	if err != nil {
		log.Fatal(err)
	}
}
