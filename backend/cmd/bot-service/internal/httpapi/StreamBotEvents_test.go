package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/kduong-dev/goutil/eventsource"
	"github.com/kduong-dev/trading-core/backend/cmd/bot-service/internal/botstore"
	. "github.com/smartystreets/goconvey/convey"
)

type fakeBotStoreQueryHandler struct{}

func (fakeBotStoreQueryHandler) Get(ctx context.Context, botID string) (*botstore.Bot, error) {
	return &botstore.Bot{ID: botID}, nil
}

func (fakeBotStoreQueryHandler) List(ctx context.Context) ([]*botstore.Bot, error) {
	return nil, nil
}

func TestStreamBotEvents(t *testing.T) {
	Convey("Given a client streaming a bot's events", t, func() {
		api := &API{
			botStoreQueryHandler: fakeBotStoreQueryHandler{},
			botEventLogFactory:   eventsource.NewInMemoryLogFactory(),
			botChannelFunc:       func(botID string) string { return "bot:" + botID },
		}
		ctx, cancel := context.WithCancel(context.Background())
		request := mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/bots/v1/bot-1/events", nil).WithContext(ctx), map[string]string{"bot_id": "bot-1"})
		recorder := httptest.NewRecorder()
		finished := make(chan struct{})
		go func() {
			api.StreamBotEvents(recorder, request)
			close(finished)
		}()
		Convey("When the client disconnects", func() {
			cancel()
			Convey("Then the stream ends without an error response", func() {
				select {
				case <-finished:
				case <-time.After(15 * time.Second):
					t.Fatal("stream didn't end after the client disconnected")
				}
				So(recorder.Code, ShouldEqual, http.StatusOK)
				So(recorder.Header().Get("Content-Type"), ShouldEqual, "text/event-stream")
			})
		})
	})
}
