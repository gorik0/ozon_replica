package grpc

import (
	"context"
	"errors"
	"log/slog"
	"ozon_replic/internal/models/models"
	"ozon_replic/internal/pkg/auth"
	"ozon_replic/internal/pkg/auth/delivery/grpc/gen"
	"ozon_replic/internal/pkg/middleware/metricsmw"
	"ozon_replic/internal/pkg/profile"
	"ozon_replic/internal/pkg/utils/logger/sl"
	"ozon_replic/proto/gmodels"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	uuid "github.com/satori/go.uuid"
)

//go:generate mockgen -source=./gen/auth_grpc.pb.go -destination=../../mocks/auth_grpc.go -package=mock

type GrpcAuthHandler struct {
	uc  auth.AuthUsecase
	log *slog.Logger

	gen.AuthServer
}

func NewGrpcAuthHandler(uc gen.AuthClient, log *slog.Logger) *GrpcAuthHandler {
	return &GrpcAuthHandler{
		uc:  uc,
		log: log,
	}
}

func (h *GrpcAuthHandler) SignIn(ctx context.Context, in *gen.SignInRequest) (*gen.SignInResponse, error) {
	h.log = h.log.With(
		slog.String("op", sl.GFN()),
	)

	userSignIn := models.SignInPayload{
		Login:    in.Login,
		Password: in.Password,
	}

	profile, token, expires, err := h.uc.SignIn(ctx, &userSignIn)
	if err != nil {
		h.log.Error("failed in uc.SignIn", sl.Err(err))

		return &gen.SignInResponse{}, metricsmw.ServerError
	}
	h.log.Info("got profile", slog.Any("profile", profile.Id))

	return &gen.SignInResponse{
		Profile: &gmodels.Profile{
			Id:          profile.Id.String(),
			Login:       profile.Login,
			Description: profile.Description,
			ImgSrc:      profile.ImgSrc,
			Phone:       profile.Phone,
		},
		Token:   token,
		Expires: expires.String(),
	}, nil
}

func (h *GrpcAuthHandler) SignUp(ctx context.Context, in *gen.SignUpRequest) (*gen.SignUpResponse, error) {
	h.log = h.log.With(
		slog.String("op", sl.GFN()),
	)

	userSignUp := models.SignUpPayload{
		Login:    in.Login,
		Password: in.Password,
		Phone:    in.Phone,
	}

	profile, token, expires, err := h.uc.SignUp(ctx, &userSignUp)
	if err != nil {
		h.log.Error("failed in uc.SignUp", sl.Err(err))
		return &gen.SignUpResponse{}, metricsmw.ServerError
	}
	h.log.Info("got profile", slog.Any("profile", profile.Id))

	return &gen.SignUpResponse{
		Profile: &gmodels.Profile{
			Id:          profile.Id.String(),
			Login:       profile.Login,
			Description: profile.Description,
			ImgSrc:      profile.ImgSrc,
			Phone:       profile.Phone,
		},

		Token:   token,
		Expires: expires.String(),
	}, nil
}

func (h *GrpcAuthHandler) CheckAuth(ctx context.Context, in *gen.CheckAuthRequst) (*gen.CheckAuthResponse, error) {
	h.log = h.log.With(
		slog.String("op", sl.GFN()),
	)

	userID, err := uuid.FromString(in.ID)
	if err != nil {
		h.log.Error("failed to get uuid from string", sl.Err(err))
		return &gen.CheckAuthResponse{}, metricsmw.ClientError
	}

	profil, err := h.uc.CheckAuth(ctx, userID)
	if err != nil {
		if errors.Is(err, profile.ErrProfileNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		h.log.Error("failed in uc.CheckAuth", sl.Err(err))
		return &gen.CheckAuthResponse{}, metricsmw.ServerError
	}
	h.log.Info("got profile", slog.Any("profile", profil.Id))

	return &gen.CheckAuthResponse{
		Profile: &gmodels.Profile{
			Id:          profil.Id.String(),
			Login:       profil.Login,
			Description: profil.Description,
			ImgSrc:      profil.ImgSrc,
			Phone:       profil.Phone,
		},
	}, nil
}
