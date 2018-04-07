package images

import (
	"errors"
	"io"
)

var (
	errEmptyFilename = errors.New("Empty filename")
	errFileExist     = errors.New("File exist")
)

type ImageStorage struct {
	folder string
}

func (s *ImageStorage) Save(filename string, r io.Reader) error {
	return nil
}

func (s *ImageStorage) Delete(filename string) error {
	return nil
}
