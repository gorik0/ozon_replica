package usecase

import (
	"context"
	"errors"
	"fmt"
	"ozon_replic/internal/models/models"
	"ozon_replic/internal/pkg/cart"
	"ozon_replic/internal/pkg/cart/repo"

	uuid "github.com/satori/go.uuid"
)

type CartUsecase struct {
	repo cart.CartRepo
}

func NewCartUsecase(repo cart.CartRepo) *CartUsecase {
	return &CartUsecase{
		repo: repo,
	}
}

func (uc *CartUsecase) GetCart(ctx context.Context, id uuid.UUID) (models.Cart, error) {
	cart, err := uc.repo.ReadCart(ctx, id)

	if err != nil {
		if errors.Is(err, repo.ErrCartNotFound) {
			cart.Id, err = uc.repo.CreateCart(ctx, id)
			cart.IsCurrent = true
			cart.Products = make([]models.CartProduct, 0)
			if err != nil {
				return models.Cart{}, fmt.Errorf("error happened in repo.CreateCart: %w", err)

			}
			return cart, nil
		}
		return models.Cart{}, fmt.Errorf("error happened in repo.ReadCart: %w", err)
	}

	return cart, nil
}

func (uc *CartUsecase) AddProduct(ctx context.Context, cart models.Cart, product models.CartProductUpdate) (models.Cart, error) {
	cartChecked, err := uc.repo.CheckCart(ctx, cart.ProfileId)

	if err != nil {
		if errors.Is(err, repo.ErrCartNotFound) {
			cart.Id, err = uc.repo.CreateCart(ctx, cart.ProfileId)
			cart.IsCurrent = true
			if err != nil {
				return models.Cart{}, fmt.Errorf("error happened in repo.GetCart: %w", err)
			}
			return uc.repo.AddProduct(ctx, cart, product)
		}
		return models.Cart{}, fmt.Errorf("error happened in repo.CheckCart: %w", err)
	}
	cart, err = uc.repo.AddProduct(ctx, cartChecked, product)
	if err != nil {
		return models.Cart{}, fmt.Errorf("error happened in repo.AddProduct: %w", err)
	}
	return cart, err

}

func (uc *CartUsecase) DeleteProduct(ctx context.Context, cart models.Cart, product models.CartProductDelete) (models.Cart, error) {
	cart, err := uc.repo.CheckCart(ctx, cart.ProfileId)
	if err != nil {
		return models.Cart{}, fmt.Errorf("error happened in repo.CheckCart: %w", err)
	}
	cart, err = uc.repo.DeleteProduct(ctx, cart, product)
	if err != nil {
		return models.Cart{}, fmt.Errorf("error happened in repo.DeleteProduct: %w", err)
	}
	return cart, err
}
func (uc *CartUsecase) UpdateCart(ctx context.Context, cart models.Cart) (models.Cart, error) {
	_, err := uc.repo.ReadCart(ctx, cart.ProfileId)
	if errors.Is(err, repo.ErrCartNotFound) {
		cart.Id, err = uc.repo.CreateCart(ctx, cart.ProfileId)
		cart.IsCurrent = true
		if err != nil {
			return models.Cart{}, err
		}
	}

	if err != nil {
		return models.Cart{}, err
	}

	cart, err = uc.repo.UpdateCart(ctx, cart)
	if err != nil {
		err = fmt.Errorf("error happened in repo.UpdateCart: %w", err)

		return models.Cart{}, err
	}

	return cart, nil
}
