package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anpryl/image-storage/imagesvc/api"
	"github.com/anpryl/image-storage/imagesvc/images"
)

const (
	imagesRoute   = "/images"
	filenameParam = "filename"
)

func main() {
	var host = flag.String("host", "127.0.0.1", "Bind to host address (default: 127.0.0.1)")
	var port = flag.Int("port", 4000, "Use port for clients (default: 4000)")
	var secret = flag.String("secret", "", "Secret used to verify auth tokens")
	var folder = flag.String("folder", "images", "Path to image storage folder (default: images)")

	flag.Parse()

	if *secret == "" {
		log.Fatalln("secret is empty, please provide some secret with -secret flag")
	}

	if err := os.MkdirAll(*folder, 0700); err != nil {
		log.Fatalln("Failed to create folder for images: ", err)
	}

	st := images.NewStorage(*folder)

	mux := api.New(st, *secret)

	addr := addr(*host, *port)
	log.Println("Starting server at", addr)
	log.Fatalln(http.ListenAndServe(addr, mux))
}

func addr(host string, port int) string {
	return host + ":" + fmt.Sprint(port)
}
