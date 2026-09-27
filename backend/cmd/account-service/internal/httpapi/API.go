package httpapi

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/ansel1/merry"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/accountstore"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/oauthstatestore"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/pendingselectionstore"
	"github.com/kduong-dev/trading-core/backend/internal/broker"
)

type API struct {
	oauthStateStore               oauthstatestore.Store
	pendingSelectionStore         pendingselectionstore.Store
	accountStoreCommandHandler    accountstore.CommandHandler
	accountStoreQueryHandler      accountstore.QueryHandler
	brokerAccountClientFactory    broker.AccountClientFactory
	brokerOnBoardingClientFactory broker.OnBoardingClientFactory
	backendRedirectURI            string
	frontendBaseURL               string
}

func GenerateStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func checkBrokerLinked(account *accountstore.Account) error {
	if !account.BrokerLinked {
		return merry.UserError("account is not linked to a broker").WithHTTPCode(http.StatusBadRequest)
	}
	return nil
}

var merrifiedSentinels = httpx.MerrifiedSentinels{
	{Sentinel: accountstore.ErrAccountNotFound, StatusCode: http.StatusNotFound, UserMessage: "account not found"},
	{Sentinel: accountstore.ErrAccountForbidden, StatusCode: http.StatusForbidden, UserMessage: "forbidden"},
	{Sentinel: accountstore.ErrBrokerAccountAlreadyLinked, StatusCode: http.StatusConflict, UserMessage: "broker already linked"},
}
