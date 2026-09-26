package main

import (
	"net/http"

	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/stock-screener/internal/fetchsentiment"
	"github.com/kduong-dev/trading-core/backend/cmd/stock-screener/internal/httpapi"
	"github.com/kduong-dev/trading-core/backend/internal/auth"
	"github.com/kduong-dev/trading-core/backend/internal/broker/alpaca"
)

func main() {
	router := httpapi.NewRouter(httpapi.NewRouterInput{
		AlpacaClient:           alpaca.ClientFromEnv(),
		AuthMiddleware:         auth.MiddlewareFromEnv(auth.AudienceStockScreenerService),
		FetchSentimentStrategy: fetchsentiment.StrategyFromEnv(),
	})
	http.ListenAndServe(":8080", httpx.HandlerWithCORS(router))
}
