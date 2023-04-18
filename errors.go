package glhf

import "errors"

var (
	ProtoErr                   = errors.New("object can not be used as proto message")
	MissingRequestBodyErr      = errors.New("missing required request body")
	UnsupportedResponseTypeErr = errors.New("response type unsupported")
	UnsupportedRequestTypeErr  = errors.New("request type unsupported")
)
