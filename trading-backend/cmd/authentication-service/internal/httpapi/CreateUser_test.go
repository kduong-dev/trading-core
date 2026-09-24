package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kduong/trading-backend/cmd/authentication-service/internal/httpapi"
	"github.com/kduong/trading-backend/cmd/authentication-service/internal/userstore"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateUser(t *testing.T) {
	Convey("Given the authentication router", t, func() {
		userStore := userstore.NewInMemoryStore()
		router := httpapi.NewRouter(httpapi.NewRouterInput{
			UserStore:   userStore,
			TokenSecret: []byte("test-secret"),
			ExpiryTTL:   time.Hour,
		})

		Convey("When a user registers", func() {
			body, _ := json.Marshal(httpapi.CreateUserInput{Email: " New.User@Example.com ", Password: "a-long-password"})
			request := httptest.NewRequest(http.MethodPost, "/auth/v1/users", bytes.NewReader(body))
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

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
	})
}
