package images

import (
	"os"

	"github.com/powerman/must"
)

func NewStorage(folder string) *ImageStorage {
	must.PanicIf(os.MkdirAll(folder, 0700))
	return &ImageStorage{folder}
}

func NewFileServer() {

}
