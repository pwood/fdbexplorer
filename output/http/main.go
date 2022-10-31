package http

import (
	"flag"
	"github.com/pwood/fdbexplorer/data"
	"net/http"
)

var httpEnable *bool
var httpAddress *string

func init() {
	httpEnable = flag.Bool("http-enable", false, "If the http output should be enabled, making the `status json` output available on /status/json.")
	httpAddress = flag.String("http-address", "127.0.0.1:8080", "Host and port number for http server to listen on, using 0.0.0.0 for all interface bind.")
}

func NewHTTP(ch chan data.State) (*HTTP, bool) {
	if !*httpEnable {
		return nil, false
	}

	return &HTTP{ch: ch, address: *httpAddress}, true
}

type HTTP struct {
	ch      chan data.State
	address string
	data    []byte
}

func (h *HTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/status/json" && r.Method == "GET" {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(h.data)
	} else {
		http.NotFound(w, r)
	}
}

func (h *HTTP) Run() {
	go func() {
		err := http.ListenAndServe(h.address, h)
		if err != nil {
			panic(err)
		}
	}()

	for {
		for s := range h.ch {
			if s.Err == nil {
				h.data = s.Data
			}
		}
	}
}
