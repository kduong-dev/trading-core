package httpapi

import (
	"net/http"

	"github.com/kduong-dev/goutil/httpx"
)

func (api *API) GetFearGreedIndex(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	output, err := api.fetchSentimentStrategy.GetFearGreedIndex(ctx)
	if err != nil {
		return
	}
	httpx.SendJSONResponse(responseWriter, http.StatusOK, output)
}
