package http

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"io"
	"log"
	"log/slog"
	"net/http"
	"ozon_replic/internal/models/models"
	"ozon_replic/internal/pkg/middleware/authmw"
	"ozon_replic/internal/pkg/middleware/logmw"
	"ozon_replic/internal/pkg/profile"
	"ozon_replic/internal/pkg/utils/logger/sl"
	resp "ozon_replic/internal/pkg/utils/responser"
)

const maxRequestBodySize = 1024 * 1024 * 5 // 5 MB

type ProfileHandler struct {
	log *slog.Logger
	uc  profile.ProfileUsecase
}

func NewProfileHandler(log *slog.Logger, uc profile.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{
		log: log,
		uc:  uc,
	}
}

// @Summary	GetProfile
// @Tags Profile
// @Description	Get profile by ID
// @Accept json
// @Produce json
// @Param id path string true "Profile UUID"
// @Success	200	{object} models.Profile "Profile"
// @Failure	400	{object} responser.response	"error messege"
// @Failure	429
// @Router	/api/profile/{id} [get]
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	h.log = h.log.With(
		slog.String("op", sl.GFN()),
		slog.String("request_id", r.Header.Get(logmw.RequestIDCtx)),
	)

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok || idStr == "" {
		h.log.Error("id is empty")
		resp.JSON(w, http.StatusBadRequest, resp.Err("invalid request"))

		return
	}
	idProfile, err := uuid.FromString(idStr)
	if err != nil {
		h.log.Error("id is invalid", sl.Err(err))
		resp.JSON(w, http.StatusBadRequest, resp.Err("invalid request"))

		return
	}

	profile, err := h.uc.GetProfile(r.Context(), idProfile)

	if err != nil {
		h.log.Error("failed to signup", sl.Err(err))
		resp.JSON(w, http.StatusBadRequest, resp.Err("invalid uuid"))

		return
	}

	h.log.Debug("got profile", slog.Any("profile", profile.Id))
	resp.JSON(w, http.StatusOK, profile)
}

// @Summary	UpdateProfileData
// @Tags Profile
// @Description	Update profile data
// @Accept json
// @Produce json
// @Param input body models.UpdateProfileDataPayload true "UpdateProfileDataPayload"
// @Success	200	{object} models.Profile "Profile"
// @Failure	400	{object} responser.response	"error messege"
// @Failure	401
// @Failure	429
// @Router	/api/profile/update-data [post]
func (h *ProfileHandler) UpdateProfileData(w http.ResponseWriter, r *http.Request) {
	h.log = h.log.With(
		slog.String("op", sl.GFN()),
		slog.String("request_id", r.Header.Get(logmw.RequestIDCtx)),
	)

	id, ok := r.Context().Value(authmw.AccessTokenCookieName).(uuid.UUID)
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
	h.log.Debug("got file from r.Body", slog.Any("request", r))

	payload := &models.UpdateProfileDataPayload{}
	err = payload.UnmarshalJSON(body)
	if err != nil {
		h.log.Error("failed to unmarshal request body", sl.Err(err))
		resp.JSONStatus(w, http.StatusTooManyRequests)

		return
	}

	profileInfo, err := h.uc.UpdateData(r.Context(), id, payload)
	if err != nil {
		h.log.Warn("failed to update profile data", sl.Err(err))
		var validationErrors validator.ValidationErrors

		switch {
		case errors.Is(err, profile.ErrBadUpdateData):
			resp.JSONStatus(w, http.StatusBadRequest)
		case errors.As(err, &validationErrors):
			resp.JSON(w, http.StatusBadRequest, resp.Err(err.Error()))
		default:
			resp.JSONStatus(w, http.StatusTooManyRequests)
		}
		return
	}

	h.log.Info("updated profile info")
	resp.JSON(w, http.StatusOK, profileInfo)
}

// @Summary	UpdatePhoto
// @Tags Profile
// @Description	Update profile photo
// @Accept json
// @Produce json
// @Param input body byte true "photo"
// @Success	200	{object} models.Profile "Profile"
// @Failure	401
// @Failure 413
// @Failure	429
// @Router	/api/profile/update-photo [post]
func (h *ProfileHandler) UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	h.log = h.log.With(
		slog.String("op", sl.GFN()),
		slog.String("request_id", r.Header.Get(logmw.RequestIDCtx)),
	)
	ID, ok := r.Context().Value(authmw.AccessTokenCookieName).(uuid.UUID)
	if !ok {
		h.log.Error("failed cast uuid from context value")
		resp.JSONStatus(w, http.StatusUnauthorized)

		return
	}

	limitedReader := http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	defer r.Body.Close()

	bodyContent, err := io.ReadAll(limitedReader)
	println("BODY ::::")
	println("BODY ::::")
	println("BODY ::::")
	log.Printf("%x", bodyContent[:512])
	println(" ::::")
	println(" ::::")
	fileFormat := http.DetectContentType(bodyContent)
	h.log.Debug("got []byte file", slog.Any("request", r))

	if err != nil && !errors.Is(err, io.EOF) {
		if errors.As(err, new(*http.MaxBytesError)) {
			resp.JSONStatus(w, http.StatusRequestEntityTooLarge)
		} else {
			resp.JSONStatus(w, http.StatusTooManyRequests)
		}

		return
	}

	h.log.Info("UpdatePhoto", "ID", ID)
	profileInfo, err := h.uc.UpdatePhoto(r.Context(), ID, bodyContent, fileFormat)
	if err != nil {
		h.log.Error("failed in uc.UpdatePhoto", sl.Err(err))
		resp.JSONStatus(w, http.StatusTooManyRequests)

		return
	}

	h.log.Info("updated profile success")
	resp.JSON(w, http.StatusOK, profileInfo)
}
