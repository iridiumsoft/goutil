package util

import (
	"archive/zip"
	"os"
	"io"
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"path/filepath"
)

func StoreJsonFile(Path string, Data interface{}) bool {
	jsonData, _ := json.Marshal(Data)
	os.Create(Path)
	file, err := os.OpenFile(Path, os.O_RDWR, os.ModePerm)
	if err != nil {
		logrus.Error("error while storing file: ", err.Error())
		return false
	}
	file.Write(jsonData)
	return true
}

// This method will ZIP all provider files with delta compression

func ZipFiles(filename string, files []string) error {

	newfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newfile.Close()

	zipWriter := zip.NewWriter(newfile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// Get the file information
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, zipfile)
		if err != nil {
			return err
		}
	}
	return nil
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)
		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {

			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)

		} else {

			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()

			if err != nil {
				return filenames, err
			}

		}
	}
	return filenames, nil
}

func ZipAndRemoveFiles(files []string, Name string) {
	if len(files) > 0 {
		ZipFiles("./tmp/"+Name+".zip", files)
		// Remove Files
		for _, file := range files {
			os.Remove(file)
		}
	}
}

func DownloadAndRemoveFile(Name string, Path string, context *gin.Context) {
	context.Header("Content-Description", "File Transfer")
	context.Header("Content-Transfer-Encoding", "binary")
	context.Header("Content-Disposition", "attachment; filename="+Name)
	context.Header("Content-Type", "application/octet-stream")
	context.File(Path)
	os.Remove(Path)
}

func ReadFile(name string, ctx *gin.Context) {

}
