package main

import (
	"net/http"

	"github.com/kduong-dev/goutil/eventsource"
	"github.com/kduong-dev/goutil/fatal"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/journal-service/internal/entrystore"
	"github.com/kduong-dev/trading-core/backend/cmd/journal-service/internal/httpapi"
	"github.com/kduong-dev/trading-core/backend/internal/auth"
)

func main() {
	logFactory, err := eventsource.LogFactoryFromEnv("TRADING_JOURNAL_EVENT_LOG", "INMEMORY")
	fatal.OnError(err)
	log, err := logFactory.Create("trading_journal:events")
	fatal.OnError(err)
	commandHandler := entrystore.NewCommandHandlerThreadSafeDecorator(entrystore.NewCommandHandlerThreadSafeDecoratorInput{
		Decorated: entrystore.NewEventSourcedCommandHandler(entrystore.NewEventSourcedCommandHandlerInput{
			Log: log,
		}),
	})
	queryHandler := entrystore.NewQueryHandlerThreadSafeDecorator(entrystore.NewQueryHandlerThreadSafeDecoratorInput{
		Decorated: entrystore.NewEventSourcedQueryHandler(entrystore.NewEventSourcedQueryHandlerInput{
			Log: log,
		}),
	})
	router := httpapi.NewRouter(httpapi.NewRouterInput{
		AuthMiddleware:      auth.MiddlewareFromEnv(auth.AudienceJournalService),
		EntryCommandHandler: commandHandler,
		EntryQueryHandler:   queryHandler,
	})
	err = http.ListenAndServe(":8084", httpx.HandlerWithCORS(router))
	fatal.OnError(err)
}
