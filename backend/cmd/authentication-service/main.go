package main

import (
	"context"
	"net/http"
	"time"

	"github.com/kduong-dev/goutil/config"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/authentication-service/internal/httpapi"
	"github.com/kduong-dev/trading-core/backend/cmd/authentication-service/internal/userstore"
)

func main() {
	ctx := context.Background()
	router := httpapi.NewRouter(httpapi.NewRouterInput{
		UserStore: userstore.NewThreadSafeDecorator(userstore.NewThreadSafeDecoratorInput{
			Decorated: userstore.FromEnv(ctx),
		}),
		TokenSecret: []byte(config.EnvStringOrFatal("TOKEN_SECRET")),
		ExpiryTTL:   1 * time.Hour,
	})
	http.ListenAndServe(":9100", httpx.HandlerWithCORS(router))
}
