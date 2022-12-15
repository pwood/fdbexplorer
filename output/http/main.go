package http

import (
	"flag"
	"github.com/pwood/fdbexplorer/input"
	"net/http"
)

var httpEnable *bool
var httpAddress *string

func init() {
	httpEnable = flag.Bool("http-enable", false, "If the http output should be enabled, making the `status json` output available on /status/json.")
	httpAddress = flag.String("http-address", "127.0.0.1:8080", "Host and port number for http server to listen on, using 0.0.0.0 for all interface bind.")
}

func NewHTTP(ds input.StatusProvider) (*HTTP, bool) {
	if !*httpEnable {
		return nil, false
	}

	return &HTTP{ds: ds, address: *httpAddress}, true
}

type HTTP struct {
	ds      input.StatusProvider
	address string
}

func (h *HTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/status/json" && r.Method == "GET" {
		if d, err := h.ds.Status(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Add("content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(d)
		}
	} else {
		http.NotFound(w, r)
	}
}

func (h *HTTP) Run() {
	err := http.ListenAndServe(h.address, h)
	if err != nil {
		panic(err)
	}
}
