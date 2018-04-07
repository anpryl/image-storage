package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/anpryl/image-storage/imagesvc/imgerrors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

const (
	filenameParam = "filename"
	authHeader    = "Authorization"
)

type ImageStorage interface {
	Save(filename string, r io.Reader) error
	Delete(filename string) error
	Get(filename string) (io.ReadCloser, error)
	Images() ([]string, error)
}

func New(s ImageStorage, secret string) http.Handler {
	r := routes(s)
	return jwtMiddleware(secret, r.ServeHTTP)
}

func routes(s ImageStorage) http.Handler {
	r := httprouter.New()
	r.POST("/images/:"+filenameParam, saveImage(s))
	r.DELETE("/images/:"+filenameParam, deleteImage(s))
	r.GET("/images/:"+filenameParam, getImage(s))
	r.GET("/images", listImages(s))
	return r
}

func jwtMiddleware(secret string, next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get(authHeader)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(rw, r)
	}
}

func saveImage(s ImageStorage) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := s.Save(ps.ByName(filenameParam), r.Body)
		if err != nil {
			errToResp(rw, err)
			return
		}
		rw.WriteHeader(http.StatusCreated)
	}
}

func deleteImage(s ImageStorage) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := s.Delete(ps.ByName(filenameParam))
		if err != nil {
			errToResp(rw, err)
			return
		}
	}
}

func getImage(s ImageStorage) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		imageReader, err := s.Get(ps.ByName(filenameParam))
		if err != nil {
			errToResp(rw, err)
			return
		}
		defer func() {
			err := imageReader.Close()
			if err != nil {
				log.Println("Failed to close image reader", err)
			}
		}()
		_, err = io.Copy(rw, imageReader)
		if err != nil {
			log.Println("Failed to send response with image", err)
		}
	}
}

type imagesResp struct {
	Images []string `json:"images"`
}

func listImages(s ImageStorage) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		images, err := s.Images()
		if err != nil {
			errToResp(rw, err)
			return
		}
		resp := imagesResp{Images: images}
		if err := json.NewEncoder(rw).Encode(resp); err != nil {
			log.Println("Failed to render response", err)
		}
	}
}

func errToResp(rw http.ResponseWriter, err error) {
	if err, ok := err.(*imgerrors.ImageServiceError); ok {
		rw.WriteHeader(err.StatusCode)
		if err := json.NewEncoder(rw).Encode(err); err != nil {
			log.Println("Failed to render response", err)
		}
		return
	}
	rw.WriteHeader(http.StatusInternalServerError)
	log.Println("Internal error", err)
}
