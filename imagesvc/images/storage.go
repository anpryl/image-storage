package images

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	errEmptyFilename = errors.New("Empty filename")
	errFileExist     = errors.New("File exist")
	errFileNotFound  = errors.New("File not found")
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
	fullPath, err := fullPath(s.folder, filename)
	if err != nil {
		return err
	}
	err = os.Remove(fullPath)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *ImageStorage) Get(filename string) (io.ReadCloser, error) {
	fullPath, err := fullPath(s.folder, filename)
	if err != nil {
		return nil, err
	}
	if err := checkIsFileNotExist(fullPath); err != nil {
		return nil, err
	}
	return os.Open(fullPath)
}

func (s *ImageStorage) Images() ([]string, error) {
	files, err := ioutil.ReadDir(s.folder)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, f := range files {
		filenames = append(filenames, f.Name())
	}
	return filenames, nil
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

func checkIsFileNotExist(fullPath string) error {
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return errFileNotFound
	}
	return nil
}

func checkIsFileExist(fullPath string) error {
	if _, err := os.Stat(fullPath); err == nil {
		return errFileExist
	}
	return nil
}
