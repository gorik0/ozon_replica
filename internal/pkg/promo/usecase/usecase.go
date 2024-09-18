package usecase

import (
	"context"
	"ozon_replic/internal/models/models"
	"ozon_replic/internal/pkg/promo"
	"time"

	uuid "github.com/satori/go.uuid"
)

type PromoUsecase struct {
	repo promo.PromoRepo
}

func NewPromoUsecase(repo promo.PromoRepo) *PromoUsecase {
	return &PromoUsecase{
		repo: repo,
	}
}

func (uc *PromoUsecase) CheckPromocode(ctx context.Context, userID uuid.UUID, name string) (*models.Promocode, error) {
	promocode, err := uc.repo.ReadPromocode(ctx, name)
	if err != nil {
		return &models.Promocode{}, err
	}

	if time.Now().After(promocode.Deadline) {
		return &models.Promocode{}, promo.ErrPromocodeExpired
	}
	if promocode.Leftover < 1 {
		return &models.Promocode{}, promo.ErrPromocodeLeftout
	}

	if err := uc.repo.CheckUniq(ctx, userID, int(promocode.Id)); err != nil {
		return &models.Promocode{}, err
	}

	return promocode, nil
}

func (uc *PromoUsecase) UsePromocode(ctx context.Context, name string) (*models.Promocode, error) {
	promocode, err := uc.repo.UsePromocode(ctx, name)
	if err != nil {
		return &models.Promocode{}, err
	}
	if time.Now().After(promocode.Deadline) {
		return &models.Promocode{}, promo.ErrPromocodeExpired
	}
	if promocode.Leftover < 1 {
		return &models.Promocode{}, promo.ErrPromocodeLeftout
	}

	return promocode, nil
}
