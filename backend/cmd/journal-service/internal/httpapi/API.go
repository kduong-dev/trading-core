package httpapi

import (
	"net/http"

	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/journal-service/internal/entrystore"
)

type API struct {
	entryCommandHandler entrystore.CommandHandler
	entryQueryHandler   entrystore.QueryHandler
}

var merrifiedSentinels = httpx.MerrifiedSentinels{
	{Sentinel: entrystore.ErrEntryNotFound, StatusCode: http.StatusNotFound, UserMessage: "entry not found"},
	{Sentinel: entrystore.ErrEntryForbidden, StatusCode: http.StatusForbidden, UserMessage: "forbidden"},
}
