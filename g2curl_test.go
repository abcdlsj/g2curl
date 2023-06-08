package g2curl

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestCURL(t *testing.T) {
	r := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "https",
			Host:   "example.com",
			Path:   "/",
		},
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(strings.NewReader(`{"foo":"bar"}`)),
	}

	c, err := New(r, Timeout(10))
	if err != nil {
		t.Fatal(err)
	}

	expected := `curl -m 10 -X 'POST' 'https://example.com/' -H 'Content-Type: application/json' -d '{"foo":"bar"}'`

	if c.String() != expected {
		t.Fatalf("expected: \n[%s]\ngot:\n [%s]", expected, c.String())
	}
}
