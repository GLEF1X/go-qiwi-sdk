package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/GLEF1X/go-qiwi-sdk/core/endpoints"
)

type Request struct {
	Payload     Payload
	APIEndpoint endpoints.Endpoint
	HttpMethod  string
}

type Payload struct {
	Headers          http.Header
	QueryParams      map[string]string
	URLConstructArgs []interface{}
	Body             interface{}
}

func (r Request) ConstructHTTPRequest(ctx context.Context, baseURL string, additionalHeaders map[string]string) (*http.Request, error) {
	baseURL += r.APIEndpoint.Resolve(r.Payload.URLConstructArgs)

	body, err := r.Payload.getEncodedBody()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, r.HttpMethod, baseURL, body)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	for key, value := range r.Payload.QueryParams {
		q.Add(key, value)
	}
	request.URL.RawQuery = q.Encode()

	headers := getDefaultHeaders()
	for key, value := range additionalHeaders {
		headers.Add(key, value)
	}
	request.Header = headers

	return request, nil
}

func (p Payload) getEncodedBody() (io.Reader, error) {
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(p.Body)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func getDefaultHeaders() http.Header {
	return http.Header{
		"Accept":       []string{"application/json"},
		"Host":         []string{"edge.qiwi.com"},
		"Content-Type": []string{"application/json"},
	}
}
