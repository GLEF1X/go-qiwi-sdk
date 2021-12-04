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
	Payload     *Payload
	APIEndpoint endpoints.Endpoint
	HttpMethod  string
}

type Payload struct {
	Headers          *http.Header
	QueryParams      map[string]string
	URLConstructArgs []interface{}
	Body             interface{}
}

func (p *Payload) GetEncodedBody() (io.Reader, error) {
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(p.Body)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (r Request) GetUrl(baseURL string) string {
	return baseURL + r.APIEndpoint.Resolve(r.Payload.URLConstructArgs)
}

func (r Request) ConstructHTTPRequest(ctx context.Context, baseURL string, additionalHeaders map[string]string) (*http.Request, error) {
	var URLConstructArgs []interface{}

	if r.Payload != nil {
		URLConstructArgs = r.Payload.URLConstructArgs
	}
	baseURL += r.APIEndpoint.Resolve(URLConstructArgs)

	request, err := http.NewRequestWithContext(ctx, r.HttpMethod, baseURL, nil)
	if err != nil {
		return nil, err
	}

	if r.Payload != nil {
		if r.Payload.Body != nil {
			body, err := r.Payload.GetEncodedBody()
			if err != nil {
				return nil, err
			}
			switch v := body.(type) {
			case *bytes.Buffer:
				request.ContentLength = int64(v.Len())
				buf := v.Bytes()
				request.GetBody = func() (io.ReadCloser, error) {
					r := bytes.NewReader(buf)
					return io.NopCloser(r), nil
				}
			default:
			}
		}

		q := url.Values{}
		for key, value := range r.Payload.QueryParams {
			q.Add(key, value)
		}
		request.URL.RawQuery = q.Encode()
	}

	headers := GetDefaultHeaders()
	for key, value := range additionalHeaders {
		headers.Add(key, value)
	}

	request.Header = headers

	return request, nil
}

func GetDefaultHeaders() http.Header {
	return http.Header{
		"Accept":       []string{"application/json"},
		"Host":         []string{"edge.qiwi.com"},
		"Content-Type": []string{"application/json"},
	}
}
