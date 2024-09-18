package http

import (
	"net/http"
	"ozon_replic/internal/pkg/category"
	"ozon_replic/internal/pkg/middleware/logmw"
	"ozon_replic/internal/pkg/utils/logger/sl"
	resp "ozon_replic/internal/pkg/utils/responser"

	"log/slog"
)

type CategoryHandler struct {
	log *slog.Logger
	uc  category.CategoryUsecase
}

func NewCategoryHandler(log *slog.Logger, uc category.CategoryUsecase) CategoryHandler {
	return CategoryHandler{
		log: log,
		uc:  uc,
	}
}

// @Summary	Category
// @Tags Category
// @Description	Get category tree
// @Accept json
// @Produce json
// @Success	200	{object} []models.Category "Category tree"
// @Failure	429
// @Router	/api/category/get_all [get]
func (h *CategoryHandler) Categories(w http.ResponseWriter, r *http.Request) {
	h.log = h.log.With(
		slog.String("op", sl.GFN()),
		slog.String("request_id", r.Header.Get(logmw.RequestIDCtx)),
	)

	categories, err := h.uc.Categories(r.Context())
	if err != nil {
		h.log.Error("failed to get categories", sl.Err(err))
		resp.JSONStatus(w, http.StatusTooManyRequests)

		return
	}

	h.log.Debug("got category", "len", len(categories))
	resp.JSON(w, http.StatusOK, &categories)
}
