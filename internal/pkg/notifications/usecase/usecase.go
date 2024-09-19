package usecase

import (
	"context"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"ozon_replic/internal/models/models"
	"ozon_replic/internal/pkg/notifications"
	"ozon_replic/internal/pkg/notifications/repo"
)

type NotificationsUsecase struct {
	repo notifications.NotificationsRepo
}

func NewNotificationsUsecase(repoNotifications notifications.NotificationsRepo) *NotificationsUsecase {
	return &NotificationsUsecase{
		repo: repoNotifications,
	}
}

func (uc *NotificationsUsecase) GetDayNotifications(ctx context.Context, userID uuid.UUID) ([]models.Message, error) {
	notifications, err := uc.repo.ReadDayNotifications(ctx, userID)
	if err != nil {
		if errors.Is(err, repo.ErrNotificationsNotFound) {
			return []models.Message{}, err
		}
		err = fmt.Errorf("error happened in repo.ReadDayNotifications: %w", err)

		return []models.Message{}, err
	}

	return notifications, nil
}
