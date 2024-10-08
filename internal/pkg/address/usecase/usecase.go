package usecase

import (
	"context"
	"errors"
	"fmt"
	"ozon_replic/internal/models/models"
	"ozon_replic/internal/pkg/address"
	"ozon_replic/internal/pkg/address/repo"

	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
)

type AddressUsecase struct {
	repo address.AddressRepo
}

func NewAddressUsecase(repoAddress address.AddressRepo) *AddressUsecase {
	return &AddressUsecase{
		repo: repoAddress,
	}
}

func (uc *AddressUsecase) AddAddress(ctx context.Context, userID uuid.UUID, addressInfo models.AddressPayload) (models.Address, error) {
	if err := validator.New().Struct(addressInfo); err != nil {
		return models.Address{}, err
	}
	addressInfo.Sanitize()

	address, err := uc.repo.CreateAddress(ctx, userID, addressInfo)
	if err != nil {
		return models.Address{}, fmt.Errorf("error happened in repo.CreateAddress: %w", err)
	}
	return address, nil
}

func (uc *AddressUsecase) UpdateAddress(ctx context.Context, addressInfo models.Address) (models.Address, error) {
	addressInfo.Sanitize()

	if err := uc.repo.UpdateAddress(ctx, addressInfo); err != nil {
		return models.Address{}, fmt.Errorf("error happened in repo.UpdateAddress: %w", err)
	}

	address, err := uc.repo.ReadAddress(ctx, addressInfo.ProfileId, addressInfo.Id)
	if err != nil {
		if errors.Is(err, repo.ErrAddressNotFound) {
			return models.Address{}, err
		}
		return models.Address{}, fmt.Errorf("error happened in repo.ReadAddress: %w", err)
	}

	return address, nil
}

func (uc *AddressUsecase) DeleteAddress(ctx context.Context, addressInfo models.AddressDelete) error {
	if err := uc.repo.DeleteAddress(ctx, addressInfo); err != nil {
		if errors.Is(err, repo.ErrNoCurrentAddressNotFound) {
			return err
		}
		return fmt.Errorf("error happened in repo.DeleteAddress: %w", err)
	}
	return nil
}

func (uc *AddressUsecase) MakeCurrentAddress(ctx context.Context, addressInfo models.AddressMakeCurrent) error {
	if err := uc.repo.MakeCurrentAddress(ctx, addressInfo); err != nil {
		if errors.Is(err, repo.ErrCurrentAddressNotFound) {
			return err
		}
		return fmt.Errorf("error happened in repo.MakeCurrentAddress: %w", err)
	}
	return nil
}

func (uc *AddressUsecase) GetCurrentAddress(ctx context.Context, userID uuid.UUID) (models.Address, error) {
	address, err := uc.repo.ReadCurrentAddress(ctx, userID)
	if err != nil {
		if errors.Is(err, repo.ErrAddressNotFound) {
			return models.Address{}, err
		}
		return models.Address{}, fmt.Errorf("error happened in repo.ReadCurrentAddress: %w", err)
	}
	return address, nil
}

func (uc *AddressUsecase) GetAllAddresses(ctx context.Context, userID uuid.UUID) ([]models.Address, error) {
	address, err := uc.repo.ReadAllAddresses(ctx, userID)
	if err != nil {
		if errors.Is(err, repo.ErrAddressesNotFound) {
			return []models.Address{}, err
		}
		return []models.Address{}, fmt.Errorf("error happened in repo.ReadAllAddresses: %w", err)
	}
	return address, nil
}
