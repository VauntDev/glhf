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
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var requestBody I
		// check for request body
		if r.Body != nil || r.ContentLength >= 0 {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				return
			}

			if err := unmarshalRequest(r.Header.Get(ContentType), b, &requestBody); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		req := &Request[I]{r: r, body: &requestBody}
		response := &Response[O]{w: w, statusCode: http.StatusOK}

		// call the handler
		fn(req, response)

		if response.body != nil {
			if response.marshal != nil {
				b, err := response.marshal(*response.body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
				w.WriteHeader(response.statusCode)
				if _, err := w.Write(b); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				return
			}
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
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			w.WriteHeader(response.statusCode)
			if _, err := w.Write(b); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
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
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		req := &Request[I]{r: r}
		response := &Response[O]{w: w, statusCode: http.StatusOK}

		// call the handler
		fn(req, response)

		if response.body != nil {
			if response.marshal != nil {
				b, err := response.marshal(*response.body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
				w.WriteHeader(response.statusCode)
				if _, err := w.Write(b); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				return
			}

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
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			w.WriteHeader(response.statusCode)
			if _, err := w.Write(b); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusNoContent)
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
		if r.Method != http.MethodPatch {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var requestBody I

		// check for request body
		if r.Body != nil || r.ContentLength >= 0 {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				return
			}

			if err := unmarshalRequest(r.Header.Get(ContentType), b, &requestBody); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			req := &Request[I]{r: r, body: &requestBody}

			response := &Response[O]{w: w, statusCode: http.StatusOK}

			// call the handler
			fn(req, response)

			if response.body != nil {
				if response.marshal != nil {
					b, err := response.marshal(*response.body)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
					}
					w.WriteHeader(response.statusCode)
					if _, err := w.Write(b); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					return
				}
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
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
				w.WriteHeader(response.statusCode)
				if _, err := w.Write(b); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		} else {
			// missing request body
			w.WriteHeader(http.StatusBadRequest)
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
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var requestBody I
		// check for request body
		if r.Body != nil || r.ContentLength >= 0 {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				return
			}

			if err := unmarshalRequest(r.Header.Get(ContentType), b, &requestBody); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		req := &Request[I]{r: r, body: &requestBody}

		response := &Response[O]{w: w, statusCode: http.StatusOK}

		// call the handler
		fn(req, response)

		if response.body != nil {
			if response.marshal != nil {
				b, err := response.marshal(*response.body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
				w.WriteHeader(response.statusCode)
				if _, err := w.Write(b); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				return
			}

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
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			w.WriteHeader(response.statusCode)
			if _, err := w.Write(b); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
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
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var requestBody I
		// check for request body
		if r.Body != nil || r.ContentLength >= 0 {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				return
			}

			if err := unmarshalRequest(r.Header.Get(ContentType), b, &requestBody); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			req := &Request[I]{r: r, body: &requestBody}

			response := &Response[O]{w: w, statusCode: http.StatusOK}

			// call the handler
			fn(req, response)

			if response.body != nil {
				if response.marshal != nil {
					b, err := response.marshal(*response.body)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
					}
					w.WriteHeader(response.statusCode)
					if _, err := w.Write(b); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					return
				}
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
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
				w.WriteHeader(response.statusCode)
				if _, err := w.Write(b); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		} else {
			// missing request body
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func unmarshalRequest(contentType string, b []byte, body Body) error {
	switch contentType {
	case ContentProto:
		msg, ok := body.(proto.Message)
		if !ok {
			return ProtoErr
		}

		if err := proto.Unmarshal(b, msg); err != nil {
			return err
		}

		body, ok = msg.(Body)
		if !ok {
			return ProtoErr
		}
		return nil
	case ContentJSON:

		// default application/json
		if err := json.Unmarshal(b, body); err != nil {
			return err
		}

		return nil
	default:
		return UnsupportedRequestTypeErr
	}
}

func marshalResponse(contentType string, body Body) ([]byte, error) {
	switch contentType {
	case ContentProto:
		msg, ok := body.(proto.Message)
		if !ok {
			return nil, ProtoErr
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
		return nil, UnsupportedRequestTypeErr
	}
}
