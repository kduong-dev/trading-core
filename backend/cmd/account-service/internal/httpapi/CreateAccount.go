package httpapi

import (
	"net/http"

	"github.com/kduong-dev/goutil/fatal"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/accountstore"
	uuid "github.com/satori/go.uuid"
)

type CreateAccountInput struct {
	AccountName string `json:"account_name"`
}

type CreateAccountOutput struct {
	AccountID   string `json:"account_id"`
	AccountName string `json:"account_name"`
}

func (api *API) CreateAccount(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	// TODO: validate input
	input, err := httpx.DecodeJSONBody[CreateAccountInput](request)
	if err != nil {
		return
	}
	accountID := uuid.NewV4().String()
	err = api.accountStoreCommandHandler.Create(ctx, accountstore.CreateInput{
		AccountID:   accountID,
		AccountName: input.AccountName,
	})
	fatal.OnError(err)
	output := CreateAccountOutput{
		AccountID:   accountID,
		AccountName: input.AccountName,
	}
	httpx.SendJSONResponse(responseWriter, http.StatusCreated, output)
}
