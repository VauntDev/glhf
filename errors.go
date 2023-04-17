package glhf

import "errors"

var (
	ProtoErr              = errors.New("object can not be used as proto message")
	MissingRequestBodyErr = errors.New("missing required request body")
)
