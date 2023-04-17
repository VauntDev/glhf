package glhf

import (
	"encoding/json"
	"io"
	"net/http"

	"google.golang.org/protobuf/proto"
)

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

type HandleFunc[I Body, O Body] func(*Request[I], *Response[O])

// DELETE deletes the specified resource. Optional request body
func Delete[I Body, O Body](fn HandleFunc[I, O], options ...Options) http.HandlerFunc {
	opts := defaultOptions()
	for _, opt := range options {
		opt.Apply(opts)
	}
	// http handler function called by server
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
		var responseBody O
		response := &Response[O]{w: w, Body: &responseBody}

		// call the handler
		fn(req, response)

		if response.Body != nil {
			b, err := marshalResponse(r.Header.Get(Accept), response.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(response.statusCode)
			if _, err := w.Write(b); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		w.WriteHeader(response.statusCode)
		return
	}
}

// Get requests a representation of the specified resource. No request body
func Get[I EmptyBody, O any](fn HandleFunc[I, O], options ...Options) http.HandlerFunc {
	opts := defaultOptions()
	for _, opt := range options {
		opt.Apply(opts)
	}
	// http handler function called by server
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		req := &Request[I]{r: r}
		var responseBody O
		response := &Response[O]{w: w, Body: &responseBody}

		// call the handler
		fn(req, response)

		if response.Body != nil {
			b, err := marshalResponse(r.Header.Get(Accept), response.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
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
	// http handler function called by server
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

			var responseBody O
			response := &Response[O]{w: w, Body: &responseBody}

			// call the handler
			fn(req, response)

			if response.Body != nil {
				b, err := marshalResponse(r.Header.Get(Accept), response.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(response.statusCode)
				if _, err := w.Write(b); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				return
			}
			w.WriteHeader(response.statusCode)
			return
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
	// http handler function called by server
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

		var responseBody O
		response := &Response[O]{w: w, Body: &responseBody}

		// call the handler
		fn(req, response)

		if response.Body != nil {
			b, err := marshalResponse(r.Header.Get(Accept), response.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(response.statusCode)
			if _, err := w.Write(b); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
		w.WriteHeader(response.statusCode)
		return
	}
}

// Put method is used to replace a resource with a similar resource that includes a different set of values. Requires request body
func Put[I Body, O Body](fn HandleFunc[I, O], options ...Options) http.HandlerFunc {
	opts := defaultOptions()
	for _, opt := range options {
		opt.Apply(opts)
	}
	// http handler function called by server
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

			var responseBody O
			response := &Response[O]{w: w, Body: &responseBody}

			// call the handler
			fn(req, response)

			if response.Body != nil {
				b, err := marshalResponse(r.Header.Get(Accept), response.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(response.statusCode)
				if _, err := w.Write(b); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				return
			}
			w.WriteHeader(response.statusCode)
			return
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
	default:

		// default applicaiton/json
		if err := json.Unmarshal(b, body); err != nil {
			return err
		}

		return nil
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
	default:
		// default applicaiton/json
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		return b, nil
	}
}
