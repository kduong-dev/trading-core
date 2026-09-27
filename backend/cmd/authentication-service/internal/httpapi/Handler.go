package httpapi

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/authentication-service/internal/userstore"
)

type NewHandlerInput struct {
	UserStore   userstore.Store
	TokenSecret []byte
	ExpiryTTL   time.Duration
}

func NewHandler(input NewHandlerInput) http.Handler {
	api := &API{
		userStore:   input.UserStore,
		tokenSecret: input.TokenSecret,
		expiryTTL:   input.ExpiryTTL,
	}
	router := mux.NewRouter().StrictSlash(true)
	publicRouter := router.PathPrefix("/auth/v1").Subrouter()
	publicRouter.HandleFunc("/users", api.CreateUser).Methods(http.MethodPost).Name("CreateUser")
	publicRouter.HandleFunc("/sessions", api.CreateSession).Methods(http.MethodPost).Name("CreateSession")
	publicRouter.HandleFunc("/sessions/refresh", api.RefreshSession).Methods(http.MethodPost).Name("RefreshSession")
	return httpx.HandlerWithCORS(router)
}
