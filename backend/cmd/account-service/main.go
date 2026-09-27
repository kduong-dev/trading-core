package main

import (
	"net/http"
	"net/url"

	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/accountstore"

	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/httpapi"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/oauthstatestore"
	"github.com/kduong-dev/trading-core/backend/cmd/account-service/internal/pendingselectionstore"

	"github.com/kduong-dev/goutil/config"
	"github.com/kduong-dev/goutil/eventsource"
	"github.com/kduong-dev/goutil/fatal"
	"github.com/kduong-dev/trading-core/backend/internal/auth"
	"github.com/kduong-dev/trading-core/backend/internal/broker"
	"github.com/kduong-dev/trading-core/backend/internal/broker/tastytrade"
)

func main() {
	logFactory, err := eventsource.LogFactoryFromEnv("ACCOUNT_EVENT_LOG", "INMEMORY")
	fatal.OnError(err)
	log, err := logFactory.Create("account:events")
	fatal.OnError(err)
	authorizationRedirectURI := url.URL{
		Scheme: config.EnvStringOrFatal("TRADING_API_SCHEME"),
		Host:   config.EnvStringOrFatal("TRADING_API_HOST"),
		Path:   "/accounts/v1/authorization_callback",
	}
	credentialsByType := auth.CredentialsByTypeFromEnv()
	tastyTradeCredentials, tastyTradeAPIURL, tastyTradeTokenManager := LoadTastyTradeConfiguration(credentialsByType, "tastytrade")
	tastyTradeSandboxCredentials, tastyTradeSandboxAPIURL, tastyTradeSandboxTokenManager := LoadTastyTradeConfiguration(credentialsByType, "tastytrade_sandbox")
	brokerAuthorizationCredentials := map[broker.AccountType]auth.Credentials{
		broker.AccountTypeTastyTrade:        tastyTradeCredentials,
		broker.AccountTypeTastyTradeSandbox: tastyTradeSandboxCredentials,
	}
	handler := httpapi.NewHandler(httpapi.NewHandlerInput{
		OAuthStateStore:       oauthstatestore.NewInMemory(),
		PendingSelectionStore: pendingselectionstore.NewInMemory(),
		AccountStoreCommandHandler: accountstore.NewCommandHandlerThreadSafeDecorator(accountstore.NewCommandHandlerThreadSafeDecoratorInput{
			Decorated: accountstore.NewEventSourcedCommandHandler(accountstore.NewEventSourcedCommandHandlerInput{
				Log: log,
			}),
		}),
		AccountStoreQueryHandler: accountstore.NewQueryHandlerThreadSafeDecorator(accountstore.NewQueryHandlerThreadSafeDecoratorInput{
			Decorated: accountstore.NewEventSourcedQueryHandler(accountstore.NewEventSourcedQueryHandlerInput{
				Log: log,
			}),
		}),
		BrokerAccountClientFactory: &BrokerAccountClientFactory{
			TastyTradeClientFactory: &tastytrade.HTTPClientFactory{
				APIURL:         tastyTradeAPIURL,
				GetAccessToken: tastyTradeTokenManager.GetAccessToken,
			},
			TastyTradeSandboxClientFactory: &tastytrade.HTTPClientFactory{
				APIURL:         tastyTradeSandboxAPIURL,
				GetAccessToken: tastyTradeSandboxTokenManager.GetAccessToken,
			},
		},
		BrokerOnBoardingClientFactory: &BrokerOnboardingClientFactory{
			BackendRedirectURI: authorizationRedirectURI.String(),
			CredentialsByType:  brokerAuthorizationCredentials,
		},
		AuthMiddleware:     auth.MiddlewareFromEnv(auth.AudienceAccountService),
		BackendRedirectURI: authorizationRedirectURI.String(),
		FrontendBaseURL:    config.EnvStringOrFatal("FRONTEND_BASE_URL"),
	})
	http.ListenAndServe(":9000", handler)
}

func LoadTastyTradeConfiguration(credentialsByType map[string]auth.Credentials, brokerType string) (auth.Credentials, *url.URL, *auth.TastyTradeTokenManager) {
	credentials, ok := credentialsByType[brokerType]
	fatal.Unlessf(ok, "credentials not configured for broker type %s", brokerType)
	apiURL, err := url.Parse(credentials.APIURL)
	fatal.OnError(err)
	tokenManager := auth.NewTastyTradeTokenManager(&credentials.AuthorizationServer)
	return credentials, apiURL, tokenManager
}
