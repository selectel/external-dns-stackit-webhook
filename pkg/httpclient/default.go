package httpclient

import (
	"net"
	"net/http"
	"time"
)

const (

	// DefaultUserAgent of CLI represents HTTP User-Agent header that will be added to auth requests.
	DefaultUserAgent = "external-dns-selectel-webhook"

	// defaultHTTPTimeout represents the default timeout (in seconds) for HTTP requests.
	defaultHTTPTimeout = 30

	// defaultDialTimeout represents the default timeout (in seconds) for HTTP connection establishments.
	defaultDialTimeout = 60

	// defaultKeepaliveTimeout represents the default keep-alive period for an active network connection.
	defaultKeepaliveTimeout = 60

	// defaultMaxIdleConns represents the maximum number of idle (keep-alive) connections.
	defaultMaxIdleConns = 100

	// defaultIdleConnTimeout represents the maximum amount of time an idle (keep-alive) connection will remain
	// idle before closing itself.
	defaultIdleConnTimeout = 100

	// defaultTLSHandshakeTimeout represents the default timeout (in seconds) for TLS handshake.
	defaultTLSHandshakeTimeout = 60

	// defaultExpectContinueTimeout represents the default amount of time to wait for a server's first
	// response headers.
	defaultExpectContinueTimeout = 1
)

// Default sets up default http client for authentication.
func Default() http.Client {
	return http.Client{
		Timeout: defaultHTTPTimeout * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   defaultDialTimeout * time.Second,
				KeepAlive: defaultKeepaliveTimeout * time.Second,
			}).DialContext,
			MaxIdleConns:          defaultMaxIdleConns,
			IdleConnTimeout:       defaultIdleConnTimeout * time.Second,
			TLSHandshakeTimeout:   defaultTLSHandshakeTimeout * time.Second,
			ExpectContinueTimeout: defaultExpectContinueTimeout * time.Second,
		},
	}
}
