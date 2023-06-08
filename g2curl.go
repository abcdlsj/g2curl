package g2curl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var longFormatOptionsMap = map[string]string{
	"-X": "--request",
	"-H": "--header",
	"-d": "--data",
	"-m": "--max-time",
	"-x": "--proxy",
}

type CURL struct {
	url    string
	body   string
	method string

	header http.Header

	multiLine      bool
	longFormat     bool
	timeout        int
	followRedirect bool
	ignoreTLS      bool
	proxy          string
	outputFile     string
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

func Format(longFormat, multiLine bool) Option {
	return func(c *CURL) {
		c.longFormat = longFormat
		c.multiLine = multiLine
	}
}

func Timeout(timeout int) Option {
	return func(c *CURL) {
		c.timeout = timeout
	}
}

func FollowRedirect() Option {
	return func(c *CURL) {
		c.followRedirect = true
	}
}

func Proxy(proxy string) Option {
	return func(c *CURL) {
		c.proxy = proxy
	}
}

func IgnoreTLS() Option {
	return func(c *CURL) {
		c.ignoreTLS = true
	}
}

func Output(file string) Option {
	return func(c *CURL) {
		c.outputFile = file
	}
}

func (c *CURL) getOptionFormat(s string) string {
	if c.longFormat {
		if longFormatOption, ok := longFormatOptionsMap[s]; ok {
			return longFormatOption
		}
	}
	return s
}

func (c *CURL) String() string {
	if c == nil {
		return ""
	}

	curl := []string{"curl"}

	if c.proxy != "" {
		curl = append(curl, c.getOptionFormat("-x"), bashStr(c.proxy))
	}

	if c.followRedirect {
		curl = append(curl, c.getOptionFormat("-L"))
	}

	if c.timeout > 0 {
		curl = append(curl, c.getOptionFormat("-m"), fmt.Sprintf("%d", c.timeout))
	}

	if c.ignoreTLS {
		curl = append(curl, c.getOptionFormat("-k"))
	}

	curl = append(curl, c.getOptionFormat("-X"), bashStr(c.method))

	curl = append(curl, bashStr(c.url))

	for k, v := range c.header {
		curl = append(curl, c.getOptionFormat("-H"), bashStr(fmt.Sprintf("%s: %s", k, v[0])))
	}

	if c.body != "" {
		curl = append(curl, c.getOptionFormat("-d"), bashStr(c.body))
	}

	return strings.Join(curl, " ")
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
