package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var host = flag.String("host", "127.0.0.1", "Bind to host address (default: 127.0.0.1)")
	var port = flag.Int("port", 80, "Use port for clients (default: 80)")

	mux := http.NewServeMux()
	mux.HandleFunc("/token", newToken())

	fmt.Println("vim-go")
	log.Fatalln(http.ListenAndServe(*host+":"+fmt.Sprint(*port), mux))
}

func newToken() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}
}
