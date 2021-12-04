package core

import (
	"context"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type WrappedHTTPClient struct {
	*http.Client
}

type Option func(*WrappedHTTPClient)

func NewHttpClient(opts ...Option) *WrappedHTTPClient {
	client := &WrappedHTTPClient{
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

func (c *WrappedHTTPClient) Close() {
	c.CloseIdleConnections()
}

func (c *WrappedHTTPClient) SendRequest(ctx context.Context, request *Request) (result []byte, err error) {
	httpRequest, err := request.Http(ctx)
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
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		if HasContentType(response, "application/json") {
			return result, nil
		}
	}
	return nil, HTTPError{Status: response.StatusCode}
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
