package http

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"ozon_replic/internal/models/models"
	"ozon_replic/internal/pkg/comments"
	"ozon_replic/internal/pkg/comments/repo"
	"ozon_replic/internal/pkg/comments/usecase"
	"ozon_replic/internal/pkg/middleware/authmw"
	"ozon_replic/internal/pkg/middleware/logmw"
	"ozon_replic/internal/pkg/utils/logger/sl"
	resp "ozon_replic/internal/pkg/utils/responser"

	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
)

type CommentsHandler struct {
	log *slog.Logger
	uc  comments.CommentsUsecase
}

func NewCommentsHandler(log *slog.Logger, uc comments.CommentsUsecase) *CommentsHandler {
	return &CommentsHandler{
		log: log,
		uc:  uc,
	}
}

// @Summary	CreateComment
// @Tags Comments
// @Description	Create Comment to product
// @Accept json
// @Produce json
// @Param input body models.CommentPayload true "cart info"
// @Success	200 {object} models.Comment "Comment created"
// @Failure	401	"User unauthorized"
// @Failure	406 "Comment is invalid"
// @Failure	413 {object} models.Comment "User already has commented this product"
// @Failure	429
// @Router	/api/comments/create [post]
func (h *CommentsHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	h.log = h.log.With(
		slog.String("op", sl.GFN()),
		slog.String("request_id", r.Header.Get(logmw.RequestIDCtx)),
	)

	userID, ok := r.Context().Value(authmw.AccessTokenCookieName).(uuid.UUID)
	if !ok {
		h.log.Error("failed cast uuid from context value")
		resp.JSONStatus(w, http.StatusUnauthorized)

		return
	}

	body, err := io.ReadAll(r.Body)
	if resp.BodyErr(err, h.log, w) {
		return
	}
	defer r.Body.Close()
	h.log.Debug("request body decoded", slog.Any("request", r))

	commentPayload := models.CommentPayload{
		UserID: userID,
	}
	err = commentPayload.UnmarshalJSON(body)
	if err != nil {
		h.log.Error("failed to unmarshal request body", sl.Err(err))
		resp.JSONStatus(w, http.StatusTooManyRequests)

		return
	}

	if err = validator.New().Struct(commentPayload); err != nil {
		h.log.Error("comment is invalid", sl.Err(err))
		resp.JSONStatus(w, http.StatusNotAcceptable)

		return
	}

	commentPayload.Sanitize()

	comment, err := h.uc.CreateComment(r.Context(), commentPayload)
	if err != nil {
		h.log.Error("failed in uc.CreateComment", sl.Err(err))
		if errors.Is(err, usecase.ErrManyCommentsToProduct) {
			resp.JSON(w, http.StatusRequestEntityTooLarge, &comment)

			return
		}
		resp.JSONStatus(w, http.StatusTooManyRequests)

		return
	}

	h.log.Debug("uc.CreateComment", "comment created", comment)
	resp.JSON(w, http.StatusOK, &comment)
}

// @Summary	GetProductComments
// @Tags Comments
// @Description	Get product comments
// @Accept json
// @Produce json
// @Param product query string true "Product ID"
// @Success	200	{object} []models.Comment "Comments array"
// @Failure	400	{object} responser.response	"error message"
// @Failure	429
// @Router	/api/comments/get_all [get]
func (h *CommentsHandler) GetProductComments(w http.ResponseWriter, r *http.Request) {
	h.log = h.log.With(
		slog.String("op", sl.GFN()),
		slog.String("request_id", r.Header.Get(logmw.RequestIDCtx)),
	)
	productIDString := r.URL.Query().Get("product")
	productID, err := uuid.FromString(productIDString)
	if err != nil || len(productID) == 0 {
		h.log.Error("productID is invalid", sl.Err(err))
		resp.JSON(w, http.StatusBadRequest, resp.Err("invalid request"))

		return
	}

	comments, err := h.uc.GetProductComments(r.Context(), productID)
	if err != nil {
		h.log.Error("failed to get comments", sl.Err(err))
		if errors.Is(err, repo.ErrCommentNotFound) {
			resp.JSONStatus(w, http.StatusNotFound)

			return
		}
		resp.JSONStatus(w, http.StatusTooManyRequests)

		return
	}

	h.log.Debug("got comments", "len", len(comments))
	resp.JSON(w, http.StatusOK, (*models.CommentSlice)(&comments))
}
