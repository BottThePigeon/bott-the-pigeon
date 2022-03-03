package utils

import (
	"fmt"
	"io/ioutil"
)

// HTTP_Request object. Body should be a KV map.
type HTTP_Request struct {
	Method string
	URL    string
	Body   interface{}
}

// Wrapper for creating and doing a HTTP Request, returning the parsed response body.
// This is a HTTP (over)simplification to speed up basic API calls.
func CreateDoHTTPRequest(params HTTP_Request, headers map[string]string, successCode int) ([]byte, error) {
	req, err := CreateHTTPRequest(params.Method, params.URL, params.Body)
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