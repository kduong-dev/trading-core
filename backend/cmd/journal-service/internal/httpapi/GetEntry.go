package httpapi

import (
	"net/http"
	"time"

	"github.com/ansel1/merry"
	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/httpx"
)

func (api *API) GetEntry(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	vars := mux.Vars(request)
	date := vars["date"]
	if _, parseErr := time.Parse(dateLayout, date); parseErr != nil {
		err = merry.UserError("date must be YYYY-MM-DD").WithHTTPCode(http.StatusBadRequest)
		return
	}
	entry, err := api.entryQueryHandler.Get(ctx, date)
	if err != nil {
		err = merrifiedSentinels.Merrify(err)
		return
	}
	httpx.SendJSONResponse(responseWriter, http.StatusOK, entry)
}
