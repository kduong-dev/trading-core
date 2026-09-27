package httpapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ansel1/merry"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/journal-service/internal/entrystore"
)

const defaultPageSize = 31
const maxPageSize = 366

func (api *API) ListEntries(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()

	from := request.URL.Query().Get("from")
	to := request.URL.Query().Get("to")
	if from != "" {
		if _, parseErr := time.Parse(dateLayout, from); parseErr != nil {
			err = merry.UserError("from must be YYYY-MM-DD").WithHTTPCode(http.StatusBadRequest)
			return
		}
	}
	if to != "" {
		if _, parseErr := time.Parse(dateLayout, to); parseErr != nil {
			err = merry.UserError("to must be YYYY-MM-DD").WithHTTPCode(http.StatusBadRequest)
			return
		}
	}

	page, err := parseQueryInt(request, "page", 0)
	if err != nil {
		return
	}
	pageSize, err := parseQueryInt(request, "page_size", defaultPageSize)
	if err != nil {
		return
	}
	if pageSize < 1 || pageSize > maxPageSize {
		err = merry.UserError("page_size must be between 1 and 366").WithHTTPCode(http.StatusBadRequest)
		return
	}
	if page < 0 {
		err = merry.UserError("page must be >= 0").WithHTTPCode(http.StatusBadRequest)
		return
	}

	result, err := api.entryQueryHandler.List(ctx, entrystore.ListInput{
		From:     from,
		To:       to,
		Page:     page,
		PageSize: pageSize,
	})
	if err = merrifiedSentinels.MerrifyOrFatal(err); err != nil {
		return
	}
	httpx.SendJSONResponse(responseWriter, http.StatusOK, result)
}

func parseQueryInt(request *http.Request, key string, defaultValue int) (int, error) {
	raw := request.URL.Query().Get(key)
	if raw == "" {
		return defaultValue, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, merry.UserErrorf("%s must be an integer", key).WithHTTPCode(http.StatusBadRequest)
	}
	return value, nil
}
