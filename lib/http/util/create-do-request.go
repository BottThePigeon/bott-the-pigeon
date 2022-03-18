package utils

import (
	http "bott-the-pigeon/lib/http"

	"fmt"
	"io/ioutil"
)

// Wrapper for creating and doing a HTTP Request, returning the parsed response body.
// This is a HTTP (over)simplification to speed up basic API calls.
func CreateDoHTTPRequest(params http.HTTP_Request, headers map[string]string, successCode int) ([]byte, error) {
	req, err := CreateHTTPRequest(params)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	res, err := DoHTTPRequest(req, successCode)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	return resBody, nil
}