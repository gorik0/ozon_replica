package http

import (
	"errors"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"ozon_replic/internal/models/models"
	"ozon_replic/internal/pkg/hub"
	"ozon_replic/internal/pkg/middleware/authmw"
	"ozon_replic/internal/pkg/middleware/logmw"
	"ozon_replic/internal/pkg/notifications"
	"ozon_replic/internal/pkg/notifications/repo"
	"ozon_replic/internal/pkg/utils/logger/sl"
	resp "ozon_replic/internal/pkg/utils/responser"

	uuid "github.com/satori/go.uuid"
)

var (
	upgrader = websocket.Upgrader{}
)

type NotificationsHandler struct {
	hub hub.HubInterface
	log *slog.Logger
	uc  notifications.NotificationsUsecase
}

func NewNotificationsHandler(hub hub.HubInterface, uc notifications.NotificationsUsecase, log *slog.Logger) *NotificationsHandler {
	return &NotificationsHandler{
		log: log,
		hub: hub,
		uc:  uc,
	}
}

func (h *NotificationsHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
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

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error("upgrade error:", err)
		return
	}
	h.log.Debug("connection upgraded: ", slog.Any("userID", userID))

	h.hub.AddClient(userID, connection)

	h.log.Debug("client disconnected: ", slog.Any("userID", userID))
}

// @Summary	GetDayNotifications
// @Tags Notifications
// @Description Get Day Notifications
// @Accept json
// @Produce json
// @Success	200	{object} models.MessageSlice "Recent today notifications"
// @Failure	401	"User unauthorized"
// @Failure	404	"Notifications not found"
// @Failure	429
// @Router	/api/notifications/get_recent [get]
func (h *NotificationsHandler) GetDayNotifications(w http.ResponseWriter, r *http.Request) {
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

	comment, err := h.uc.GetDayNotifications(r.Context(), userID)
	if err != nil {
		h.log.Error("failed in uc.GetDayNotifications", sl.Err(err))
		if errors.Is(err, repo.ErrNotificationsNotFound) {
			resp.JSONStatus(w, http.StatusNotFound)

			return
		}
		resp.JSONStatus(w, http.StatusTooManyRequests)

		return
	}

	h.log.Debug("uc.GetDayNotifications", "got notifications: ", len(comment))
	resp.JSON(w, http.StatusOK, (*models.MessageSlice)(&comment))
}
