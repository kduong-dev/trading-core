package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/internal/authz"
)

func (handler *Handler) GetJob(responseWriter http.ResponseWriter, request *http.Request) {
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
	job, err := handler.jobQueryHandler.Get(ctx, jobID)
	if err != nil {
		err = merrifiedSentinels.Merrify(err)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(job)
}
