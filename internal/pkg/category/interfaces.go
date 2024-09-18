package category

import (
	"context"
	"ozon_replic/internal/models/models"
)

//go:generate mockgen -source interfaces.go -destination ./mocks/category_mock.go -package mock

type CategoryUsecase interface {
	Categories(context.Context) (models.CategoryTree, error)
}

type CategoryRepo interface {
	ReadCategories(context.Context) (models.CategoryTree, error)
}
