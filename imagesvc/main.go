package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

const (
	imagesRoute   = "/images"
	filenameParam = "filename"
)

func main() {
	var host = flag.String("host", "127.0.0.1", "Bind to host address (default: 127.0.0.1)")
	var port = flag.Int("port", 80, "Use port for clients (default: 80)")
	var secret = flag.String("secret", "", "Secret used to verify auth tokens")
	var folder = flag.String("folder", "images", "Path to image storage folder (default: images)")

	flag.Parse()

	if *secret == "" {
		log.Fatalln("secret is empty, please provide some secret with -secret flag")
	}

	if err := os.MkdirAll(*folder, 0700); err != nil {
		log.Fatalln("Failed to create folder for images: ", err)
	}

	r := httprouter.New()

	r.POST("/images/:"+filenameParam, saveImageHandle(*folder))

	// https://github.com/julienschmidt/httprouter#chaining-with-the-notfound-handler
	r.HandleMethodNotAllowed = false
	r.NotFound = http.FileServer(http.Dir(*folder)).ServeHTTP

	log.Fatalln(http.ListenAndServe(addr(*host, *port), r))
}

func addr(host string, port int) string {
	return host + ":" + fmt.Sprint(port)
}

func saveImageHandle(folder string) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		panic("not implemented")
	}
}
