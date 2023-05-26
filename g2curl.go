package g2curl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type CURL struct {
	tlsFlag bool
	url     string
	header  http.Header
	body    string
	method  string
}

func (c *CURL) String() string {
	if c == nil {
		return ""
	}

	curl := []string{"curl"}

	if c.tlsFlag {
		curl = append(curl, "-k")
	}

	curl = append(curl, "-X", bashStr(c.method))

	curl = append(curl, bashStr(c.url))

	for k, v := range c.header {
		curl = append(curl, "-H", bashStr(fmt.Sprintf("%s: %s", k, v[0])))
	}

	if c.body != "" {
		curl = append(curl, "-d", bashStr(c.body))
	}

	curl = append(curl, "--compressed")

	return strings.Join(curl, " ")
}

type Option func(*CURL)

func New(r *http.Request, opts ...Option) (*CURL, error) {
	curl, err := build(r)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(curl)
	}

	return curl, nil
}

func bashStr(s string) string {
	return fmt.Sprintf("'%s'", strings.Replace(s, "'", "'\"'\"'", -1))
}

func build(r *http.Request) (*CURL, error) {
	c := CURL{
		method: r.Method,
		header: r.Header,
	}
	schema := r.URL.Scheme
	requestURL := r.URL.String()
	if schema == "" {
		schema = "http"
		if r.TLS != nil {
			schema = "https"
		}
		requestURL = schema + "://" + r.Host + r.URL.Path
	}

	c.url = requestURL

	if schema == "https" {
		c.tlsFlag = true
	}

	if r.Body != nil {
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			return nil, err
		}
		r.Body = io.NopCloser(bytes.NewBuffer(buf.Bytes()))
		if len(buf.String()) > 0 {
			c.body = buf.String()
		}
	}

	return &c, nil
}
