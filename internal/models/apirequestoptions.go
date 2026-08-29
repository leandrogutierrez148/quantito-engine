package models

import "net/url"

type APIRequestOptions struct {
	Query     url.Values
	Headers   map[string]string
	Body      []byte
	Auth      bool
	BasicAuth string
}

type APIClientOptionsFunc func(*APIRequestOptions) error
