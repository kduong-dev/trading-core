package httpapi

import (
	"net/http"
	"strconv"

	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/internal/broker/alpaca"
)

func (api *API) GetTopStockMovers(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	query := request.URL.Query()
	limit := 10
	if v := query.Get("limit"); len(v) > 0 {
		limit, err = strconv.Atoi(v)
		if err != nil {
			return
		}
	}
	output, err := api.alpacaClient.GetTopStockMovers(ctx, alpaca.GetTopStockMoversInput{
		Limit: limit,
	})
	if err != nil {
		return
	}
	httpx.SendJSONResponse(responseWriter, http.StatusOK, output)
}
