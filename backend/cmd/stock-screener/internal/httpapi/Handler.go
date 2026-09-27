package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/stock-screener/internal/fetchsentiment"
	"github.com/kduong-dev/trading-core/backend/internal/auth"
	"github.com/kduong-dev/trading-core/backend/internal/broker/alpaca"
)

type NewHandlerInput struct {
	AuthMiddleware         *auth.Middleware
	AlpacaClient           alpaca.Client
	FetchSentimentStrategy fetchsentiment.Strategy
}

func NewHandler(input NewHandlerInput) http.Handler {
	api := &API{
		alpacaClient:           input.AlpacaClient,
		fetchSentimentStrategy: input.FetchSentimentStrategy,
	}
	router := mux.NewRouter().StrictSlash(true)
	publicRouter := router.PathPrefix("/stock-screener/v1").Subrouter()
	publicRouter.Use(input.AuthMiddleware.Handle)
	publicRouter.HandleFunc("/most-actives", api.GetActiveStocks).Methods(http.MethodGet).Name("GetActiveStocks")
	publicRouter.HandleFunc("/movers", api.GetTopStockMovers).Methods(http.MethodGet).Name("GetTopStockMovers")
	publicRouter.HandleFunc("/news", api.GetStockNews).Methods(http.MethodGet).Name("GetStockNews")
	publicRouter.HandleFunc("/sentiments/fear-greed", api.GetFearGreedIndex).Methods(http.MethodGet).Name("GetFearGreedIndex")
	publicRouter.HandleFunc("/stocks/{symbol}/bars", api.GetStockBars).Methods(http.MethodGet).Name("GetStockBars")
	publicRouter.HandleFunc("/stocks/{symbol}/snapshot", api.GetStockSnapshot).Methods(http.MethodGet).Name("GetStockSnapshot")
	return httpx.HandlerWithCORS(router)
}
