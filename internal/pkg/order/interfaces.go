package order

import (
	"context"
	"ozon_replic/internal/models/models"
	"time"

	uuid "github.com/satori/go.uuid"
)

//go:generate mockgen -source interfaces.go -destination ./mocks/order_mock.go -package mock

type OrderUsecase interface {
	CreateOrder(context.Context, uuid.UUID, string, string, string) (models.Order, error)
	GetCurrentOrder(context.Context, uuid.UUID) (models.Order, error)
	GetOrders(context.Context, uuid.UUID) ([]models.Order, error)
}

type OrderRepo interface {
	CreateOrder(context.Context, models.Cart, uuid.UUID, uuid.UUID, int64, string, string) (models.Order, error)
	ReadOrderID(context.Context, uuid.UUID) (uuid.UUID, error)
	ReadOrder(context.Context, uuid.UUID) (models.Order, error)
	ReadOrdersID(context.Context, uuid.UUID) ([]uuid.UUID, error)
	GetUpdates(context.Context, uuid.UUID, time.Time) ([]models.Message, error)
	SetPromoOrder(context.Context, int, uuid.UUID) error
	SetOrderStatus(context.Context, time.Time) error
}
