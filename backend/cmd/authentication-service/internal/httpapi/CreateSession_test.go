package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kduong-dev/trading-core/backend/cmd/authentication-service/internal/httpapi"
	"github.com/kduong-dev/trading-core/backend/cmd/authentication-service/internal/userstore"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateSession(t *testing.T) {
	Convey("Given the authentication handler with a registered user", t, func() {
		handler := httpapi.NewHandler(httpapi.NewHandlerInput{
			UserStore:   userstore.NewInMemoryStore(),
			TokenSecret: []byte("test-secret"),
			ExpiryTTL:   time.Hour,
		})
		body, _ := json.Marshal(httpapi.CreateUserInput{Email: "user@example.com", Password: "a-long-password"})
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/auth/v1/users", bytes.NewReader(body)))
		createSession := func(input httpapi.CreateSessionInput) *httptest.ResponseRecorder {
			body, _ := json.Marshal(input)
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/auth/v1/sessions", bytes.NewReader(body)))
			return recorder
		}

		Convey("When the user signs in with the correct password", func() {
			recorder := createSession(httpapi.CreateSessionInput{Email: "user@example.com", Password: "a-long-password"})

			Convey("Then it responds with an access token", func() {
				So(recorder.Code, ShouldEqual, http.StatusOK)
				var output httpapi.CreateSessionOutput
				So(json.Unmarshal(recorder.Body.Bytes(), &output), ShouldBeNil)
				So(output.AccessToken, ShouldNotBeEmpty)
			})
		})

		Convey("When the user signs in with a blank password", func() {
			recorder := createSession(httpapi.CreateSessionInput{Email: "user@example.com"})

			Convey("Then it responds with bad request", func() {
				So(recorder.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		unauthorized := map[string]httpapi.CreateSessionInput{
			"the wrong password": {Email: "user@example.com", Password: "wrong-password"},
			"an unknown email":   {Email: "unknown@example.com", Password: "a-long-password"},
		}
		for description, input := range unauthorized {
			Convey("When the user signs in with "+description, func() {
				recorder := createSession(input)

				Convey("Then it responds with unauthorized", func() {
					So(recorder.Code, ShouldEqual, http.StatusUnauthorized)
					So(recorder.Body.String(), ShouldContainSubstring, "invalid credentials")
				})
			})
		}
	})
}
