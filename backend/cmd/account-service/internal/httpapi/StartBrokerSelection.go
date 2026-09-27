package httpapi

import (
	"net/http"
	"time"

	"github.com/ansel1/merry"
	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/accountstore"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/oauthstatestore"
	"github.com/kduong-dev/trading-core/backend/internal/broker"
	"github.com/kduong-dev/trading-core/backend/internal/contextx"
)

type StartBrokerSelectionInput struct {
	Broker broker.AccountType `json:"broker"`
}

type StartBrokerSelectionOutput struct {
	AuthorizationURL string `json:"authorization_url"`
}

func (api *API) StartBrokerSelection(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	userID := contextx.GetUserID(ctx)
	vars := mux.Vars(request)
	accountID := vars["account_id"]
	input, err := httpx.DecodeJSONBody[StartBrokerSelectionInput](request)
	if err != nil {
		return
	}
	if input.Broker == "" {
		err = merry.UserError("broker is required").WithHTTPCode(http.StatusBadRequest)
		return
	}
	_, err = api.accountStoreQueryHandler.Get(ctx, accountstore.GetInput{
		AccountID: accountID,
	})
	if err != nil {
		err = merrifiedSentinels.Merrify(err)
		return
	}
	stateToken, err := GenerateStateToken()
	if err != nil {
		return
	}
	authorizationClient, err := api.brokerOnBoardingClientFactory.GetAuthorizationClient(input.Broker)
	if err != nil {
		err = merry.Wrap(err).WithHTTPCode(http.StatusBadRequest).WithUserMessage("unsupported broker")
		return
	}
	api.oauthStateStore.Put(stateToken, oauthstatestore.Entry{
		AccountID: accountID,
		UserID:    userID,
		Broker:    input.Broker,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	})
	httpx.SendJSONResponse(responseWriter, http.StatusOK, StartBrokerSelectionOutput{
		AuthorizationURL: authorizationClient.BuildAuthorizationURL(stateToken),
	})
}
