package glhf

import "net/http"

type Request[T any] struct {
	*http.Request
}

type Response[T any] struct {
	http.ResponseWriter
}

type HandleFunc[I any, O any] func(*Request[I]) *Response[O]

// Delete is a generic function that takes cares of the repeative marshaling logic for request/response objects
func Delete[I any, O any](fn HandleFunc[I, O]) http.HandlerFunc {
	// http handler function called by server
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

// Get is a generic function that takes cares of the repeative marshaling logic for request/response objects
func Get[I any, O any](fn HandleFunc[I, O]) http.HandlerFunc {
	// http handler function called by server
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

// Patch is a generic function that takes cares of the repeative marshaling logic for request/response objects
func Patch[I any, O any](fn HandleFunc[I, O]) http.HandlerFunc {
	// http handler function called by server
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

// Post is a generic function that takes cares of the repeative marshaling logic for request/response objects
func Post[I any, O any](fn HandleFunc[I, O]) http.HandlerFunc {
	// http handler function called by server
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

// Put is a generic function that takes cares of the repeative marshaling logic for request/response objects
func Put[I any, O any](fn HandleFunc[I, O], requiredParameters ...string) http.HandlerFunc {
	// http handler function called by server
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
