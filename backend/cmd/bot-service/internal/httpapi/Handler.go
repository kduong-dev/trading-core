package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/eventsource"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/pkg/accountservice"
	"github.com/kduong-dev/trading-core/backend/cmd/bot-service/internal/botstore"
	"github.com/kduong-dev/trading-core/backend/cmd/bot-service/internal/symbolvalidator"
	"github.com/kduong-dev/trading-core/backend/internal/auth"
)

type NewHandlerInput struct {
	AuthMiddleware         *auth.Middleware
	AccountServiceClient   accountservice.Client
	SymbolValidator        symbolvalidator.SymbolValidator
	BotStoreCommandHandler botstore.CommandHandler
	BotStoreQueryHandler   botstore.QueryHandler
	BotEventLogFactory     eventsource.LogFactory
	BotChannelFunc         func(botID string) string
}

func NewHandler(input NewHandlerInput) http.Handler {
	symbolValidator := input.SymbolValidator
	if symbolValidator == nil {
		symbolValidator = symbolvalidator.NoopSymbolValidator{}
	}
	api := &API{
		accountServiceClient:   input.AccountServiceClient,
		symbolValidator:        symbolValidator,
		botStoreCommandHandler: input.BotStoreCommandHandler,
		botStoreQueryHandler:   input.BotStoreQueryHandler,
		botEventLogFactory:     input.BotEventLogFactory,
		botChannelFunc:         input.BotChannelFunc,
	}
	router := mux.NewRouter().StrictSlash(true)
	publicRouter := router.PathPrefix("/bots/v1").Subrouter()
	publicRouter.Use(input.AuthMiddleware.Handle)
	publicRouter.HandleFunc("/bots", api.CreateBot).Methods(http.MethodPost).Name("CreateBot")
	publicRouter.HandleFunc("/bots", api.ListBots).Methods(http.MethodGet).Name("ListBots")
	publicRouter.HandleFunc("/bots/{bot_id}", api.GetBot).Methods(http.MethodGet).Name("GetBot")
	publicRouter.HandleFunc("/bots/{bot_id}/stream", api.StreamBotEvents).Methods(http.MethodGet).Name("StreamBotEvents")
	publicRouter.HandleFunc("/bots/{bot_id}", api.UpdateBot).Methods(http.MethodPatch).Name("UpdateBot")
	publicRouter.HandleFunc("/bots/{bot_id}", api.DeleteBot).Methods(http.MethodDelete).Name("DeleteBot")
	return httpx.HandlerWithCORS(router)
}
