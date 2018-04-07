package images

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	errEmptyFilename = errors.New("Empty filename")
	errFileExist     = errors.New("File exist")
)

type ImageStorage struct {
	folder string
}

func (s *ImageStorage) Save(filename string, r io.Reader) error {
	fullPath, err := fullPath(s.folder, filename)
	if err != nil {
		return err
	}
	if err := checkIsFileExist(fullPath); err != nil {
		return err
	}
	return writeFile(fullPath, r)
}

func (s *ImageStorage) Delete(filename string) error {
	return nil
}

func fullPath(folder, filename string) (string, error) {
	filename = filepath.Base(filename)
	if filename == "" || filename == "." {
		return "", errEmptyFilename
	}
	return filepath.Join(folder, filename), nil
}

func writeFile(fullPath string, r io.Reader) error {
	f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)

	if err != nil {
		if removeErr := os.Remove(fullPath); removeErr != nil {
			log.Println("Failed to write file", removeErr)
		}
		return err
	}
	return nil
}

func checkIsFileExist(fullPath string) error {
	if _, err := os.Stat(fullPath); err == nil {
		return errFileExist
	}
	return nil
}
