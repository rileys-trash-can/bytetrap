package main

import (
	bt "github.com/derzombiiie/bytetrap"

	"flag"
	"log"
	"net/http"
)

const DefaultListen = "0.0.0.0:8080"

func main() {
	var (
		flagListen  = flag.String("listen", DefaultListen, "listen address")
		flagCheckUA = flag.Bool("check-ua", false, "only spam on 'Bytespider' UA")
	)

	flag.Parse()

	var (
		Listen  = *flagListen
		CheckUA = *flagCheckUA
	)

	log.Printf("Listening on %s", Listen)

	var handler http.Handler
	if CheckUA {
		handler = bt.Middleware(http.HandlerFunc(OkHandler))
	} else {
		handler = http.HandlerFunc(bt.Handler)
	}

	err := http.ListenAndServe(Listen, handler)
	log.Fatalf("Failed to listen on %s: %s", Listen, err)
}

func OkHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("not bytespider\n"))
}
