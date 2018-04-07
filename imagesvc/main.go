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
	var secret = flag.String("secret", "", "Secret used to verify auth tokens")

	if *secret == "" {
		log.Fatalln("secret is empty, please provide some secret with -secret flag")
	}

	flag.Parse()

	mux := http.NewServeMux()

	log.Fatalln(http.ListenAndServe(*host+":"+fmt.Sprint(*port), mux))
}
