package auth

import (
	"context"
	"ozon_replic/internal/models/models"
	"time"

	uuid "github.com/satori/go.uuid"
)

//go:generate mockgen -source interfaces.go -destination ./mocks/auth_mock.go -package mock

type AuthUsecase interface {
	SignIn(context.Context, *models.SignInPayload) (*models.Profile, string, time.Time, error)
	SignUp(context.Context, *models.SignUpPayload) (*models.Profile, string, time.Time, error)
	CheckAuth(context.Context, uuid.UUID) (*models.Profile, error)
}
