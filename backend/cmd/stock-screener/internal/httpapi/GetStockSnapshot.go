package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/internal/broker/alpaca"
)

func (api *API) GetStockSnapshot(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	vars := mux.Vars(request)
	symbol := vars["symbol"]
	output, err := api.alpacaClient.GetStockSnapshot(ctx, alpaca.GetStockSnapshotInput{
		Symbol: symbol,
	})
	if err != nil {
		return
	}
	httpx.SendJSONResponse(responseWriter, http.StatusOK, output)
}
