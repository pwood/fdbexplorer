package url

import (
	"flag"
	"fmt"
	"github.com/pwood/fdbexplorer/data"
	"io"
	"net/http"
	"time"
)

var url *string

func init() {
	url = flag.String("url", "", "URL to fetch status json from periodically.")
}

func NewURL(ch chan data.State, interval time.Duration) (*URL, bool) {
	if len(*url) == 0 {
		return nil, false
	}

	return &URL{ch: ch, url: *url, interval: interval}, true
}

type URL struct {
	ch       chan data.State
	url      string
	interval time.Duration
}

func (f *URL) Run() {
	timer := time.NewTicker(f.interval)

	nowCh := make(chan struct{}, 1)
	nowCh <- struct{}{}

	for {
		select {
		case <-nowCh:
			f.poll()
		case <-timer.C:
			f.poll()
		}
	}
}

func (f *URL) poll() {
	start := time.Now()

	d, err := f.get()

	if err != nil {
		f.ch <- data.State{
			Err: fmt.Errorf("url fetch err: %w", err),
		}
		return
	}

	f.ch <- data.State{
		Duration: time.Now().Sub(start),
		Interval: f.interval,
		Data:     d,
	}
}

func (f *URL) get() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, f.url, nil)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response: not 200, was %d", res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("http body: %w", err)
	}

	return resBody, nil
}
