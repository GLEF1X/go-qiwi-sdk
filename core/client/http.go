package client

import (
	"context"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/GLEF1X/go-qiwi-sdk/core"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Http struct {
	*http.Client
	baseURL        string
	defaultHeaders map[string]string
}

type Option func(*Http)

func WithBaseURL(url string) func(*Http) {
	return func(h *Http) {
		h.baseURL = url
	}
}

func WithDefaultHeaders(headers map[string]string) func(*Http) {
	return func(h *Http) {
		h.defaultHeaders = headers
	}
}

func NewHttp(opts ...Option) *Http {
	client := &Http{
		Client: &http.Client{
			Transport: &http.Transport{
				DisableCompression: true,
				IdleConnTimeout:    30 * time.Second,
				MaxIdleConns:       150,
			},
		},
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (c *Http) Close() {
	c.CloseIdleConnections()
}

func (c *Http) SendRequest(ctx context.Context, request *Request) (result []byte, err error) {
	httpRequest, err := request.ConstructHTTPRequest(ctx, c.baseURL, c.defaultHeaders)
	if err != nil {
		return nil, err
	}
	response, err := c.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	if result, err = getResponseBody(response); err != nil {
		return nil, err
	}
	defer func() {
		closeErr := response.Body.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()
	return result, nil
}

func getResponseBody(response *http.Response) ([]byte, error) {
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if ResponseIsUnsatisfactory(response) {
		return nil, core.ErrBadResponse{Status: response.StatusCode}
	}
	return result, nil
}

func ResponseIsUnsatisfactory(r *http.Response) bool {
	if !(r.StatusCode >= 200 && r.StatusCode < 300) {
		return false
	} else if !HasContentType(r, "application/json") {
		return false
	}
	return true
}

// HasContentType determines whether the request `content-type` includes a
// server-acceptable mime-type
func HasContentType(r *http.Response, mimetype string) bool {
	contentType := r.Header.Get("Content-type")
	if contentType == "" {
		return mimetype == "application/octet-stream"
	}

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == mimetype {
			return true
		}
	}
	return false
}
