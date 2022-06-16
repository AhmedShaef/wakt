// Package upload provides a simple interface for uploading files to the server.
package upload

import (
	"io"
	"io/ioutil"
	"os"
)

//Image upload image to the server.
func Image(file io.Reader) (name string, err error) {
	tempFile, err := ioutil.TempFile("assets", "image-*.png")
	if err != nil {
		return "", err
	}
	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			panic(err)
		}
	}(tempFile)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

//Logo upload image to the server.
func Logo(file io.Reader) (name string, err error) {
	tempFile, err := ioutil.TempFile("assets", "logo-*.png")
	if err != nil {
		return "", err
	}
	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			panic(err)
		}
	}(tempFile)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}
