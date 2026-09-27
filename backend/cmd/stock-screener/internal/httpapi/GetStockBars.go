package httpapi

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/internal/broker/alpaca"
)

func (api *API) GetStockBars(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	query := request.URL.Query()
	vars := mux.Vars(request)
	timeframe := query.Get("timeframe")
	if timeframe == "" {
		timeframe = "1Day"
	}
	limit := 365
	if v := query.Get("limit"); len(v) > 0 {
		limit, err = strconv.Atoi(v)
		if err != nil {
			return
		}
	}
	feed := query.Get("feed")
	if feed == "" {
		feed = "iex"
	}
	start := query.Get("start")
	end := query.Get("end")
	output, err := api.alpacaClient.GetStockBars(ctx, alpaca.GetStockBarsInput{
		Symbol:    vars["symbol"],
		Timeframe: timeframe,
		Limit:     limit,
		Feed:      feed,
		Start:     start,
		End:       end,
	})
	if err != nil {
		return
	}
	httpx.SendJSONResponse(responseWriter, http.StatusOK, output)
}
