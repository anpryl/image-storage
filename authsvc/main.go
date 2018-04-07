package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func main() {
	var host = flag.String("host", "127.0.0.1", "Bind to host address (default: 127.0.0.1)")
	var port = flag.Int("port", 80, "Use port for clients (default: 80)")
	var secret = flag.String("secret", "", "Secret used to sign auth tokens")
	var tokenDuration = flag.Duration(
		"token_duration",
		5*time.Minute,
		"Duration after token expires (default: 5m)",
	)

	flag.Parse()

	if *secret == "" {
		log.Fatalln("secret is empty, please provide some secret with -secret flag")
	}

	mux := http.NewServeMux()

	signToken := tokenSigner(*secret, *tokenDuration)
	mux.HandleFunc("/token", newTokenHandler(signToken))

	log.Fatalln(http.ListenAndServe(*host+":"+fmt.Sprint(*port), mux))
}

type signTokenFunc func() (string, error)

func tokenSigner(secret string, tokenDuration time.Duration) signTokenFunc {
	return func() (string, error) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenDuration).Unix(),
		})
		return token.SignedString([]byte(secret))
	}
}

type TokenResp struct {
	Token string `json:"token"`
}

func newTokenHandler(signToken signTokenFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
		}
		token, err := signToken()
		if err != nil {
			log.Println("Error on token signing: ", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = json.NewEncoder(rw).Encode(TokenResp{Token: token}); err != nil {
			log.Println("Error on encoding resp: ", err)
		}
	}
}
