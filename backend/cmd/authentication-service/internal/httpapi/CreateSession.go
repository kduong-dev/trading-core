package httpapi

import (
	"net/http"
	"strings"
	"time"

	"github.com/ansel1/merry"
	"github.com/kduong-dev/goutil/fatal"
	"github.com/kduong-dev/goutil/httpx"
	"golang.org/x/crypto/bcrypt"
)

type CreateSessionInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateSessionOutput struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   string `json:"expires_at"`
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
}

func (api *API) CreateSession(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	input, err := httpx.DecodeJSONBody[CreateSessionInput](request)
	if err != nil {
		return
	}
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if len(input.Email) == 0 || len(input.Password) == 0 {
		err = merry.UserError("email and password are required").WithHTTPCode(http.StatusBadRequest)
		return
	}
	object, err := api.userStore.GetByEmail(ctx, input.Email)
	if err = merrifiedSentinels.MerrifyOrFatal(err); err != nil {
		return
	}
	isPasswordValid := VerifyPassword(input.Password, object.PasswordHash)
	if !isPasswordValid {
		err = merry.New("invalid password").WithHTTPCode(http.StatusUnauthorized).WithUserMessage("invalid credentials")
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

func VerifyPassword(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
