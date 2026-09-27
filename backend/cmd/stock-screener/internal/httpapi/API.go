package httpapi

import (
	"github.com/kduong-dev/trading-core/backend/cmd/stock-screener/internal/fetchsentiment"
	"github.com/kduong-dev/trading-core/backend/internal/broker/alpaca"
)

type API struct {
	alpacaClient           alpaca.Client
	fetchSentimentStrategy fetchsentiment.Strategy
}
