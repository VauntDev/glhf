package glhf

import "errors"

var (
	ErrProto                   = errors.New("value can not be used as proto message, invalid type")
	ErrUnsupportedResponseType = errors.New("response type unsupported")
	ErrUnsupportedRequestType  = errors.New("request type unsupported")
)

// errorResponse is a optional response that can be
// used to debug GLHF
type errorResponse struct {
	// Code is the HTTP Status code
	Code int `json:"code"`
	// Message is a developer-facing human-readable error message.
	Message string `json:"message"`
}
