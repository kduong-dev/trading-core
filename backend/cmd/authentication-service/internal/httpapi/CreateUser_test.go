package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kduong-dev/trading-core/backend/cmd/authentication-service/internal/httpapi"
	"github.com/kduong-dev/trading-core/backend/cmd/authentication-service/internal/userstore"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateUser(t *testing.T) {
	Convey("Given the authentication handler", t, func() {
		userStore := userstore.NewInMemoryStore()
		handler := httpapi.NewHandler(httpapi.NewHandlerInput{
			UserStore:   userStore,
			TokenSecret: []byte("test-secret"),
			ExpiryTTL:   time.Hour,
		})

		Convey("When a user registers", func() {
			body, _ := json.Marshal(httpapi.CreateUserInput{Email: " New.User@Example.com ", Password: "a-long-password"})
			request := httptest.NewRequest(http.MethodPost, "/auth/v1/users", bytes.NewReader(body))
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)

			Convey("Then it responds with the created user", func() {
				So(recorder.Code, ShouldEqual, http.StatusOK)
				var output httpapi.CreateUserOutput
				So(json.Unmarshal(recorder.Body.Bytes(), &output), ShouldBeNil)
				So(output.ID, ShouldNotBeEmpty)
				So(output.Email, ShouldEqual, "new.user@example.com")
			})

			Convey("Then the response does not contain the password hash", func() {
				var fields map[string]any
				So(json.Unmarshal(recorder.Body.Bytes(), &fields), ShouldBeNil)
				So(fields, ShouldHaveLength, 3)
				So(fields, ShouldContainKey, "id")
				So(fields, ShouldContainKey, "email")
				So(fields, ShouldContainKey, "created_at")
				So(recorder.Body.String(), ShouldNotContainSubstring, "$2a$")
			})

			Convey("Then the stored user keeps its password hash", func() {
				user, err := userStore.GetByEmail(request.Context(), "new.user@example.com")
				So(err, ShouldBeNil)
				So(user.PasswordHash, ShouldStartWith, "$2a$")
			})
		})

		rejected := map[string]httpapi.CreateUserInput{
			"a blank email":            {Email: "  ", Password: "a-long-password"},
			"a 7 character password":   {Email: "short@example.com", Password: "seven77"},
			"a password over 72 bytes": {Email: "long@example.com", Password: strings.Repeat("a", 73)},
		}
		for description, input := range rejected {
			Convey("When a user registers with "+description, func() {
				body, _ := json.Marshal(input)
				request := httptest.NewRequest(http.MethodPost, "/auth/v1/users", bytes.NewReader(body))
				recorder := httptest.NewRecorder()
				handler.ServeHTTP(recorder, request)

				Convey("Then it responds with bad request and stores nothing", func() {
					So(recorder.Code, ShouldEqual, http.StatusBadRequest)
					_, err := userStore.GetByEmail(request.Context(), strings.TrimSpace(input.Email))
					So(err, ShouldNotBeNil)
				})
			})
		}

		Convey("When a user registers with an 8 character password", func() {
			body, _ := json.Marshal(httpapi.CreateUserInput{Email: "eight@example.com", Password: "eight888"})
			request := httptest.NewRequest(http.MethodPost, "/auth/v1/users", bytes.NewReader(body))
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)

			Convey("Then it is accepted", func() {
				So(recorder.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}
