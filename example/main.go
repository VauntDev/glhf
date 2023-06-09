package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/unix"
	"google.golang.org/protobuf/proto"

	"github.com/VauntDev/glhf"
	"github.com/VauntDev/glhf/example/pb"
)

type TodoService struct {
	todos map[string]*pb.Todo
}

func (ts *TodoService) Add(t *pb.Todo) error {
	ts.todos[t.Id] = t
	return nil
}

func (ts *TodoService) Get(id string) (*pb.Todo, error) {
	t, ok := ts.todos[id]
	if !ok {
		return nil, fmt.Errorf("no todo")
	}
	return t, nil
}

type Handlers struct {
	service *TodoService
}

func (h *Handlers) StandardLookupTodo(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)

	id, ok := p["id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	todo, err := h.service.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Header.Get("accept") {
	case "application/json":
		b, err := json.Marshal(todo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	case "application/proto":
		b, err := proto.Marshal(todo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func (h *Handlers) StandardCreateTodo(w http.ResponseWriter, r *http.Request) {
	todo := &pb.Todo{}

	switch r.Header.Get("content-type") {
	case "application/json":
		decode := json.NewDecoder(r.Body)
		if err := decode.Decode(todo); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "application/proto":
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := proto.Unmarshal(b, todo); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.service.Add(todo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) GLHFLookupTodo(r *glhf.Request[glhf.EmptyBody], w *glhf.Response[pb.Todo]) {
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

	w.SetBody(todo)
	w.SetStatus(http.StatusOK)
}

func (h *Handlers) GLHFCreateTodo(r *glhf.Request[pb.Todo], w *glhf.Response[glhf.EmptyBody]) {
	if r.Body() == nil {
		w.SetStatus(http.StatusBadRequest)
		return
	}

	if err := h.service.Add(r.Body()); err != nil {
		w.SetStatus(http.StatusInternalServerError)
		return
	}
	w.SetStatus(http.StatusOK)
}

func main() {
	TodoService := &TodoService{
		todos: make(map[string]*pb.Todo),
	}
	h := &Handlers{service: TodoService}

	mux := mux.NewRouter()
	mux.HandleFunc("/standard/todo", h.StandardCreateTodo)
	mux.HandleFunc("/standard/todo/{id}", h.StandardLookupTodo)
	mux.HandleFunc("/glhf/todo", glhf.Post(h.GLHFCreateTodo))
	mux.HandleFunc("/glhf/todo/{id}", glhf.Get(h.GLHFLookupTodo))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		log.Println("starting server")
		if err := server.ListenAndServe(); err != nil {
			return nil
		}
		return nil
	})

	g.Go(func() error {
		sigs := make(chan os.Signal, 1)
		// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
		// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
		signal.Notify(sigs, os.Interrupt, unix.SIGTERM)
		// block until we receive our signal or context is done
		select {
		case <-ctx.Done():
			log.Println("ctx done, shutting down server")
		case <-sigs:
			log.Println("caught sig, shutting down server")
		}
		// Create a deadline to wait for cleanup
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("error in server shutdown: %w", err)
		}
		return nil
	})

	// ----- Client Code ----- //

	// wait for server to start
	time.Sleep(time.Second * 1)

	client := http.DefaultClient

	id := uuid.NewString()
	todo := &pb.Todo{
		Id: id,
		Item: &pb.Item{
			Name:    "Marshaling Example",
			Message: "This todo is used to demo the marshalling functionality of glhf",
		},
	}

	b, err := proto.Marshal(todo)
	if err != nil {
		log.Fatal("failed to marshal proto")
	}

	postRequest, err := http.NewRequest("POST", "http://localhost:8080/glhf/todo", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal("failed to create post request")
	}

	postRequest.Header.Add("Content-Type", "application/proto") // send protobuff

	log.Println("sending post request to create todo")

	postResp, err := client.Do(postRequest)
	if err != nil {
		log.Fatal("failed to do post request", err)
	}

	if postResp.StatusCode != http.StatusOK {
		log.Fatal("post request failed with", postResp.StatusCode)
	}

	getRequest, err := http.NewRequest("GET", "http://localhost:8080/glhf/todo/"+id, nil)
	if err != nil {
		log.Fatal("failed to create get request")
	}

	getRequest.Header.Add("Accept", "application/json") // get json

	log.Println("sending get request to lookup todo")
	getResp, err := client.Do(getRequest)
	if err != nil {
		log.Fatal("failed to do get request", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		log.Fatal("get failed with", getResp.StatusCode)
	}

	body, err := ioutil.ReadAll(getResp.Body)
	if err != nil {
		log.Fatal("failed to read response body", err)
	}

	log.Println(string(body))
}
