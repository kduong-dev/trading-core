package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/fatal"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/internal/broker/alpaca"
)

func (handler *Handler) GetStockSnapshot(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	vars := mux.Vars(request)
	symbol := vars["symbol"]
	output, err := handler.alpacaClient.GetStockSnapshot(ctx, alpaca.GetStockSnapshotInput{
		Symbol: symbol,
	})
	if err != nil {
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(responseWriter).Encode(output)
	fatal.OnErrorUnlessDone(ctx, err)
}
