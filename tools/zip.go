package main

import (
	"archive/zip"
	"fmt"
	"io/fs"
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
		fmt.Printf("Failed to read file %s with error: %s", fileName, err)
	}
	zipFile, err := os.Create(zipName)
	if err != nil {
		fmt.Printf("Failed to create dest zip file with error: %s", err)
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
		fmt.Printf("Failed to create header file when zipping with error: %s", err)
	}
	_, err = dstZip.Write(srcFile)
	if err != nil {
		fmt.Printf("Failed to write to zip file for compressing: %s", err)
	}
}
