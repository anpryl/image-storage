package imgerrors

import (
	"net/http"
)

type ImageServiceError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (i *ImageServiceError) Error() string {
	return i.Message
}

var (
	EmptyFilename = &ImageServiceError{
		StatusCode: http.StatusBadRequest,
		Message:    "Empty filename",
	}
	FileExist = &ImageServiceError{
		StatusCode: http.StatusBadRequest,
		Message:    "File exist",
	}
	FileNotFound = &ImageServiceError{
		StatusCode: http.StatusNotFound,
		Message:    "File not found",
	}
)
