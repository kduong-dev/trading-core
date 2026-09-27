package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/accountstore"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/oauthstatestore"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/pendingselectionstore"
	"github.com/kduong-dev/trading-core/backend/internal/auth"
	"github.com/kduong-dev/trading-core/backend/internal/broker"
)

type NewHandlerInput struct {
	AuthMiddleware                *auth.Middleware
	OAuthStateStore               oauthstatestore.Store
	PendingSelectionStore         pendingselectionstore.Store
	AccountStoreCommandHandler    accountstore.CommandHandler
	AccountStoreQueryHandler      accountstore.QueryHandler
	BrokerAccountClientFactory    broker.AccountClientFactory
	BrokerOnBoardingClientFactory broker.OnBoardingClientFactory
	BackendRedirectURI            string
	FrontendBaseURL               string
}

func NewHandler(input NewHandlerInput) http.Handler {
	api := &API{
		oauthStateStore:               input.OAuthStateStore,
		pendingSelectionStore:         input.PendingSelectionStore,
		accountStoreCommandHandler:    input.AccountStoreCommandHandler,
		accountStoreQueryHandler:      input.AccountStoreQueryHandler,
		brokerAccountClientFactory:    input.BrokerAccountClientFactory,
		brokerOnBoardingClientFactory: input.BrokerOnBoardingClientFactory,
		backendRedirectURI:            input.BackendRedirectURI,
		frontendBaseURL:               input.FrontendBaseURL,
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/accounts/v1/authorization_callback", api.HandleAuthorizationCallback).Methods(http.MethodGet).Name("HandleAuthorizationCallback")
	publicRouter := router.PathPrefix("/accounts/v1").Subrouter()
	publicRouter.Use(input.AuthMiddleware.Handle)
	publicRouter.HandleFunc("/accounts", api.CreateAccount).Methods(http.MethodPost).Name("CreateAccount")
	publicRouter.HandleFunc("/accounts", api.ListAccounts).Methods(http.MethodGet).Name("ListAccounts")
	publicRouter.HandleFunc("/accounts/{account_id}", api.GetAccount).Methods(http.MethodGet).Name("GetAccount")
	publicRouter.HandleFunc("/accounts/{account_id}/balances", api.GetAccountBalance).Methods(http.MethodGet).Name("GetAccountBalance")
	publicRouter.HandleFunc("/accounts/{account_id}/pnl/daily", api.GetDailyPnL).Methods(http.MethodGet).Name("GetDailyPnL")
	publicRouter.HandleFunc("/accounts/{account_id}/brokers", api.StartBrokerSelection).Methods(http.MethodPost).Name("StartBrokerSelection")
	publicRouter.HandleFunc("/accounts/{account_id}/brokers", api.GetPendingBrokerSelection).Methods(http.MethodGet).Name("GetPendingBrokerSelection")
	publicRouter.HandleFunc("/accounts/{account_id}/brokers", api.CompleteBrokerSelection).Methods(http.MethodPut).Name("CompleteBrokerSelection")
	return httpx.HandlerWithCORS(router)
}
