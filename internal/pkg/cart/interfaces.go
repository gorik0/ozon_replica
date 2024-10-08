package cart

import (
	"context"
	"ozon_replic/internal/models/models"

	"github.com/satori/go.uuid"
)

//go:generate mockgen -source interfaces.go -destination ./mocks/cart_mock.go -package mock

type CartUsecase interface {
	GetCart(context.Context, uuid.UUID) (models.Cart, error)
	AddProduct(context.Context, models.Cart, models.CartProductUpdate) (models.Cart, error)
	DeleteProduct(context.Context, models.Cart, models.CartProductDelete) (models.Cart, error)
	UpdateCart(context.Context, models.Cart) (models.Cart, error)
}

type CartRepo interface {
	CreateCart(context.Context, uuid.UUID) (uuid.UUID, error)
	CheckCart(context.Context, uuid.UUID) (models.Cart, error)
	ReadCart(context.Context, uuid.UUID) (models.Cart, error)
	ReadCartProducts(context.Context, models.Cart) (models.Cart, error)
	AddProduct(context.Context, models.Cart, models.CartProductUpdate) (models.Cart, error)
	DeleteProduct(context.Context, models.Cart, models.CartProductDelete) (models.Cart, error)
	DeleteCart(ctx context.Context, cartID uuid.UUID) error
	UpdateCart(context.Context, models.Cart) (models.Cart, error)
}
