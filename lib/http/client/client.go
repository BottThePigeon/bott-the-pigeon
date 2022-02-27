package client

import (
	"net/http"
	"time"
)

// The HTTP client pointer is stored, and can be accessed later.
var client *http.Client

// Returns the stored HTTP client or creates one if it doesn't exist.
func GetHTTPClient() *http.Client {
	if client != nil {
		return client
	} else {
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.MaxIdleConns, t.MaxConnsPerHost, t.MaxIdleConnsPerHost = 100, 100, 100
		client = &http.Client{
			Timeout:   10 * time.Second,
			Transport: t,
		}
		return client
	}
}