package client

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
