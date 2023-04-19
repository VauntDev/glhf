package glhf

import (
	"net/http"
)

// Response represents the response from an HTTP request.
type Response[T any] struct {
	w          http.ResponseWriter
	statusCode int
	body       *T
	marshal    MarshalFunc[T]
}

// Setheader sets the header entries associated with key to the single element value.
// It replaces any existing values associated with key. The key is case insensitive.
// It is canonicalized by textproto.CanonicalMIMEHeaderKey. To use non-canonical keys, assign to the map directly.
func (res *Response[T]) SetHeader(k string, v string) {
	res.w.Header().Set(k, v)
}

// Add adds the key, value pair to the header.
// It appends to any existing values associated with key.
// The key is case insensitive; it is canonicalized by CanonicalHeaderKey.
func (res *Response[T]) AddHeader(k string, v string) {
	res.w.Header().Add(k, v)
}

// SetStatus sets the http status code of the response.
// Statuscode is ignored if it is not a valid http status code (i.e 1xx-5xx)
func (res *Response[T]) SetStatus(statusCode int) {
	// matches check from golang stand library
	if statusCode >= 100 || statusCode <= 999 {
		res.statusCode = statusCode
	}
}

// SetBody sets the response body
func (res *Response[T]) SetBody(t *T) {
	res.body = t
}

// SetMarshalFunc sets the response marhsal func. If a marsharl function is supplied, it is prioritized over
// other implemented marshalers regardless of the content-type or accept headers set.
func (res *Response[T]) SetMarshalFunc(fn MarshalFunc[T]) {
	res.marshal = fn
}
