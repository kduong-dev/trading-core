package httpapi

import (
	"net/http"
	"strconv"

	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/internal/broker/alpaca"
)

func (api *API) GetActiveStocks(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	query := request.URL.Query()
	rankBy := alpaca.RankBy(query.Get("rank_by"))
	if len(rankBy) == 0 {
		rankBy = alpaca.RankByTradeCount
	}
	limit := 10
	if v := query.Get("limit"); len(v) > 0 {
		limit, err = strconv.Atoi(v)
		if err != nil {
			return
		}
	}
	output, err := api.alpacaClient.GetActiveStocks(ctx, alpaca.GetActiveStocksInput{
		RankBy: rankBy,
		Limit:  limit,
	})
	if err != nil {
		return
	}
	httpx.SendJSONResponse(responseWriter, http.StatusOK, output)
}
