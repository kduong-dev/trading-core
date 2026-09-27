package httpapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/kduong-dev/goutil/eventsource"
	"github.com/kduong-dev/goutil/fatal"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/pkg/accountservice"
	"github.com/kduong-dev/trading-core/backend/cmd/bot-service/internal/botstore"
	"github.com/kduong-dev/trading-core/backend/cmd/bot-service/internal/symbolvalidator"
	"github.com/kduong-dev/trading-core/backend/internal/contextx"
)

const MaxActiveAllocationPercent = 80.0

type API struct {
	accountServiceClient   accountservice.Client
	symbolValidator        symbolvalidator.SymbolValidator
	botStoreCommandHandler botstore.CommandHandler
	botStoreQueryHandler   botstore.QueryHandler
	botEventLogFactory     eventsource.LogFactory
	botChannelFunc         func(botID string) string
}

func ContextWithAccessTokenFromRequestHeader(ctx context.Context, request *http.Request) context.Context {
	authorization := request.Header.Get("Authorization")
	parts := strings.SplitN(authorization, " ", 2)
	fatal.Unless(len(parts) == 2, "invalid authorization header format")
	return contextx.WithAccessToken(ctx, parts[1])
}

var merrifiedSentinels = httpx.MerrifiedSentinels{
	{Sentinel: botstore.ErrBotNotFound, StatusCode: http.StatusNotFound, UserMessage: "bot not found"},
	{Sentinel: botstore.ErrBotForbidden, StatusCode: http.StatusForbidden, UserMessage: "forbidden"},
	{Sentinel: accountservice.ErrAccountNotFound, StatusCode: http.StatusNotFound, UserMessage: "account not found"},
	{Sentinel: accountservice.ErrAccountForbidden, StatusCode: http.StatusForbidden, UserMessage: "forbidden"},
	{Sentinel: accountservice.ErrServerError, StatusCode: http.StatusInternalServerError, UserMessage: "account service error"},
}
