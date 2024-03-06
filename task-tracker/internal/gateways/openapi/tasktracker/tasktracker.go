package tasktracker

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"

	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/task"
)

type service interface {
	AddTask(context.Context, task.Task) error
	ShuffleTasks(context.Context) error
	CloseTask(context.Context, uuid.UUID) error
	ListTasks(context.Context, uuid.UUID) ([]task.Task, error)
}

type server struct {
	service
}

func New(service service) *mux.Router {
	router := mux.NewRouter()

	srv := server{service: service}

	router.HandleFunc("/tasks", srv.ListTasks).Methods(http.MethodGet)
	router.HandleFunc("/tasks/add", srv.AddTask).Methods(http.MethodPost)
	router.HandleFunc("/tasks/shuffle", srv.ShuffleTasks).Methods(http.MethodPatch)
	router.HandleFunc("/tasks/close", srv.CloseTask).Methods(http.MethodPatch)

	return router
}

func (s server) ListTasks(writer http.ResponseWriter, request *http.Request) {
	accountPublicIDString := request.URL.Query().Get("account_public_id")

	accountPublicID, err := uuid.FromString(accountPublicIDString)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	tasks, err := s.service.ListTasks(request.Context(), accountPublicID)
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
	if err := s.service.ShuffleTasks(request.Context()); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (s server) CloseTask(writer http.ResponseWriter, request *http.Request) {
	taskPublicIDString := request.URL.Query().Get("public_id")

	taskPublicID, err := uuid.FromString(taskPublicIDString)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	if err = s.service.CloseTask(request.Context(), taskPublicID); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
