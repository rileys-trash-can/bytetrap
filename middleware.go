package bytetrap

import (
	"net/http"
	"strings"
)

const BytespiderUA = "Bytespider"

var bytespiderUA = strings.ToLower(BytespiderUA)

// checks if useragent contains BytespiderUA
func IsBytespider(useragent string) bool {
	// simple and effective
	return strings.Contains(strings.ToLower(useragent), bytespiderUA)
}

// see github.com/gorilla/mux#MiddlewareFunc
type MiddlewareFunc func(http.Handler) http.Handler

type middleware struct {
	handler http.Handler
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if IsBytespider(r.UserAgent()) {
		Write(w)

		return
	}

	m.handler.ServeHTTP(w, r)
}

// sends all `Bytespider`-UAs unlimited copypasta
func Middleware(next http.Handler) http.Handler {
	return &middleware{next}
}

// spams copypasta as response
func Handler(w http.ResponseWriter, r *http.Request) {
	Write(w)
}
