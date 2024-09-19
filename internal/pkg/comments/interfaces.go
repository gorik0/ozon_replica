package comments

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"ozon_replic/internal/models/models"
)

//go:generate mockgen -source interfaces.go -destination ./mocks/comments_mock.go -package mock

type CommentsUsecase interface {
	CreateComment(context.Context, models.CommentPayload) (models.Comment, error)
	GetProductComments(context.Context, uuid.UUID) ([]models.Comment, error)
}

type CommentsRepo interface {
	ReadCountOfCommentsToProduct(context.Context, uuid.UUID, uuid.UUID) (int, models.Comment, error)
	MakeComment(context.Context, models.CommentPayload) (models.Comment, error)
	ReadProductComments(context.Context, uuid.UUID) ([]models.Comment, error)
}
