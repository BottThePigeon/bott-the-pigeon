package utils

import (
	c "bott-the-pigeon/lib/http/client"

	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// Make the provided HTTP request and validate its success using the provided mustStatus value.
func DoHTTPRequest(req *http.Request, successCode int) (*http.Response, error) {
	client := c.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed http request: %v", err)
	} else if resp.StatusCode != successCode {
		err := errors.New("got bad status code, expected " +
		strconv.Itoa(successCode) + ", got " + 
		resp.Status)
		return nil, fmt.Errorf("failed http request: %v", err)
	}
	return resp, nil
}