package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ansel1/merry"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/reporting-service/internal/jobstore"
	"github.com/kduong-dev/trading-core/backend/internal/authz"
)

const defaultPageSize = 10
const maxPageSize = 100

func (handler *Handler) ListJobs(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	if err = authz.RequireScope(ctx, authz.ScopeJobsRead); err != nil {
		return
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
		err = merry.UserError("page_size must be between 1 and 100").WithHTTPCode(http.StatusBadRequest)
		return
	}
	if page < 0 {
		err = merry.UserError("page must be >= 0").WithHTTPCode(http.StatusBadRequest)
		return
	}

	result, err := handler.jobQueryHandler.List(ctx, jobstore.ListInput{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(result)
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
