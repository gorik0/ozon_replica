package search

import (
	"context"
	"ozon_replic/internal/models/models"
)

//go:generate mockgen -source interfaces.go -destination ./mocks/search_mock.go -package mock

type SearchUsecase interface {
	SearchProducts(context.Context, string) ([]models.Product, error)
}

type SearchRepo interface {
	ReadProductsByName(context.Context, string) ([]models.Product, error)
}
