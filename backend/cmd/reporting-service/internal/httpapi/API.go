package httpapi

import (
	"net/http"

	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/storage-service/pkg/storageservice"
	"github.com/kduong-dev/trading-core/backend/cmd/reporting-service/internal/jobstore"
)

type API struct {
	jobCommandHandler jobstore.CommandHandler
	jobQueryHandler   jobstore.QueryHandler
	storageClient     storageservice.Client
	enqueueJob        func(job *jobstore.Job)
}

var merrifiedSentinels = httpx.MerrifiedSentinels{
	{Sentinel: jobstore.ErrJobNotFound, StatusCode: http.StatusNotFound, UserMessage: "job not found"},
	{Sentinel: jobstore.ErrJobForbidden, StatusCode: http.StatusForbidden, UserMessage: "forbidden"},
}
