package main

import (
	"context"
	"net/http"
	"os"

	"github.com/kduong-dev/goutil/config"
	"github.com/kduong-dev/goutil/eventsource"
	"github.com/kduong-dev/goutil/fatal"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/storage-service/pkg/storageservice"
	"github.com/kduong-dev/trading-core/backend/cmd/reporting-service/internal/httpapi"
	"github.com/kduong-dev/trading-core/backend/cmd/reporting-service/internal/jobstore"
	"github.com/kduong-dev/trading-core/backend/cmd/reporting-service/internal/jobsync"
	"github.com/kduong-dev/trading-core/backend/internal/auth"
)

func main() {
	ctx := context.Background()
	outputsDirectory := config.EnvString("REPORTING_OUTPUTS_DIRECTORY", "./tmp/reports")
	err := os.MkdirAll(outputsDirectory, 0o755)
	fatal.OnError(err)
	logFactory, err := eventsource.LogFactoryFromEnv("TRADING_REPORT_EVENT_LOG", "INMEMORY")
	fatal.OnError(err)
	log, err := logFactory.Create("trading_report:events")
	fatal.OnError(err)
	commandHandler := jobstore.NewCommandHandlerThreadSafeDecorator(jobstore.NewCommandHandlerThreadSafeDecoratorInput{
		Decorated: jobstore.NewEventSourcedCommandHandler(jobstore.NewEventSourcedCommandHandlerInput{
			Log: log,
		}),
	})
	queryHandler := jobstore.NewQueryHandlerThreadSafeDecorator(jobstore.NewQueryHandlerThreadSafeDecoratorInput{
		Decorated: jobstore.NewEventSourcedQueryHandler(jobstore.NewEventSourcedQueryHandlerInput{
			Log: log,
		}),
	})
	storageClient := storageservice.ClientFromEnv()
	actor := jobsync.NewActor(jobsync.NewActorInput{
		CommandHandler:   commandHandler,
		StorageClient:    storageClient,
		OutputsDirectory: outputsDirectory,
		Log:              log,
	})
	actor.CatchUp(ctx)
	actor.CompleteCatchup(ctx)
	go actor.Run(ctx)
	router := httpapi.NewRouter(httpapi.NewRouterInput{
		AuthMiddleware:    auth.MiddlewareFromEnv(auth.AudienceReportingService),
		JobCommandHandler: commandHandler,
		JobQueryHandler:   queryHandler,
		StorageClient:     storageClient,
		EnqueueJob:        actor.Notify,
	})
	err = http.ListenAndServe(":8082", httpx.HandlerWithCORS(router))
	fatal.OnError(err)
}
