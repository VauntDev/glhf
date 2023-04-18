# glhf
Generic Light-weight Handler Framework


## Experimental

This is a **experimental** library. The goal is to evaluate leveraging generics to reduce duplicate code around deserializing and serializing http requests/responses.

## Installation
glhf is versioned using go modules. To install it, run

`go get -u github.com/VauntDev/glhf`

## Use

glhf is a simple library that abstracts common http patterns while aiming to not prevent more complex use cases.


### Marshaling

Request and Response marhaling is handled in glfh by utilizing the following HTTP headers

- Content-Type: glhf Content-type to determine how to marshal and unmarshal the request and response.
- Accept : glhf uses the request Accept header to determine what Content-Type should be used by the response.

Currenly glhf only support `Application/json` and `Application/proto`. The default is `Application/json`.


### HTTP Routers

glhf works with any http router that uses `http.handlerFunc` functions.

Standard Library HTTP Mux
```go

	mux := http.NewServeMux()
	mux.HandleFunc("/todo", glhf.Post(h.CreateTodo))
	mux.HandleFunc("/todo/{id}", glhf.Get(h.LookupTodo))

```

Gorilla mux
```go

	mux := mux.NewRouter()
	mux.HandleFunc("/todo", glhf.Post(h.CreateTodo))
	mux.HandleFunc("/todo/{id}", glhf.Get(h.LookupTodo))

```

## Future Work
- glhf router
- cache support ( https://www.rfc-editor.org/rfc/rfc9111.html )

