# GLHF

Generic Light-weight Handler Framework

## Experimental

This is a **experimental** library. The goal is to evaluate leveraging generics to reduce duplicate code around deserializing and serializing http requests/responses.

## Installation

GLHF is versioned using go modules. To install it, run

`go get -u github.com/VauntDev/glhf`

## Use

GLHF is a simple library that abstracts common http patterns while aiming to not prevent more complex use cases.

### Marshaling

Request and Response marshaling is handled in glfh by utilizing the following HTTP headers

- Content-Type: GLHF Content-type to determine how to marshal and unmarshal the request and response.
- Accept : GLHF uses the request Accept header to determine what Content-Type should be used by the response.

Currenly glhf only support `Application/json` and `Application/proto`. The default is `Application/json`.

### HTTP Routers

GLHF works with any http router that uses `http.handlerFunc` functions.

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

- GLHF router
- cache support ( https://www.rfc-editor.org/rfc/rfc9111.html )

## Examples

A sample application can be found in the [example](./example/main.go) directory.

The following is an example GET handler. The functions expects an empty body and Todo body in response.

```go

func (h *Handlers) LookupTodo(r *glhf.Request[glhf.EmptyBody], w *glhf.Response[pb.Todo]) {
    p := mux.Vars(r.HTTPRequest())

    id, ok := p["id"]
    if !ok {
        w.SetStatus(http.StatusInternalServerError)
        return
    }

    todo, err := h.service.Get(id)
    if err != nil {
        w.SetStatus(http.StatusNotFound)
        return
    }

    w.Body = todo
    log.Println("external handler", w.Body)
    w.SetStatus(http.StatusOK)
    return

}
```

The following is an example POST handler. The function expects a Todo body and Todo response.

```go

func (h *Handlers) CreateTodo(r *glhf.Request[pb.Todo], w *glhf.Response[glhf.EmptyBody]) {
    t, err := r.Body()
    if err != nil {
        w.SetStatus(http.StatusBadRequest)
        return
    }

    if err := h.service.Add(t); err != nil {
        w.SetStatus(http.StatusInternalServerError)
        return
    }
    w.SetStatus(http.StatusOK)
    return
}
```
