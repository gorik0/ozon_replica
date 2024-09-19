package notifications

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"ozon_replic/internal/models/models"
)

//go:generate mockgen -source interfaces.go -destination ./mocks/notifications_mock.go -package mock

type NotificationsUsecase interface {
	GetDayNotifications(context.Context, uuid.UUID) ([]models.Message, error)
}

type NotificationsRepo interface {
	ReadDayNotifications(context.Context, uuid.UUID) ([]models.Message, error)
}
