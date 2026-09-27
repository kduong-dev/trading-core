package httpapi

import (
	"net/http"
	"strings"
	"time"

	"github.com/ansel1/merry"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kduong-dev/goutil/fatal"
	"github.com/kduong-dev/goutil/httpx"
)

func (api *API) RefreshSession(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	authorization := request.Header.Get("Authorization")
	if authorization == "" {
		err = merry.New("missing authorization header").WithHTTPCode(http.StatusUnauthorized).WithUserMessage("unauthorized")
		return
	}
	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		err = merry.New("invalid authorization header").WithHTTPCode(http.StatusUnauthorized).WithUserMessage("unauthorized")
		return
	}
	tokenString := parts[1]
	claims := &jwt.RegisteredClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, merry.New("unexpected signing method").WithHTTPCode(http.StatusUnauthorized).WithUserMessage("unauthorized")
		}
		return api.tokenSecret, nil
	})
	if err != nil {
		err = merry.Wrap(err).WithHTTPCode(http.StatusUnauthorized).WithUserMessage("unauthorized")
		return
	}
	userID := claims.Subject
	object, err := api.userStore.GetByID(ctx, userID)
	if err = merrifiedSentinels.MerrifyOrFatal(err); err != nil {
		return
	}
	token, expiresAt, err := api.GenerateToken(object)
	fatal.OnError(err)
	output := CreateSessionOutput{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt.Format(time.RFC3339),
		UserID:      object.ID,
		Email:       object.Email,
	}
	httpx.SendJSONResponse(responseWriter, http.StatusOK, output)
}
