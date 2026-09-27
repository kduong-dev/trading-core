package httpapi

import (
	"io"
	"net/http"
	"strings"

	"github.com/ansel1/merry"
	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/storage-service/pkg/storageservice"
	"github.com/kduong-dev/trading-core/backend/cmd/reporting-service/internal/jobstore"
	"github.com/kduong-dev/trading-core/backend/internal/authz"
)

func (api *API) DownloadJob(responseWriter http.ResponseWriter, request *http.Request) {
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
	if err = merrifiedSentinels.MerrifyOrFatal(err); err != nil {
		return
	}
	if job.Status != jobstore.JobStatusCompleted {
		err = merry.UserError("job is not yet available for download").WithHTTPCode(http.StatusConflict)
		return
	}
	fileID := extractFileID(job.DownloadURL)
	if fileID == "" {
		err = merry.UserError("job has no file attached").WithHTTPCode(http.StatusNotFound)
		return
	}
	// jobQueryHandler.Get has already enforced that the caller owns this job;
	// storage-service only knows about trading-core as a whole.
	download, err := api.storageClient.DownloadFile(ctx, storageservice.DownloadFileInput{FileID: fileID})
	if err != nil {
		return
	}
	defer download.Body.Close()
	responseWriter.Header().Set("Content-Type", download.ContentType)
	if download.ContentDisposition != "" {
		responseWriter.Header().Set("Content-Disposition", download.ContentDisposition)
	}
	// The response has started, so a copy error can only mean the client or upstream has gone.
	_, _ = io.Copy(responseWriter, download.Body)
}

// extractFileID parses the file ID from a storage-service path of the form
// /storage/v1/files/<id>.
func extractFileID(downloadURL string) string {
	const prefix = "/storage/v1/files/"
	if !strings.HasPrefix(downloadURL, prefix) {
		return ""
	}
	return strings.TrimPrefix(downloadURL, prefix)
}
