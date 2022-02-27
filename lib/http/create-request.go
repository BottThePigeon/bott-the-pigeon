package httputils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Wraps http.NewRequest with the removal of the need for a premade buffer.
func CreateHTTPRequest(method string, url string, body interface{}) (*http.Request, error) {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request JSON body: %v", err)
	}
	bodyBuffer := bytes.NewBuffer(bodyJson)
	req, err := http.NewRequest(method, url, bodyBuffer)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}
	return req, nil
}