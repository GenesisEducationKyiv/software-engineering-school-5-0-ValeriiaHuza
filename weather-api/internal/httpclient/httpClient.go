package httpclient

import (
	"net"
	"net/http"
	"time"
)

func InitHttpClient() http.Client {
	return http.Client{
		Timeout: 2 * time.Second,

		Transport: &http.Transport{
			// Max idle connections
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,

			// Keep connections open
			IdleConnTimeout: 90 * time.Second,

			// Optional: TLS handshake, dial timeouts
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,

			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		},
	}
}
