package usecase

import (
	"context"
	"errors"
	"fmt"
	"ozon_replic/internal/models/models"
	"ozon_replic/internal/pkg/products"
	"ozon_replic/internal/pkg/search"
	"ozon_replic/internal/pkg/search/repo"
)

type SearchUsecase struct {
	repoSearch   search.SearchRepo
	repoProducts products.ProductsRepo
}

func NewSearchUsecase(repoSearch search.SearchRepo, repoProducts products.ProductsRepo) *SearchUsecase {
	return &SearchUsecase{
		repoSearch:   repoSearch,
		repoProducts: repoProducts,
	}
}

func (uc *SearchUsecase) SearchProducts(ctx context.Context, productName string) ([]models.Product, error) {
	var err error
	var productsSlice []models.Product
	if productName == "" {
		productsSlice, err = uc.repoProducts.ReadProducts(ctx, 0, 10)
		if err != nil {
			if errors.Is(err, repo.ErrProductNotFound) {
				return []models.Product{}, err
			}
			err = fmt.Errorf("error happened in repoProducts.ReadProducts: %w", err)

			return []models.Product{}, err
		}

		return productsSlice, nil
	}
	productsSlice, err = uc.repoSearch.ReadProductsByName(ctx, productName)
	if err != nil {
		if errors.Is(err, repo.ErrProductNotFound) {
			return []models.Product{}, err
		}
		err = fmt.Errorf("error happened in repoSearch.SearchProducts: %w", err)

		return []models.Product{}, err
	}

	return productsSlice, nil
}
