package glhf

import (
	"net/http"
)

type Response[T any] struct {
	w          http.ResponseWriter
	statusCode int
	Body       *T
}

func (res *Response[T]) Header() http.Header {
	return res.w.Header()
}

func (res *Response[T]) Status(statusCode int) {
	res.statusCode = statusCode
}
