package utils

import (
	httputil "bott-the-pigeon/lib/http"

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Wraps http.NewRequest with removal of the need for a premade buffer.
func CreateHTTPRequest(params httputil.HTTP_Request) (*http.Request, error) {
	bodyJson, err := json.Marshal(params.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request JSON body: %v", err)
	}
	bodyBuffer := bytes.NewBuffer(bodyJson)
	req, err := http.NewRequest(params.Method, params.URL, bodyBuffer)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}
	return req, nil
}
