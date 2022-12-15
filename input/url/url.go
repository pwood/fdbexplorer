package url

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

var url *string

func init() {
	url = flag.String("url", "", "URL to fetch status json from periodically.")
}

func NewURL() (*URL, bool) {
	if len(*url) == 0 {
		return nil, false
	}

	return &URL{url: *url}, true
}

type URL struct {
	url string
}

func (f *URL) Status() (json.RawMessage, error) {
	if d, err := f.get(); err != nil {
		return nil, fmt.Errorf("url fetch err: %w", err)
	} else {
		return d, nil
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
