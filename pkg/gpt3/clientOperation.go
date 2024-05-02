package gpt3

import (
	"net/http"
	"time"
)

// ClientOption are options that can be passed when creating a new gpt client.
type ClientOption func(*client) error

// WithAPIVersion is a client option that allows you to override the default api version of the client.
//
//	@param apiVersion
//	@return ClientOption
func WithAPIVersion(apiVersion string) ClientOption {
	return func(c *client) error {
		c.apiVersion = apiVersion
		return nil
	}
}

// WithUserAgent is a client option that allows you to override the default user agent of the client.
//
//	@param userAgent
//	@return ClientOption
func WithUserAgent(userAgent string) ClientOption {
	return func(c *client) error {
		c.userAgent = userAgent
		return nil
	}
}

// WithHTTPClient allows you to override the internal http.Client used.
//
//	@param httpClient
//	@return ClientOption
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *client) error {
		c.httpClient = httpClient
		return nil
	}
}

// WithTimeout is a client option that allows you to override the default timeout duration of requests
// for the client.
//
//	@param timeout
//	@return ClientOption
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *client) error {
		c.httpClient.Timeout = timeout
		return nil
	}
}
