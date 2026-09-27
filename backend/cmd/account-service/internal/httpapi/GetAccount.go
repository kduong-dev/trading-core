package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/accountstore"
)

func (api *API) GetAccount(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	vars := mux.Vars(request)
	accountID := vars["account_id"]
	account, err := api.accountStoreQueryHandler.Get(ctx, accountstore.GetInput{
		AccountID: accountID,
	})
	if err != nil {
		err = merrifiedSentinels.Merrify(err)
		return
	}
	httpx.SendJSONResponse(responseWriter, http.StatusOK, account)
}
