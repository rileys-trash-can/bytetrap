package main

import (
	"github.com/google/uuid"
	bt "github.com/rileys-trash-can/bytetrap"

	"flag"
	"fmt"
	"log"
	"net/http"
)

const DefaultListen = "0.0.0.0:8080"

func main() {
	var (
		flagListen  = flag.String("listen", DefaultListen, "listen address")
		flagCheckUA = flag.Bool("check-ua", false, "only spam on 'Bytespider' UA")

		flagSchaffenburg = flag.Bool("schaffenburg", false, "schaffenburg wiki mode")
	)

	flag.Parse()

	var (
		Listen       = *flagListen
		CheckUA      = *flagCheckUA
		schaffenburg = *flagSchaffenburg
	)

	log.Printf("Listening on %s", Listen)

	var handler http.Handler
	handler = http.HandlerFunc(bt.Handler)

	if schaffenburg {
		handler = RandomizeHandler(handler, http.HandlerFunc(LinkPageHandler))
	}

	if CheckUA {
		handler = bt.Middleware(http.HandlerFunc(OkHandler))
	}

	err := http.ListenAndServe(Listen, handler)
	log.Fatalf("Failed to listen on %s: %s", Listen, err)
}

func RandomizeHandler(h1, h2 http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if bt.Rand.Intn(5) == 0 {
			h2.ServeHTTP(w, r)
		} else {
			h1.ServeHTTP(w, r)
		}
	})
}

func LinkPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)

	for i := 0; i < 100000; i++ {
		uuid, _ := uuid.NewRandom()

		fmt.Fprintf(w, "<p><a href=\"//wiki.schaffenburg.org/%s\">Ja halt inhalt Nr. %d</a></p>",
			uuid, i)
	}
}

func OkHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("not bytespider\n"))
}
