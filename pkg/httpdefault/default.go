package httpdefault

import (
	"net"
	"net/http"
	"time"
)

// UserAgent represents HTTP User-Agent header that should be added to requests to Selectel API.
// You should add this header by yourself or use Headers function.
const UserAgent = "external-dns-selectel-webhook"

const (
	// httpTimeout represents the default timeout (in seconds) for HTTP requests.
	httpTimeout = 30

	// dialTimeout represents the default timeout (in seconds) for HTTP connection establishments.
	dialTimeout = 60

	// keepaliveTimeout represents the default keep-alive period for an active network connection.
	keepaliveTimeout = 60

	// maxIdleConns represents the maximum number of idle (keep-alive) connections.
	maxIdleConns = 100

	// idleConnTimeout represents the maximum amount of time an idle (keep-alive) connection will remain
	// idle before closing itself.
	idleConnTimeout = 100

	// tlsHandshakeTimeout represents the default timeout (in seconds) for TLS handshake.
	tlsHandshakeTimeout = 60

	// expectContinueTimeout represents the default amount of time to wait for a server's first response headers.
	expectContinueTimeout = 5
)

// Client returns default HTTP client for requests to Selectel API. It does not add User-Agent header, so
// you should add it by yourself (by default it is UserAgent) or use Headers function.
func Client() http.Client {
	return http.Client{
		Timeout: httpTimeout * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   dialTimeout * time.Second,
				KeepAlive: keepaliveTimeout * time.Second,
			}).DialContext,
			MaxIdleConns:          maxIdleConns,
			IdleConnTimeout:       idleConnTimeout * time.Second,
			TLSHandshakeTimeout:   tlsHandshakeTimeout * time.Second,
			ExpectContinueTimeout: expectContinueTimeout * time.Second,
		},
	}
}

// Headers returns default HTTP headers for requests to Selectel API.
func Headers() http.Header {
	h := http.Header{}
	h.Add("User-Agent", UserAgent)

	return h
}
