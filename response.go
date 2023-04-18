package glhf

import (
	"net/http"
)

// Response represents the response from an HTTP request.
type Response[T any] struct {
	w          http.ResponseWriter
	statusCode int
	Body       *T
}

// Header returns the underlying http response header object
func (res *Response[T]) Header() http.Header {
	return res.w.Header()
}

// SetStatus sets the http status code of the response
func (res *Response[T]) SetStatus(statusCode int) {
	res.statusCode = statusCode
}
