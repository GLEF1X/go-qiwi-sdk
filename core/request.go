package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/GLEF1X/go-qiwi-sdk/core/endpoints"
)

var headersPool = sync.Pool{
	New: func() interface{} {
		return http.Header{
			"Accept":       []string{"application/json"},
			"Host":         []string{"edge.qiwi.com"},
			"Content-Type": []string{"application/json"},
		}
	},
}

type Request struct {
	Payload            *Payload
	BaseUrl            string
	APIEndpoint        endpoints.Endpoint
	HttpMethod         string
	AuthorizationToken string
}

type Payload struct {
	Headers          *http.Header
	QueryParams      map[string]string
	URLConstructArgs []interface{}
	Body             interface{}
}

func (p *Payload) GetBody() (io.Reader, error) {
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(p.Body)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (r Request) GetUrl() string {
	return r.BaseUrl + r.APIEndpoint.Resolve(r.Payload.URLConstructArgs)
}

func (r Request) Http(ctx context.Context) (*http.Request, error) {
	endpoint := r.BaseUrl + r.APIEndpoint.Resolve(r.Payload.URLConstructArgs)
	body, err := r.Payload.GetBody()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(r.HttpMethod, endpoint, body)
	if err != nil {
		return nil, err
	}
	request.Header = createDefaultHeaders(r.AuthorizationToken)
	q := url.Values{}
	for key, value := range r.Payload.QueryParams {
		q.Add(key, value)
	}
	request.URL.RawQuery = q.Encode()
	request = request.WithContext(ctx)
	return request, nil
}

func createDefaultHeaders(token string) http.Header {
	authorizationHeaderValue := fmt.Sprintf("Bearer %s", token)
	headers := headersPool.Get().(http.Header)
	defer headersPool.Put(headers)
	headers.Add("Authorization", authorizationHeaderValue)
	return headers
}
