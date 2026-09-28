package httpapi

import (
	"net/http"
	"time"

	"github.com/ansel1/merry"
	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/accountstore"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/pnlaggregator"
	"github.com/kduong-dev/trading-core/backend/internal/broker"
)

const dailyPnLDateLayout = "2006-01-02"
const dailyPnLMaxRangeDays = 366

// dailyPnLMatchingLookbackDays widens the broker fetch backwards from the
// requested `from` date so that closes within the requested window can be
// FIFO-matched against opens that happened earlier. The extra rows are
// discarded after matching and never appear in the response.
const dailyPnLMatchingLookbackDays = 365

type GetDailyPnLResponse struct {
	Currency string                   `json:"currency"`
	Days     []pnlaggregator.DailyPnL `json:"days"`
	Summary  pnlaggregator.Summary    `json:"summary"`
}

func (api *API) GetDailyPnL(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			merrifiedSentinels.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	vars := mux.Vars(request)
	accountID := vars["account_id"]

	from := request.URL.Query().Get("from")
	to := request.URL.Query().Get("to")
	if from == "" || to == "" {
		err = merry.UserError("from and to are required (YYYY-MM-DD)").WithHTTPCode(http.StatusBadRequest)
		return
	}
	fromDate, parseErr := time.Parse(dailyPnLDateLayout, from)
	if parseErr != nil {
		err = merry.UserError("from must be YYYY-MM-DD").WithHTTPCode(http.StatusBadRequest)
		return
	}
	toDate, parseErr := time.Parse(dailyPnLDateLayout, to)
	if parseErr != nil {
		err = merry.UserError("to must be YYYY-MM-DD").WithHTTPCode(http.StatusBadRequest)
		return
	}
	if toDate.Before(fromDate) {
		err = merry.UserError("to must be on or after from").WithHTTPCode(http.StatusBadRequest)
		return
	}
	if toDate.Sub(fromDate).Hours()/24 > dailyPnLMaxRangeDays {
		err = merry.UserError("date range must be at most 366 days").WithHTTPCode(http.StatusBadRequest)
		return
	}

	account, err := api.accountStoreQueryHandler.Get(ctx, accountstore.GetInput{
		AccountID: accountID,
	})
	if err != nil {
		return
	}
	err = checkBrokerLinked(account)
	if err != nil {
		return
	}
	accountClient := api.brokerAccountClientFactory.Get(ctx, account.BrokerAccount)
	matchingFrom := fromDate.AddDate(0, 0, -dailyPnLMatchingLookbackDays).Format(dailyPnLDateLayout)
	transactionsOutput, err := accountClient.GetTransactions(ctx, broker.GetTransactionsInput{
		From: matchingFrom,
		To:   to,
	})
	if err != nil {
		return
	}
	pnlaggregator.MatchRealizedPnL(transactionsOutput.Transactions)
	withinRequestedWindow := pnlaggregator.FilterByDateRange(transactionsOutput.Transactions, from, to)
	aggregated := pnlaggregator.Aggregate(withinRequestedWindow)
	summary := pnlaggregator.Summarize(withinRequestedWindow)

	balance, err := accountClient.GetBalance(ctx)
	currency := "USD"
	if err == nil && balance != nil && balance.Currency != "" {
		currency = balance.Currency
	}
	err = nil

	response := GetDailyPnLResponse{
		Currency: currency,
		Days:     aggregated.Days,
		Summary:  summary,
	}
	httpx.SendJSONResponse(responseWriter, http.StatusOK, response)
}
