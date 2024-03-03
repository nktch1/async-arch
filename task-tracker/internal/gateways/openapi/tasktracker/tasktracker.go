package tasktracker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/task"
)

type service interface {
	AddTask(context.Context, task.Task) error
	ShuffleTasks(context.Context) error
	CloseTask(context.Context) error
	ListTasks(context.Context) ([]task.Task, error)
}

type server struct {
	service
}

func New(service service) *mux.Router {
	router := mux.NewRouter()

	srv := server{service: service}

	router.HandleFunc("/tasks", srv.ListTasks).Methods(http.MethodGet)
	router.HandleFunc("/tasks/add", srv.AddTask).Methods(http.MethodPost)
	router.HandleFunc("/tasks/shuffle", srv.ShuffleTasks)
	router.HandleFunc("/tasks/close", srv.CloseTask)

	return router
}

func (s server) ListTasks(writer http.ResponseWriter, request *http.Request) {
	tasks, err := s.service.ListTasks(request.Context())
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(writer).Encode(tasks); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (s server) AddTask(writer http.ResponseWriter, request *http.Request) {
	bodyContent, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	var task task.Task

	if err = json.Unmarshal(bodyContent, &task); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	if err = s.service.AddTask(request.Context(), task); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (s server) ShuffleTasks(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "ShuffleTasks")

}

func (s server) CloseTask(writer http.ResponseWriter, request *http.Request) {

	fmt.Fprintln(writer, "CloseTask")
}
