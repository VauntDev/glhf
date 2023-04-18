package glhf

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/url"
)

// A Request represents an HTTP request received by a server
type Request[T Body] struct {
	r    *http.Request
	body *T
}

// HTTPRequest returns the raw HTTP Request. If the request contains a body,
// it will be nil and only available through the glhf request.Body() method.
// This method is a temporary function and will be deprecated in future releases.
func (req *Request[T]) HTTPRequest() *http.Request {
	return req.r
}

// Body returns the generic http body, or nil if no body exists.
func (req *Request[T]) Body() (*T, error) {
	if req.body == nil {
		return nil, MissingRequestBodyErr
	}
	return req.body, nil
}

// Header wraps http.Request.Header
func (req *Request[T]) Header() http.Header {
	return req.r.Header
}

// URL wraps http.Request.URL
func (req *Request[T]) URL() *url.URL {
	return req.r.URL
}

// Value returns the value associated with the request context for key, or nil if no value is associated with key.
// Successive calls to Value with the same key returns the same result.
func (req *Request[T]) Value(key any) any {
	return req.r.Context().Value(key)
}

// Context returns the underlying request's context.
func (req *Request[T]) Context() context.Context {
	return req.r.Context()
}

// Cookie warps the underlying request's cookie and returns the named cookie provided in the request or ErrNoCookie if not found.
// If multiple cookies match the given name, only one cookie will be returned.
func (req *Request[T]) Cookie(name string) (*http.Cookie, error) {
	return req.r.Cookie(name)
}

// Cookies wraps the underlying request and returns the HTTP cookies sent with the request.
func (req *Request[T]) Cookies() []*http.Cookie {
	return req.r.Cookies()
}

//FormFile wraps the underlying request and returns the first file for the provided form key.
// FormFile calls ParseMultipartForm and ParseForm if necessary.
func (req *Request[T]) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return req.r.FormFile(key)
}

// FormValue wraps the underlying request and returns the first value for the named component of the query.
// POST and PUT body parameters take precedence over URL query string values.
// FormValue calls ParseMultipartForm and ParseForm if necessary and ignores any errors returned by these functions.
// If key is not present, FormValue returns the empty string. To access multiple values of the same key, call ParseForm and then inspect Request.Form directly.
func (req *Request[T]) FormValue(key string) string {
	return req.r.FormValue(key)
}

// MultipartReader wraps the underlying request and returns a MIME multipart reader
// if this is a multipart/form-data or a multipart/mixed POST request, else returns nil and an error.
// Use this function instead of ParseMultipartForm to process the request body as a stream.
func (req *Request[T]) MultipartReader() (*multipart.Reader, error) {
	return req.r.MultipartReader()
}

// ParseForm  wraps the underlying request and populates r.Form and r.PostForm.
func (req *Request[T]) ParseForm() error {
	return req.r.ParseForm()
}

//ParseMultipartForm wraps the underlying request and parses a request body as multipart/form-data.
// The whole request body is parsed and up to a total of maxMemory bytes of its file parts are stored in memory, with the remainder stored on disk in temporary files.
// ParseMultipartForm calls ParseForm if necessary. If ParseForm returns an error, ParseMultipartForm returns it but also continues parsing the request body.
// After one call to ParseMultipartForm, subsequent calls have no effect.
func (req *Request[T]) ParseMultipartForm(maxMemory int64) error {
	return req.r.ParseMultipartForm(maxMemory)
}

// PostFormValue wraps the underlying request and returns the first value for the named component of the POST, PATCH, or PUT request body. URL query parameters are ignored.
// PostFormValue calls ParseMultipartForm and ParseForm if necessary and ignores any errors returned by these functions.
// If key is not present, PostFormValue returns the empty string.
func (req *Request[T]) PostFormValue(key string) string {
	return req.r.PostFormValue(key)
}

// ProtoAtLeast wraps the underlying request and reports whether the HTTP protocol used in the request is at least major.minor.
func (req *Request[T]) ProtoAtLeast(major, minor int) bool {
	return req.r.ProtoAtLeast(major, minor)
}

// Referer wraps the underlying request and returns the referring URL, if sent in the request.
func (req *Request[T]) Referer() string {
	return req.r.Referer()
}

// UserAgent wraps the underlying request and
// returns the client's User-Agent, if sent in the request.
func (req *Request[T]) UserAgent() string {
	return req.r.UserAgent()
}
