package tasktracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/task"
)

type service interface {
	AddTask(task.Task) error
	ShuffleTasks() error
	CloseTask() error
	ListTasks() ([]task.Task, error)
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
	tasks, err := s.service.ListTasks()
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

	if err = s.service.AddTask(task); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (s server) ShuffleTasks(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "ShuffleTasks")

}

func (s server) CloseTask(writer http.ResponseWriter, request *http.Request) {

	fmt.Fprintln(writer, "CloseTask")
}
