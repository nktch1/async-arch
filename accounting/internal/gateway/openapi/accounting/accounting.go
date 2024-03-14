package accounting

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"

	"github.com/nikitych1/awesome-task-exchange-system/accounting/internal/usecase/accounting"
)

type service interface {
	ListTransactions(context.Context) (accounting.ListTransactionsResponse, error)
	ListTransactionsByAccount(context.Context, uuid.UUID) (accounting.ListTransactionsByAccountResponse, error)
}

type server struct {
	service
}

func New(service service) *mux.Router {
	router := mux.NewRouter()

	srv := server{service: service}

	router.HandleFunc("/transactions", srv.ListTransactionsByAccount).Methods(http.MethodGet)
	router.HandleFunc("/all_transactions", srv.ListTransactions).Methods(http.MethodGet)

	return router
}

func (s server) ListTransactions(writer http.ResponseWriter, request *http.Request) {
	transactions, err := s.service.ListTransactions(request.Context())
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(writer).Encode(transactions); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s server) ListTransactionsByAccount(writer http.ResponseWriter, request *http.Request) {
	accountPublicIDString := request.URL.Query().Get("account_public_id")

	accountPublicID, err := uuid.FromString(accountPublicIDString)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	transactions, err := s.service.ListTransactionsByAccount(request.Context(), accountPublicID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(writer).Encode(transactions); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
