package httpapi

import (
	"net/http"

	"github.com/kduong-dev/goutil/httpx"
)

func (api *API) ListAccounts(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	accounts, err := api.accountStoreQueryHandler.List(ctx)
	if err != nil {
		return
	}
	httpx.SendJSONResponse(responseWriter, http.StatusOK, accounts)
}
