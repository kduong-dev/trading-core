package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/journal-service/internal/entrystore"
	"github.com/kduong-dev/trading-core/backend/internal/auth"
)

type NewHandlerInput struct {
	AuthMiddleware      *auth.Middleware
	EntryCommandHandler entrystore.CommandHandler
	EntryQueryHandler   entrystore.QueryHandler
}

func NewHandler(input NewHandlerInput) http.Handler {
	api := &API{
		entryCommandHandler: input.EntryCommandHandler,
		entryQueryHandler:   input.EntryQueryHandler,
	}
	router := mux.NewRouter().StrictSlash(true)
	publicRouter := router.PathPrefix("/journal/v1").Subrouter()
	publicRouter.Use(input.AuthMiddleware.Handle)
	publicRouter.HandleFunc("/entries", api.ListEntries).Methods(http.MethodGet).Name("ListEntries")
	publicRouter.HandleFunc("/entries/{date}", api.GetEntry).Methods(http.MethodGet).Name("GetEntry")
	publicRouter.HandleFunc("/entries/{date}", api.UpsertEntry).Methods(http.MethodPut).Name("UpsertEntry")
	publicRouter.HandleFunc("/entries/{date}", api.DeleteEntry).Methods(http.MethodDelete).Name("DeleteEntry")
	return httpx.HandlerWithCORS(router)
}
