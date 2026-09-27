package httpapi

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kduong-dev/goutil/httpx"
	"github.com/kduong-dev/trading-core/backend/cmd/authentication-service/internal/userstore"
	"github.com/kduong-dev/trading-core/backend/internal/auth"
	"github.com/kduong-dev/trading-core/backend/internal/authz"
)

type API struct {
	userStore   userstore.Store
	tokenSecret []byte
	expiryTTL   time.Duration
}

func (api *API) GenerateToken(user *userstore.User) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(api.expiryTTL)
	claims := auth.Claims{
		Scope: strings.Join(authz.UserScopes, " "),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   user.ID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(api.tokenSecret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

var merrifiedSentinels = httpx.MerrifiedSentinels{
	{Sentinel: userstore.ErrNotFound, StatusCode: http.StatusUnauthorized, UserMessage: "invalid credentials"},
	{Sentinel: userstore.ErrAlreadyExists, StatusCode: http.StatusConflict, UserMessage: "user already exists"},
}
