package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ansel1/merry"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/kduong/trading-backend/cmd/authentication-service/internal/userstore"
	"github.com/kduong/trading-backend/internal/fatal"
	"github.com/kduong/trading-backend/internal/httpx"
)

type CreateUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const (
	MinimumPasswordLength = 8
	// bcrypt only hashes the first 72 bytes and rejects anything longer.
	MaximumPasswordBytes = 72
)

// CreateUserOutput is the user as the client sees it: never the password hash.
type CreateUserOutput struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (handler *Handler) CreateUser(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpx.SendErrorResponse(responseWriter, err)
		}
	}()
	ctx := request.Context()
	input, err := httpx.DecodeJSONBody[CreateUserInput](request)
	if err != nil {
		return
	}
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if len(email) == 0 {
		err = merry.New("email is required").WithHTTPCode(http.StatusBadRequest).WithUserMessage("email is required")
		return
	}
	if len([]rune(input.Password)) < MinimumPasswordLength {
		message := fmt.Sprintf("password must be at least %d characters", MinimumPasswordLength)
		err = merry.New(message).WithHTTPCode(http.StatusBadRequest).WithUserMessage(message)
		return
	}
	if len(input.Password) > MaximumPasswordBytes {
		message := fmt.Sprintf("password must be at most %d bytes", MaximumPasswordBytes)
		err = merry.New(message).WithHTTPCode(http.StatusBadRequest).WithUserMessage(message)
		return
	}
	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return
	}
	object := userstore.User{
		ID:           uuid.NewV4().String(),
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
	}
	err = handler.userStore.Put(ctx, object)
	if err != nil {
		if errors.Is(err, userstore.ErrAlreadyExists) {
			err = merry.Wrap(err).WithHTTPCode(http.StatusConflict).WithUserMessage("user already exists")
			return
		}
		return
	}
	output := CreateUserOutput{
		ID:        object.ID,
		Email:     object.Email,
		CreatedAt: object.CreatedAt,
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(responseWriter).Encode(&output)
	fatal.OnErrorUnlessDone(ctx, err)
}

func HashPassword(password string) (hashPassword string, err error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	hashPassword = string(hashBytes)
	return
}
