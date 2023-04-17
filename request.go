package glhf

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
)

type contextKey int

const (
	varsKey contextKey = iota
)

type Request[T Body] struct {
	r    *http.Request
	body *T
}

func (req *Request[T]) HTTPRequest() *http.Request {
	return req.r
}

func (req *Request[T]) Body() (*T, error) {
	if req.body == nil {
		return nil, fmt.Errorf("nil body")
	}
	return req.body, nil
}

func (req *Request[T]) Header() http.Header {
	return req.r.Header
}

func (req *Request[T]) URL() *url.URL {
	return req.r.URL
}

func (req *Request[T]) Value(key any) any {
	return req.r.Context().Value(key)
}

func (req *Request[T]) Context() context.Context {
	return req.r.Context()
}

func (req *Request[T]) Cookie(name string) (*http.Cookie, error) {
	return req.r.Cookie(name)
}

func (req *Request[T]) Cookies() []*http.Cookie {
	return req.r.Cookies()
}

func (req *Request[T]) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return req.r.FormFile(key)
}

func (req *Request[T]) FormValue(key string) string {
	return req.r.FormValue(key)
}

func (req *Request[T]) MultipartReader() (*multipart.Reader, error) {
	return req.MultipartReader()
}

func (req *Request[T]) ParseForm() error {
	return req.r.ParseForm()
}

func (req *Request[T]) ParseMultipartForm(maxMemory int64) error {
	return req.r.ParseMultipartForm(maxMemory)
}

func (req *Request[T]) PostFormValue(key string) string {
	return req.r.PostFormValue(key)
}

func (req *Request[T]) ProtoAtLeast(major, minor int) bool {
	return req.r.ProtoAtLeast(major, minor)
}

func (req *Request[T]) Referer() string {
	return req.r.Referer()
}

func (req *Request[T]) UserAgent() string {
	return req.r.UserAgent()
}
