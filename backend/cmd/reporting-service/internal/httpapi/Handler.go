package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/storage-service/pkg/storageservice"
	"github.com/kduong-dev/trading-core/backend/cmd/reporting-service/internal/jobstore"
	"github.com/kduong-dev/trading-core/backend/internal/auth"
)

type NewHandlerInput struct {
	AuthMiddleware    *auth.Middleware
	JobCommandHandler jobstore.CommandHandler
	JobQueryHandler   jobstore.QueryHandler
	StorageClient     storageservice.Client
	EnqueueJob        func(job *jobstore.Job)
}

func NewHandler(input NewHandlerInput) http.Handler {
	api := &API{
		jobCommandHandler: input.JobCommandHandler,
		jobQueryHandler:   input.JobQueryHandler,
		storageClient:     input.StorageClient,
		enqueueJob:        input.EnqueueJob,
	}
	router := mux.NewRouter().StrictSlash(true)
	publicRouter := router.PathPrefix("/reports/v1").Subrouter()
	publicRouter.Use(input.AuthMiddleware.Handle)
	publicRouter.HandleFunc("/jobs", api.CreateJob).Methods(http.MethodPost).Name("CreateJob")
	publicRouter.HandleFunc("/jobs", api.ListJobs).Methods(http.MethodGet).Name("ListJobs")
	publicRouter.HandleFunc("/jobs/{job_id}", api.GetJob).Methods(http.MethodGet).Name("GetJob")
	publicRouter.HandleFunc("/jobs/{job_id}/download", api.DownloadJob).Methods(http.MethodGet).Name("DownloadJob")
	return httpx.HandlerWithCORS(router)
}
