package httpapi

import (
	"context"
	"net/http"

	"github.com/ansel1/merry"
	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/fatal"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/bot-service/internal/botstore"
)

type UpdateBotInput struct {
	Status string `json:"status"`
}

func (api *API) UpdateBot(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	vars := mux.Vars(request)
	botID := vars["bot_id"]
	body, err := httpx.DecodeJSONBody[UpdateBotInput](request)
	if err != nil {
		return
	}
	status := botstore.BotStatus(body.Status)
	switch status {
	case botstore.BotStatusRunning:
		err = api.ensureAllocationPolicy(ctx, request, botID)
		if err != nil {
			return
		}
	case botstore.BotStatusStopped:
	default:
		err = merry.UserError(`status must be "running" or "stopped"`).WithHTTPCode(http.StatusBadRequest)
		return
	}
	err = api.botStoreCommandHandler.UpdateBotStatus(ctx, botID, status)
	if err = merrifiedSentinels.MerrifyOrFatal(err); err != nil {
		return
	}
	httpx.SendJSONResponse(responseWriter, http.StatusOK, nil)
}

func (api *API) ensureAllocationPolicy(ctx context.Context, request *http.Request, botID string) (err error) {
	bot, err := api.botStoreQueryHandler.Get(ctx, botID)
	if err = merrifiedSentinels.MerrifyOrFatal(err); err != nil {
		return
	}
	ctx = ContextWithAccessTokenFromRequestHeader(ctx, request)
	balance, err := api.accountServiceClient.GetAccountBalance(ctx, bot.AccountID)
	if err != nil {
		err = merrifiedSentinels.Merrify(err)
		return
	}
	if balance.CashBalance <= 0 {
		err = merry.UserError("account has no available cash balance").WithHTTPCode(http.StatusBadRequest)
		return
	}
	bots, err := api.botStoreQueryHandler.List(ctx)
	fatal.OnError(err)
	activeAllocationPercent := 0.0
	for _, botItem := range bots {
		if botItem.ID == botID {
			continue
		}
		if botItem.AccountID != bot.AccountID {
			continue
		}
		if botItem.Status != botstore.BotStatusRunning {
			continue
		}
		activeAllocationPercent += botItem.AllocationPercent
	}
	if activeAllocationPercent+bot.AllocationPercent > MaxActiveAllocationPercent {
		err = merry.UserError("active bot allocation exceeds 80% for this account").WithHTTPCode(http.StatusBadRequest)
		return
	}
	return
}
