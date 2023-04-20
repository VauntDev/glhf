package glhf

import (
	"encoding/json"
	"io"
	"net/http"

	"google.golang.org/protobuf/proto"
)

// Body is the request's body.
type Body any

// EmptyBody represents an empty http request body
type EmptyBody struct{}

const (
	// ContentType header constant.
	ContentType = "Content-Type"
	Accept      = "Accept"

	// ContentJSON header value for JSON data.
	ContentJSON = "application/json"
	// ContentProto header value for proto buff
	ContentProto = "application/proto"

	// TODO :: Add additional content type support
	// ContentBinary header value for binary data.
	ContentBinary = "application/octet-stream"
	// ContentHTML header value for HTML data.
	ContentHTML = "text/html"
	// ContentText header value for Text data.
	ContentText = "text/plain"
	// ContentXHTML header value for XHTML data.
	ContentXHTML = "application/xhtml+xml"
	// ContentXML header value for XML data.
	ContentXML = "text/xml"
)

// HandleFunc responds to an HTTP request.
// I and O represent the request body or response body.
type HandleFunc[I Body, O Body] func(*Request[I], *Response[O])

// MarshalFunc defines how a body should be marshaled into bytes
type MarshalFunc[I Body] func(I) ([]byte, error)

// Delete deletes the specified resource. The underlying request body is optional.
func Delete[I Body, O Body](fn HandleFunc[I, O], options ...Options) http.HandlerFunc {
	opts := defaultOptions()
	for _, opt := range options {
		opt.Apply(opts)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var errResp *errorResponse
		if r.Method != http.MethodDelete {
			errResp = &errorResponse{
				Code:    http.StatusMethodNotAllowed,
				Message: "invalid method used, expected DELETE found " + r.Method,
			}
		}
		var requestBody I
		// check for request body
		if r.Body != nil || r.ContentLength >= 0 {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				errResp = &errorResponse{
					Code:    http.StatusInternalServerError,
					Message: "failed to read request body",
				}
			}

			if err := unmarshalRequest(r.Header.Get(ContentType), b, &requestBody); err != nil {
				errResp = &errorResponse{
					Code:    http.StatusInternalServerError,
					Message: "failed to unmarshal request with content-type " + r.Header.Get(ContentType),
				}
			}
		}

		// request failed to unmarshal, return with failure
		if errResp != nil {
			w.WriteHeader(errResp.Code)
			if opts.verbose {
				b, _ := json.Marshal(errResp)
				w.Write(b)
			}
			return
		}

		req := &Request[I]{r: r, body: &requestBody}
		response := &Response[O]{w: w}

		// call the handler
		fn(req, response)

		if response.body != nil {
			var bodyBytes []byte
			// if there is a custom marshaler, prioritize it
			if response.marshal != nil {
				b, err := response.marshal(*response.body)
				if err != nil {
					errResp = &errorResponse{
						Code:    http.StatusInternalServerError,
						Message: "failed to marshal response with custom marhsaler",
					}
				}
				bodyBytes = b
			} else {
				// client preferred content-type
				b, err := marshalResponse(r.Header.Get(Accept), response.body)
				if err != nil {
					// server preferred content-type
					contentType := response.w.Header().Get(ContentType)
					if len(contentType) == 0 {
						contentType = opts.defaultContentType
					}
					b, err = marshalResponse(contentType, response.body)
					if err != nil {
						errResp = &errorResponse{
							Code:    http.StatusInternalServerError,
							Message: "failed to marshal response with content-type: " + contentType,
						}
					}
				}
				bodyBytes = b
			}
			// Response failed to marshal
			if errResp != nil {
				w.WriteHeader(errResp.Code)
				if opts.verbose {
					b, _ := json.Marshal(errResp)
					w.Write(b)
				}
				return
			}
			// ensure user supplied status code is valid
			if validStatusCode(response.statusCode) {
				w.WriteHeader(response.statusCode)
			}
			if len(bodyBytes) > 0 {
				w.Write(bodyBytes)
			}
			return
		}
	}
}

// Get requests a representation of the specified resource. Expects an empty request body. If a request
// body is set, it will be ignored.
func Get[I EmptyBody, O any](fn HandleFunc[I, O], options ...Options) http.HandlerFunc {
	opts := defaultOptions()
	for _, opt := range options {
		opt.Apply(opts)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var errResp *errorResponse
		if r.Method != http.MethodGet {
			errResp = &errorResponse{
				Code:    http.StatusMethodNotAllowed,
				Message: "invalid method used, expected GET found " + r.Method,
			}
		}

		req := &Request[I]{r: r}
		response := &Response[O]{w: w, statusCode: http.StatusOK}

		// call the handler
		fn(req, response)

		if response.body != nil {
			var bodyBytes []byte
			// if there is a custom marshaler, prioritize it
			if response.marshal != nil {
				b, err := response.marshal(*response.body)
				if err != nil {
					errResp = &errorResponse{
						Code:    http.StatusInternalServerError,
						Message: "failed to marshal response with custom marhsaler",
					}
				}
				bodyBytes = b
			} else {
				// client preferred content-type
				b, err := marshalResponse(r.Header.Get(Accept), response.body)
				if err != nil {
					// server preferred content-type
					contentType := response.w.Header().Get(ContentType)
					if len(contentType) == 0 {
						contentType = opts.defaultContentType
					}
					b, err = marshalResponse(contentType, response.body)
					if err != nil {
						errResp = &errorResponse{
							Code:    http.StatusInternalServerError,
							Message: "failed to marshal response with content-type: " + contentType,
						}
					}
				}
				bodyBytes = b
			}
			// Response failed to marshal
			if errResp != nil {
				w.WriteHeader(errResp.Code)
				if opts.verbose {
					b, _ := json.Marshal(errResp)
					w.Write(b)
				}
				return
			}
			// ensure user supplied status code is valid
			if validStatusCode(response.statusCode) {
				w.WriteHeader(response.statusCode)
			}
			if len(bodyBytes) > 0 {
				w.Write(bodyBytes)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
			return
		}
	}
}

// Patch method is used to apply partial modifications to a resource. Required Request Body
func Patch[I Body, O Body](fn HandleFunc[I, O], options ...Options) http.HandlerFunc {
	opts := defaultOptions()
	for _, opt := range options {
		opt.Apply(opts)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var errResp *errorResponse
		if r.Method != http.MethodPatch {
			errResp = &errorResponse{
				Code:    http.StatusMethodNotAllowed,
				Message: "invalid method used, expected PATCH found " + r.Method,
			}
		}
		var requestBody I
		// check for request body
		if r.Body != nil || r.ContentLength >= 0 {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				errResp = &errorResponse{
					Code:    http.StatusInternalServerError,
					Message: "failed to read request body",
				}
			}

			if err := unmarshalRequest(r.Header.Get(ContentType), b, &requestBody); err != nil {
				errResp = &errorResponse{
					Code:    http.StatusInternalServerError,
					Message: "failed to unmarshal request with content-type " + r.Header.Get(ContentType),
				}
			}
		} else {
			errResp = &errorResponse{
				Code:    http.StatusBadRequest,
				Message: "missing request body",
			}
		}

		// request failed to unmarshal, return with failure
		if errResp != nil {
			w.WriteHeader(errResp.Code)
			if opts.verbose {
				b, _ := json.Marshal(errResp)
				w.Write(b)
			}
			return
		}

		req := &Request[I]{r: r, body: &requestBody}
		response := &Response[O]{w: w, statusCode: http.StatusOK}

		// call the handler
		fn(req, response)

		if response.body != nil {
			var bodyBytes []byte
			// if there is a custom marshaler, prioritize it
			if response.marshal != nil {
				b, err := response.marshal(*response.body)
				if err != nil {
					errResp = &errorResponse{
						Code:    http.StatusInternalServerError,
						Message: "failed to marshal response with custom marhsaler",
					}
				}
				bodyBytes = b
			} else {
				// client preferred content-type
				b, err := marshalResponse(r.Header.Get(Accept), response.body)
				if err != nil {
					// server preferred content-type
					contentType := response.w.Header().Get(ContentType)
					if len(contentType) == 0 {
						contentType = opts.defaultContentType
					}
					b, err = marshalResponse(contentType, response.body)
					if err != nil {
						errResp = &errorResponse{
							Code:    http.StatusInternalServerError,
							Message: "failed to marshal response with content-type: " + contentType,
						}
					}
				}
				bodyBytes = b
			}
			// Response failed to marshal
			if errResp != nil {
				w.WriteHeader(errResp.Code)
				if opts.verbose {
					b, _ := json.Marshal(errResp)
					w.Write(b)
				}
				return
			}
			// ensure user supplied status code is valid
			if validStatusCode(response.statusCode) {
				w.WriteHeader(response.statusCode)
			}
			if len(bodyBytes) > 0 {
				w.Write(bodyBytes)
			}
			return
		}
	}
}

// Post method can be used in two different ways, create a resource or perform and operation:. Optional request body
func Post[I Body, O Body](fn HandleFunc[I, O], options ...Options) http.HandlerFunc {
	opts := defaultOptions()
	for _, opt := range options {
		opt.Apply(opts)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var errResp *errorResponse
		if r.Method != http.MethodPost {
			errResp = &errorResponse{
				Code:    http.StatusMethodNotAllowed,
				Message: "invalid method used, expected POST found " + r.Method,
			}
		}
		var requestBody I
		// check for request body
		if r.Body != nil || r.ContentLength >= 0 {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				errResp = &errorResponse{
					Code:    http.StatusInternalServerError,
					Message: "failed to read request body",
				}
			}

			if err := unmarshalRequest(r.Header.Get(ContentType), b, &requestBody); err != nil {
				errResp = &errorResponse{
					Code:    http.StatusInternalServerError,
					Message: "failed to unmarshal request with content-type " + r.Header.Get(ContentType),
				}
			}
		}

		// request failed to unmarshal, return with failure
		if errResp != nil {
			w.WriteHeader(errResp.Code)
			if opts.verbose {
				b, _ := json.Marshal(errResp)
				w.Write(b)
			}
			return
		}

		req := &Request[I]{r: r, body: &requestBody}

		response := &Response[O]{w: w, statusCode: http.StatusOK}

		// call the handler
		fn(req, response)

		if response.body != nil {
			var bodyBytes []byte
			// if there is a custom marshaler, prioritize it
			if response.marshal != nil {
				b, err := response.marshal(*response.body)
				if err != nil {
					errResp = &errorResponse{
						Code:    http.StatusInternalServerError,
						Message: "failed to marshal response with custom marhsaler",
					}
				}
				bodyBytes = b
			} else {
				// client preferred content-type
				b, err := marshalResponse(r.Header.Get(Accept), response.body)
				if err != nil {
					// server preferred content-type
					contentType := response.w.Header().Get(ContentType)
					if len(contentType) == 0 {
						contentType = opts.defaultContentType
					}
					b, err = marshalResponse(contentType, response.body)
					if err != nil {
						errResp = &errorResponse{
							Code:    http.StatusInternalServerError,
							Message: "failed to marshal response with content-type: " + contentType,
						}
					}
				}
				bodyBytes = b
			}
			// Response failed to marshal
			if errResp != nil {
				w.WriteHeader(errResp.Code)
				if opts.verbose {
					b, _ := json.Marshal(errResp)
					w.Write(b)
				}
				return
			}
			// ensure user supplied status code is valid
			if validStatusCode(response.statusCode) {
				w.WriteHeader(response.statusCode)
			}
			if len(bodyBytes) > 0 {
				w.Write(bodyBytes)
			}
			return
		}
	}
}

// Put method is used to replace a resource with a similar resource that includes a different set of values. Requires request body
func Put[I Body, O Body](fn HandleFunc[I, O], options ...Options) http.HandlerFunc {
	opts := defaultOptions()
	for _, opt := range options {
		opt.Apply(opts)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var errResp *errorResponse
		if r.Method != http.MethodPut {
			errResp = &errorResponse{
				Code:    http.StatusMethodNotAllowed,
				Message: "invalid method used, expected PUT found " + r.Method,
			}
		}
		var requestBody I
		// check for request body
		if r.Body != nil || r.ContentLength >= 0 {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				errResp = &errorResponse{
					Code:    http.StatusInternalServerError,
					Message: "failed to read request body",
				}
			}

			if err := unmarshalRequest(r.Header.Get(ContentType), b, &requestBody); err != nil {
				errResp = &errorResponse{
					Code:    http.StatusInternalServerError,
					Message: "failed to unmarshal request with content-type " + r.Header.Get(ContentType),
				}
			}
		} else {
			errResp = &errorResponse{
				Code:    http.StatusBadRequest,
				Message: "missing request body",
			}
		}

		// request failed to unmarshal, return with failure
		if errResp != nil {
			w.WriteHeader(errResp.Code)
			if opts.verbose {
				b, _ := json.Marshal(errResp)
				w.Write(b)
			}
			return
		}

		req := &Request[I]{r: r, body: &requestBody}
		response := &Response[O]{w: w, statusCode: http.StatusOK}

		// call the handler
		fn(req, response)

		if response.body != nil {
			var bodyBytes []byte
			// if there is a custom marshaler, prioritize it
			if response.marshal != nil {
				b, err := response.marshal(*response.body)
				if err != nil {
					errResp = &errorResponse{
						Code:    http.StatusInternalServerError,
						Message: "failed to marshal response with custom marhsaler",
					}
				}
				bodyBytes = b
			} else {
				// client preferred content-type
				b, err := marshalResponse(r.Header.Get(Accept), response.body)
				if err != nil {
					// server preferred content-type
					contentType := response.w.Header().Get(ContentType)
					if len(contentType) == 0 {
						contentType = opts.defaultContentType
					}
					b, err = marshalResponse(contentType, response.body)
					if err != nil {
						errResp = &errorResponse{
							Code:    http.StatusInternalServerError,
							Message: "failed to marshal response with content-type: " + contentType,
						}
					}
				}
				bodyBytes = b
			}
			// Response failed to marshal
			if errResp != nil {
				w.WriteHeader(errResp.Code)
				if opts.verbose {
					b, _ := json.Marshal(errResp)
					w.Write(b)
				}
				return
			}
			// ensure user supplied status code is valid
			if validStatusCode(response.statusCode) {
				w.WriteHeader(response.statusCode)
			}
			if len(bodyBytes) > 0 {
				w.Write(bodyBytes)
			}
			return
		}
	}
}

func unmarshalRequest(contentType string, b []byte, body Body) error {
	switch contentType {
	case ContentProto:
		// msg pointer matches body
		msg, ok := body.(proto.Message)
		if !ok {
			return ErrProto
		}

		if err := proto.Unmarshal(b, msg); err != nil {
			return err
		}

		return nil
	case ContentJSON:

		// default application/json
		if err := json.Unmarshal(b, body); err != nil {
			return err
		}

		return nil
	default:
		return ErrUnsupportedRequestType
	}
}

func marshalResponse(contentType string, body Body) ([]byte, error) {
	switch contentType {
	case ContentProto:
		msg, ok := body.(proto.Message)
		if !ok {
			return nil, ErrProto
		}

		b, err := proto.Marshal(msg)
		if err != nil {
			return nil, err
		}
		return b, nil
	case ContentJSON:
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		return b, nil
	default:
		return nil, ErrUnsupportedRequestType
	}
}

func validStatusCode(statusCode int) bool {
	return (statusCode >= 100 && statusCode <= 999)
}
