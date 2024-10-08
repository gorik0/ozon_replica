package grpc

import (
	"context"
	"log/slog"
	"ozon_replic/internal/pkg/middleware/metricsmw"
	"ozon_replic/internal/pkg/products"
	"ozon_replic/internal/pkg/products/delivery/grpc/gen"
	"ozon_replic/internal/pkg/utils/logger/sl"
	"ozon_replic/proto/gmodels"

	uuid "github.com/satori/go.uuid"
)

//go:generate mockgen -source=./gen/products_grpc.pb.go -destination=../../mocks/products_grpc.go -package=mock

type GrpcProductsHandler struct {
	uc  products.ProductsUsecase
	log *slog.Logger

	gen.ProductsServer
}

func NewGrpcProductsHandler(uc products.ProductsUsecase, log *slog.Logger) *GrpcProductsHandler {
	return &GrpcProductsHandler{
		uc:  uc,
		log: log,
	}
}

func (h GrpcProductsHandler) GetProduct(ctx context.Context,
	in *gen.ProductRequest) (*gen.ProductResponse, error) {
	h.log = h.log.With(
		slog.String("op", sl.GFN()),
	)

	id, err := uuid.FromString(in.Id)
	if err != nil {
		h.log.Error("failed to get uuid from string", sl.Err(err))
		return nil, metricsmw.ClientError
	}

	product, err := h.uc.GetProduct(ctx, id)
	if err != nil {
		h.log.Error("failed to get product", sl.Err(err))
		return nil, metricsmw.ServerError
	}

	gproduct := &gmodels.Product{
		Id:            product.Id.String(),
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		ImgSrc:        product.ImgSrc,
		Rating:        product.Rating,
		CountComments: product.CountComments,
		Category: &gmodels.Category{
			Id:     product.Category.Id,
			Name:   product.Category.Name,
			Parent: product.Parent,
		},
	}

	return &gen.ProductResponse{Product: gproduct}, nil
}

func (h GrpcProductsHandler) GetProducts(ctx context.Context,
	in *gen.ProductsRequest) (*gen.ProductsResponse, error) {
	h.log = h.log.With(
		slog.String("op", sl.GFN()),
	)

	products, err := h.uc.GetProducts(ctx, in.Paging, in.Count, in.RatingBy, in.PriceBy)
	if err != nil {
		h.log.Error("failed to get products", sl.Err(err))
		return nil, metricsmw.ServerError
	}

	gproducts := make([]*gmodels.Product, len(products))
	for i, product := range products {
		gproduct := &gmodels.Product{
			Id:            product.Id.String(),
			Name:          product.Name,
			Description:   product.Description,
			Price:         product.Price,
			ImgSrc:        product.ImgSrc,
			Rating:        product.Rating,
			CountComments: product.CountComments,
			Category: &gmodels.Category{
				Id:     product.Category.Id,
				Name:   product.Category.Name,
				Parent: product.Parent,
			},
		}
		gproducts[i] = gproduct
	}
	return &gen.ProductsResponse{Products: gproducts}, nil
}

func (h GrpcProductsHandler) GetCategory(ctx context.Context,
	in *gen.CategoryRequest) (*gen.CategoryResponse, error) {
	h.log = h.log.With(
		slog.String("op", sl.GFN()),
	)

	products, err := h.uc.GetCategory(ctx, int(in.Id), in.Paging, in.Count, in.RatingBy, in.PriceBy)
	if err != nil {
		h.log.Error("failed in h.uc.GetCategory", sl.Err(err))
		return nil, metricsmw.ServerError
	}

	gproducts := make([]*gmodels.Product, len(products))
	for i, product := range products {
		gproduct := &gmodels.Product{
			Id:            product.Id.String(),
			Name:          product.Name,
			Description:   product.Description,
			Price:         product.Price,
			ImgSrc:        product.ImgSrc,
			Rating:        product.Rating,
			CountComments: product.CountComments,
			Category: &gmodels.Category{
				Id:     product.Category.Id,
				Name:   product.Category.Name,
				Parent: product.Parent,
			},
		}
		gproducts[i] = gproduct
	}
	return &gen.CategoryResponse{Products: gproducts}, nil
}
