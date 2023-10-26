package bytetrap

import (
	"net/http"
	"strings"
	"time"

	"fmt"
	"log"
)

// int in bytes
var statsCh chan int64

// stolen from https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func startStats() {
	ticker := time.NewTicker(time.Second * 10)

	statsCh = make(chan int64)
	var total, i, alltime int64
	
	starttime := time.Now()


	for {
		select {
		case i = <-statsCh:
			total += i

		case <-ticker.C:
			if total != 0 {
				dur := int64(time.Now().Sub(starttime).Seconds())
				alltime += total

				log.Printf("[bytetrap/stats] sent %s (%d) bytes in the last 10s (%s over %ds -> %sbit/s)",
					ByteCountSI(total), total,
					ByteCountSI(alltime), dur,
					ByteCountSI(alltime/dur*8),
					)

				total = 0
			}
		}
	}
}

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
	write(w, true)
}
