package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/internal/authz"
)

func (api *API) GetJob(responseWriter http.ResponseWriter, request *http.Request) {
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
	vars := mux.Vars(request)
	jobID := vars["job_id"]
	job, err := api.jobQueryHandler.Get(ctx, jobID)
	if err != nil {
		err = merrifiedSentinels.Merrify(err)
		return
	}
	httpx.SendJSONResponse(responseWriter, http.StatusOK, job)
}
