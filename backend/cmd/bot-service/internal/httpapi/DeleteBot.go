package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/httpx"
)

func (api *API) DeleteBot(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	vars := mux.Vars(request)
	botID := vars["bot_id"]
	err = api.botStoreCommandHandler.Delete(ctx, botID)
	if err != nil {
		err = merrifiedSentinels.Merrify(err)
		return
	}
	httpx.SendJSONResponse(responseWriter, http.StatusOK, nil)
}
