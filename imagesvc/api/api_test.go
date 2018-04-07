package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anpryl/image-storage/imagesvc/imgerrors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/powerman/check"
	"github.com/powerman/must"
)

const secret = "secret_phrase"

func TestJwtFail(tt *testing.T) {
	t := check.T{tt}
	st := &TestImageStorage{}
	mux := New(st, secret)
	srv := httptest.NewServer(mux)
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/images/1")
	if t.Nil(err) {
		t.EQ(resp.StatusCode, http.StatusUnauthorized)
	}

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/images/1", nil)
	must.NoErr(err)
	// Token already expired
	token, err := signToken(secret, -1*time.Hour)
	must.NoErr(err)
	req.Header.Set("Authorization", token)
	resp, err = http.DefaultClient.Do(req)
	if t.Nil(err) {
		t.EQ(resp.StatusCode, http.StatusUnauthorized)
	}
}

func TestJwtSuccess(tt *testing.T) {
	t := check.T{tt}
	st := &TestImageStorage{}
	mux := New(st, secret)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/images/1", nil)
	must.NoErr(err)
	// Token should expire before check
	token, err := signToken(secret, 1*time.Hour)
	fmt.Println(token)
	must.NoErr(err)
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if t.Nil(err) {
		t.EQ(resp.StatusCode, http.StatusNotFound)
	}
}

func TestGetNotFound(tt *testing.T) {
	t := check.T{tt}
	st := &TestImageStorage{}
	mux := routes(st)
	srv := httptest.NewServer(mux)
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/images/1")
	if t.Nil(err) {
		t.EQ(resp.StatusCode, http.StatusNotFound)
	}
}

func TestGet(tt *testing.T) {
	t := check.T{tt}
	file := []byte{1}
	name := "1"
	st := &TestImageStorage{
		images: map[string][]byte{
			name: file,
		},
	}
	mux := routes(st)
	srv := httptest.NewServer(mux)
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/images/" + name)
	if t.Nil(err) {
		t.EQ(resp.StatusCode, http.StatusOK)
	}
	defer resp.Body.Close()
	img := must.ReadAll(resp.Body)
	t.Zero(bytes.Compare(img, file))
}

func TestCreate(tt *testing.T) {
	t := check.T{tt}
	file := []byte{1}
	name := "1"
	st := &TestImageStorage{
		images: map[string][]byte{},
	}
	mux := routes(st)
	srv := httptest.NewServer(mux)
	defer srv.Close()
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/images/"+name, bytes.NewReader(file))
	must.NoErr(err)
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	t.Zero(bytes.Compare(st.images[name], file))
}

func TestDelete(tt *testing.T) {
	t := check.T{tt}
	file := []byte{1}
	name := "1"
	st := &TestImageStorage{
		images: map[string][]byte{
			name: file,
		},
	}
	mux := routes(st)
	srv := httptest.NewServer(mux)
	defer srv.Close()
	req, err := http.NewRequest(http.MethodDelete, srv.URL+"/images/"+name, nil)
	must.NoErr(err)
	resp, err := http.DefaultClient.Do(req)
	if t.Nil(err) {
		t.EQ(resp.StatusCode, http.StatusOK)
	}
	defer resp.Body.Close()

	_, ok := st.images[name]
	t.False(ok)
}

func TestImages(tt *testing.T) {
	t := check.T{tt}
	st := &TestImageStorage{
		images: map[string][]byte{
			"1": nil,
			"2": nil,
			"3": nil,
			"4": nil,
		},
	}
	mux := routes(st)
	srv := httptest.NewServer(mux)
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/images")
	if t.Nil(err) {
		t.EQ(resp.StatusCode, http.StatusOK)
	}
	defer resp.Body.Close()
	var imgResp imagesResp
	must.NoErr(json.NewDecoder(resp.Body).Decode(&imgResp))
	for k := range st.images {
		t.Contains(imgResp.Images, k)
	}
}

func TestConvertErr(tt *testing.T) {
	t := check.T{tt}
	err := &imgerrors.ImageServiceError{
		StatusCode: http.StatusTeapot,
		Message:    "Cookies!",
	}
	rw := httptest.NewRecorder()
	errToResp(rw, err)
	t.EQ(rw.Result().StatusCode, err.StatusCode)
}

func TestConvertErrInternalError(tt *testing.T) {
	t := check.T{tt}
	err := errors.New("Hi")
	rw := httptest.NewRecorder()
	errToResp(rw, err)
	t.EQ(rw.Result().StatusCode, http.StatusInternalServerError)
}

type TestImageStorage struct {
	images map[string][]byte

	err error
}

func (t *TestImageStorage) Save(filename string, r io.Reader) error {
	t.images[filename] = must.ReadAll(r)
	return t.err
}

func (t *TestImageStorage) Delete(filename string) error {
	delete(t.images, filename)
	return t.err
}

func (t *TestImageStorage) Get(filename string) (io.ReadCloser, error) {
	bs, ok := t.images[filename]
	if !ok {
		return nil, imgerrors.FileNotFound
	}
	return ioutil.NopCloser(bytes.NewReader(bs)), nil
}

func (t *TestImageStorage) Images() ([]string, error) {
	var images []string
	for k := range t.images {
		images = append(images, k)
	}
	return images, t.err
}

func signToken(secret string, tokenDuration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenDuration).Unix(),
	})
	return token.SignedString([]byte(secret))
}
