package httpapi

import (
	"net/http"

	"github.com/kduong-dev/goutil/httpx"
)

func (api *API) ListBots(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			merrifiedSentinels.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	bots, err := api.botStoreQueryHandler.List(ctx)
	if err != nil {
		return
	}
	httpx.SendJSONResponse(responseWriter, http.StatusOK, bots)
}
