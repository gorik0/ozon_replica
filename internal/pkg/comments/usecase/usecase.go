package usecase

import (
	"context"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"ozon_replic/internal/models/models"
	"ozon_replic/internal/pkg/comments"
	"ozon_replic/internal/pkg/comments/repo"
)

var (
	ErrManyCommentsToProduct = errors.New("many comments to one product")
)

type CommentsUsecase struct {
	repo comments.CommentsRepo
}

func NewCommentsUsecase(repoComments comments.CommentsRepo) *CommentsUsecase {
	return &CommentsUsecase{
		repo: repoComments,
	}
}

func (uc *CommentsUsecase) CreateComment(ctx context.Context, commentPayload models.CommentPayload) (
	models.Comment, error) {
	count, comment, err := uc.repo.ReadCountOfCommentsToProduct(ctx, commentPayload.UserID, commentPayload.ProductID)
	if !errors.Is(err, repo.ErrCommentNotFound) && err != nil {
		err = fmt.Errorf("error happened in repo.MakeComment: %w", err)

		return models.Comment{}, err
	}
	if count > 0 {
		return comment, ErrManyCommentsToProduct
	}
	comment, err = uc.repo.MakeComment(ctx, commentPayload)
	if err != nil {
		err = fmt.Errorf("error happened in repo.MakeComment: %w", err)

		return models.Comment{}, err
	}

	return comment, nil
}

func (uc *CommentsUsecase) GetProductComments(ctx context.Context, productID uuid.UUID) ([]models.Comment, error) {
	comments, err := uc.repo.ReadProductComments(ctx, productID)
	if err != nil {
		if errors.Is(err, repo.ErrCommentNotFound) {
			return []models.Comment{}, err
		}
		err = fmt.Errorf("error happened in repo.ReadProductComments: %w", err)

		return []models.Comment{}, err
	}

	return comments, nil
}
