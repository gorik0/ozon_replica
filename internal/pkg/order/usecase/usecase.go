package usecase

import (
	"context"
	"errors"
	"fmt"
	"ozon_replic/internal/models/models"
	"ozon_replic/internal/pkg/address"
	addressRepo "ozon_replic/internal/pkg/address/repo"

	"ozon_replic/internal/pkg/cart"
	cartRepo "ozon_replic/internal/pkg/cart/repo"
	"ozon_replic/internal/pkg/order"
	orderRepo "ozon_replic/internal/pkg/order/repo"
	"ozon_replic/internal/pkg/promo"
	"time"

	uuid "github.com/satori/go.uuid"
)

type OrderUsecase struct {
	repoOrder   order.OrderRepo
	repoCart    cart.CartRepo
	repoAddress address.AddressRepo
	repoPromo   promo.PromoRepo
}

func NewOrderUsecase(
	repoOrder order.OrderRepo,
	repoCart cart.CartRepo,
	repoAddress address.AddressRepo,
	repoPromo promo.PromoRepo) *OrderUsecase {
	return &OrderUsecase{
		repoOrder:   repoOrder,
		repoCart:    repoCart,
		repoAddress: repoAddress,
		repoPromo:   repoPromo,
	}
}

func (uc *OrderUsecase) CreateOrder(
	ctx context.Context, userID uuid.UUID,
	deliveryTime, deliveryDate, promocodeName string) (models.Order, error) {

	address, err := uc.repoAddress.ReadCurrentAddress(ctx, userID)
	if err != nil {
		if errors.Is(err, addressRepo.ErrAddressesNotFound) {
			return models.Order{}, err
		}
		return models.Order{}, fmt.Errorf("error happened in repoAddress.ReadCurrentAddressID: %w", err)
	}

	cart, err := uc.repoCart.ReadCart(ctx, userID)
	if err != nil {
		if errors.Is(err, cartRepo.ErrCartNotFound) {
			return models.Order{}, err
		}
		return models.Order{}, fmt.Errorf("error happened in repoCart.ReadCart: %w", err)
	}

	order, err := uc.repoOrder.CreateOrder(ctx, cart, address.Id, userID, 1, deliveryTime, deliveryDate) //Status ID
	if err != nil {
		if errors.Is(err, orderRepo.ErrPoductNotFound) {
			return models.Order{}, err
		}
		return models.Order{}, fmt.Errorf("error happened in repo.CreateOrder: %w", err)
	}

	if promocodeName != "" {
		promocode, err := uc.repoPromo.UsePromocode(ctx, promocodeName)
		if err != nil {
			return models.Order{}, err
		}
		if time.Now().After(promocode.Deadline) {
			return models.Order{}, promo.ErrPromocodeExpired
		}
		if promocode.Leftover < 1 {
			return models.Order{}, promo.ErrPromocodeLeftout
		}
		if err = uc.repoOrder.SetPromoOrder(ctx, int(promocode.Id), order.Id); err != nil {
			return models.Order{}, err
		}
	}
	order.Address = address
	order.PomocodeName = promocodeName
	order.DeliveryDate = deliveryDate
	order.DeliveryTime = deliveryTime
	order = setStatus(order)

	err = uc.repoCart.DeleteCart(ctx, cart.Id)
	if err != nil {
		err = fmt.Errorf("error happened in repoCart.DeleteCart: %w", err)

		return models.Order{}, err
	}

	return order, nil
}

func setStatus(order models.Order) models.Order {
	switch order.StatusId {
	case 1:
		order.Status = "Создан"
	case 2:
		order.Status = "В обработке"
	default:
		order.Status = "Уточняется"

	}
	return order
}
func (uc *OrderUsecase) GetOrders(ctx context.Context, userID uuid.UUID) ([]models.Order, error) {
	ordersID, err := uc.repoOrder.ReadOrdersID(ctx, userID)
	if err != nil {
		if errors.Is(err, orderRepo.ErrOrdersNotFound) {
			return []models.Order{}, err
		}
		return nil, fmt.Errorf("error happened in repoOrder.ReadOrder: %w", err)
	}

	orders := make([]models.Order, len(ordersID))
	for i, orderID := range ordersID {
		orders[i], err = uc.repoOrder.ReadOrder(ctx, orderID)
		orders[i] = setStatus(orders[i])
		if err != nil {
			if errors.Is(err, orderRepo.ErrPoductsInOrderNotFound) {
				return []models.Order{}, err
			}
			return nil, fmt.Errorf("error happened in repoOrder.ReadOrder: %w", err)
		}
	}

	return orders, nil
}

func (uc *OrderUsecase) GetCurrentOrder(ctx context.Context, userID uuid.UUID) (models.Order, error) {
	orderID, err := uc.repoOrder.ReadOrderID(ctx, userID)
	if err != nil {
		if errors.Is(err, orderRepo.ErrOrderNotFound) {
			return models.Order{}, err
		}
		err = fmt.Errorf("error happened in repoOrder.ReadOrderID: %w", err)

		return models.Order{}, err
	}

	order, err := uc.repoOrder.ReadOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, orderRepo.ErrPoductsInOrderNotFound) {
			return models.Order{}, err
		}
		err = fmt.Errorf("error happened in repoOrder.ReadOrder: %w", err)

		return models.Order{}, err
	}

	order = setStatus(order)

	return order, nil
}
