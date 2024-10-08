package grpc

import (
	"context"
	"errors"
	"log"
	"log/slog"
	addressRepo "ozon_replic/internal/pkg/address/repo"
	cartRepo "ozon_replic/internal/pkg/cart/repo"
	"ozon_replic/internal/pkg/middleware/metricsmw"
	"ozon_replic/internal/pkg/order"
	"ozon_replic/internal/pkg/order/delivery/grpc/gen"
	orderRepo "ozon_replic/internal/pkg/order/repo"
	"ozon_replic/internal/pkg/promo"
	"ozon_replic/internal/pkg/utils/logger/sl"
	gmodels "ozon_replic/proto/gmodels"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	uuid "github.com/satori/go.uuid"
)

//go:generate mockgen -source=./gen/order_grpc.pb.go -destination=../../mocks/order_grpc.go -package=mock

type GrpcOrderHandler struct {
	log *slog.Logger
	uc  order.OrderUsecase

	gen.OrderServer
}

func NewGrpcOrderHandler(uc order.OrderUsecase, log *slog.Logger) *GrpcOrderHandler {
	return &GrpcOrderHandler{
		uc:  uc,
		log: log,
	}
}

func (h GrpcOrderHandler) CreateOrder(ctx context.Context, in *gen.CreateOrderRequest) (*gen.CreateOrderResponse, error) {
	h.log = h.log.With(
		slog.String("op", sl.GFN()),
	)

	userId, err := uuid.FromString(in.Id)
	log.Println("USERid :::", userId)
	if err != nil {
		h.log.Error("failed to get uuid from string", sl.Err(err))
		return nil, metricsmw.ClientError
	}
	order, err := h.uc.CreateOrder(ctx, userId, in.DeliveryTime, in.DeliveryDate, in.PromocodeName)

	switch {
	case err == nil:
		break
	case errors.Is(err, promo.ErrPromocodeNotFound):
		h.log.Error("failed to CreateOrder", sl.Err(err))
		return nil, status.Error(codes.NotFound, err.Error())
	case errors.Is(err, promo.ErrPromocodeLeftout):
		h.log.Error("failed to CreateOrder", sl.Err(err))
		return nil, status.Error(codes.OutOfRange, err.Error())
	case errors.Is(err, promo.ErrPromocodeExpired):
		h.log.Error("failed to CreateOrder", sl.Err(err))
		return nil, status.Error(codes.DeadlineExceeded, err.Error())
	case errors.Is(err, orderRepo.ErrPoductNotFound):
		h.log.Error("failed to CreateOrder", sl.Err(err))
		return nil, status.Error(codes.Unavailable, err.Error())
	case errors.Is(err, cartRepo.ErrCartNotFound):
		h.log.Error("failed to CreateOrder", sl.Err(err))
		return nil, status.Error(codes.Aborted, err.Error())
	case errors.Is(err, addressRepo.ErrAddressNotFound):
		h.log.Error("failed to CreateOrder", sl.Err(err))
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	default:
		h.log.Error("failed in h.uc.CreateOrder", sl.Err(err))
		return nil, metricsmw.ServerError
	}

	gorder := gen.CreateOrderResponse{
		Order: &gmodels.Order{
			Id:            order.Id.String(),
			Status:        order.Status,
			CreationAt:    order.CreationAt.Format(time.RFC3339),
			DeliveryTime:  order.DeliveryTime,
			DeliveryDate:  order.DeliveryDate,
			PromocodeName: order.PomocodeName,
			Address: &gmodels.Address{
				Id:        order.Address.Id.String(),
				ProfileId: order.Address.ProfileId.String(),
				City:      order.Address.City,
				Street:    order.Address.Street,
				House:     order.Address.House,
				Flat:      order.Address.Flat,
				IsCurrent: order.Address.IsCurrent,
			},
			Products: make([]*gmodels.OrderProduct, len(order.Products)),
		},
	}

	for i, product := range order.Products {
		gorder.Order.Products[i] = &gmodels.OrderProduct{
			Quantity: product.Quantity,
			Product: &gmodels.Product{
				Id:          product.Id.String(),
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				ImgSrc:      product.ImgSrc,
				Rating:      product.Rating,
				Category: &gmodels.Category{
					Id:     product.Category.Id,
					Name:   product.Category.Name,
					Parent: product.Category.Parent,
				},
			},
		}
	}

	return &gorder, nil
}

func (h GrpcOrderHandler) GetOrders(ctx context.Context, in *gen.OrdersRequest) (*gen.OrdersResponse, error) {
	h.log = h.log.With(
		slog.String("op", sl.GFN()),
	)

	userId, err := uuid.FromString(in.Id)
	if err != nil {
		h.log.Error("failed to get uuid from string", sl.Err(err))
		return nil, metricsmw.ClientError
	}
	orders, err := h.uc.GetOrders(ctx, userId)
	if err != nil {
		if errors.Is(err, orderRepo.ErrOrdersNotFound) {
			h.log.Error("GetOrders empty order found", sl.Err(err))
			return nil, status.Error(codes.NotFound, "empty order found")
		}
		h.log.Error("failed in h.uc.GetOrders", sl.Err(err))
		return nil, metricsmw.ServerError
	}

	gorders := make([]*gmodels.Order, len(orders))

	for i, order := range orders {
		gorders[i] = &gmodels.Order{
			Id:            order.Id.String(),
			Status:        order.Status,
			CreationAt:    order.CreationAt.Format(time.RFC3339),
			DeliveryTime:  order.DeliveryTime,
			DeliveryDate:  order.DeliveryDate,
			PromocodeName: order.PomocodeName,
			Address: &gmodels.Address{
				Id:        order.Address.Id.String(),
				ProfileId: order.Address.ProfileId.String(),
				City:      order.Address.City,
				Street:    order.Address.Street,
				House:     order.Address.House,
				Flat:      order.Address.Flat,
				IsCurrent: order.Address.IsCurrent,
			},
			Products: make([]*gmodels.OrderProduct, len(order.Products)),
		}

		for j, product := range order.Products {
			gorders[i].Products[j] = &gmodels.OrderProduct{
				Quantity: product.Quantity,
				Product: &gmodels.Product{
					Id:          product.Id.String(),
					Name:        product.Name,
					Description: product.Description,
					Price:       product.Price,
					ImgSrc:      product.ImgSrc,
					Rating:      product.Rating,
					Category: &gmodels.Category{
						Id:     product.Category.Id,
						Name:   product.Category.Name,
						Parent: product.Category.Parent,
					},
				},
			}
		}
	}

	return &gen.OrdersResponse{Orders: gorders}, nil
}
