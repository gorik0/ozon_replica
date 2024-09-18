package usecase

import (
	"context"
	"fmt"
	"ozon_replic/internal/models/models"
	"ozon_replic/internal/pkg/category"
)

type CategoryUsecase struct {
	repo category.CategoryRepo
}

func NewCategoryUsecase(repo category.CategoryRepo) *CategoryUsecase {
	return &CategoryUsecase{
		repo: repo,
	}
}

func (uc *CategoryUsecase) Categories(ctx context.Context) (models.CategoryTree, error) {
	tree, err := uc.repo.ReadCategories(ctx)
	if err != nil {
		err = fmt.Errorf("error happened in repo.ReadCategories: %w", err)

		return models.CategoryTree{}, err
	}

	return tree, nil
}
