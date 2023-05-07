package http

// HTTP_Request object. Body should be a KV map.
type HTTP_Request struct {
	Method string
	URL    string
	Body   interface{}
}
