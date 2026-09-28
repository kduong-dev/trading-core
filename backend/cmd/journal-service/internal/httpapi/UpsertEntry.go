package httpapi

import (
	"net/http"
	"time"

	"github.com/ansel1/merry"
	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/journal-service/internal/entrystore"
)

const dateLayout = "2006-01-02"

type UpsertEntryInput struct {
	Notes             string   `json:"notes,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	Mood              string   `json:"mood,omitempty"`
	DisciplineScore   int      `json:"discipline_score,omitempty"`
	ScreenshotFileIDs []string `json:"screenshot_file_ids,omitempty"`
}

func (input *UpsertEntryInput) Validate() error {
	if input.DisciplineScore < 0 || input.DisciplineScore > 10 {
		return merry.UserError("discipline_score must be between 0 and 10").WithHTTPCode(http.StatusBadRequest)
	}
	return nil
}

func (api *API) UpsertEntry(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			merrifiedSentinels.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	vars := mux.Vars(request)
	date := vars["date"]
	if _, parseErr := time.Parse(dateLayout, date); parseErr != nil {
		err = merry.UserError("date must be YYYY-MM-DD").WithHTTPCode(http.StatusBadRequest)
		return
	}
	input, err := httpx.DecodeJSONBody[UpsertEntryInput](request)
	if err != nil {
		return
	}
	err = input.Validate()
	if err != nil {
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	entry := &entrystore.Entry{
		Date:              date,
		Notes:             input.Notes,
		Tags:              input.Tags,
		Mood:              input.Mood,
		DisciplineScore:   input.DisciplineScore,
		ScreenshotFileIDs: input.ScreenshotFileIDs,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	err = api.entryCommandHandler.UpsertEntry(ctx, entry)
	if err != nil {
		return
	}
	httpx.SendJSONResponse(responseWriter, http.StatusAccepted, entry)
}
