package httpapi

import (
	"net/http"
	"slices"

	"github.com/ansel1/merry"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/accountstore"
	"github.com/kduong-dev/trading-core/backend/internal/broker"
	"github.com/kduong-dev/trading-core/backend/internal/contextx"
)

type CompleteBrokerSelectionInput struct {
	PendingToken    string `json:"pending_token"`
	BrokerAccountID string `json:"broker_account_id"`
}

type CompleteBrokerSelectionOutput struct {
	AccountID     string         `json:"account_id"`
	BrokerAccount broker.Account `json:"broker_account"`
}

func (api *API) CompleteBrokerSelection(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	userID := contextx.GetUserID(ctx)
	input, err := httpx.DecodeJSONBody[CompleteBrokerSelectionInput](request)
	if err != nil {
		return
	}
	if input.PendingToken == "" || input.BrokerAccountID == "" {
		err = merry.UserError("pending_token and broker_account_id are required").WithHTTPCode(http.StatusBadRequest)
		return
	}
	entry, ok := api.pendingSelectionStore.Get(input.PendingToken)
	if !ok {
		err = merry.UserError("pending broker selection not found").WithHTTPCode(http.StatusNotFound)
		return
	}
	if entry.UserID != userID {
		err = merry.UserError("forbidden").WithHTTPCode(http.StatusForbidden)
		return
	}
	isValidBrokerAccount := slices.Contains(entry.BrokerAccounts, input.BrokerAccountID)
	if !isValidBrokerAccount {
		err = merry.UserError("broker account is not available for this selection").WithHTTPCode(http.StatusBadRequest)
		return
	}
	brokerAccount := &broker.Account{
		Type: entry.Broker,
		ID:   input.BrokerAccountID,
	}
	ctx = contextx.WithUserID(ctx, entry.UserID)
	err = api.accountStoreCommandHandler.LinkBrokerAccount(ctx, accountstore.LinkBrokerAccountInput{
		AccountID:     entry.AccountID,
		BrokerAccount: brokerAccount,
	})
	if err != nil {
		err = merrifiedSentinels.Merrify(err)
		return
	}
	api.pendingSelectionStore.Delete(input.PendingToken)
	httpx.SendJSONResponse(responseWriter, http.StatusOK, CompleteBrokerSelectionOutput{
		AccountID:     entry.AccountID,
		BrokerAccount: *brokerAccount,
	})
}
