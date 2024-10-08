package http

import (
	"io"
	"log/slog"
	"net/http"
	"ozon_replic/internal/models/models"
	"ozon_replic/internal/pkg/auth/delivery/grpc/gen"
	"ozon_replic/internal/pkg/middleware/authmw"
	"ozon_replic/internal/pkg/middleware/logmw"
	"ozon_replic/internal/pkg/utils/logger/sl"
	resp "ozon_replic/internal/pkg/utils/responser"
	"time"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	client gen.AuthClient
	log    *slog.Logger
}

const customTimeFormat = "2006-01-02 15:04:05.999999999 -0700 UTC"

func NewAuthHandler(cl gen.AuthClient, log *slog.Logger) *AuthHandler {
	return &AuthHandler{
		client: cl,
		log:    log,
	}
}

// @Summary	SignIn(get)
// @Tags Auth
// @Description	Login to Account
// @Produce json
// @Success	200	{object} models.Profile "Profile"
// @Failure	400	{object} responser.response	"error messege"
// @Failure	429
// @Router	/api/auth/signin [get]
func foo() {

}

// @Summary	SignIn
// @Tags Auth
// @Description	Login to Account
// @Accept json
// @Produce json
// @Param input body models.SignInPayload true "SignInPayload"
// @Param X-CSRF-Token header string true "X-CSRF-Token"
// @Param Cookie header string true "Cookie"
// @Success	200	{object} models.Profile "Profile"
// @Failure	400	{object} responser.response	"error messege"
// @Failure	429
// @Router	/api/auth/signin [post]
func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	h.log = h.log.With(
		slog.String("op", sl.GFN()),
		slog.String("request_id", r.Header.Get(logmw.RequestIDCtx)),
	)

	body, err := io.ReadAll(r.Body)
	if resp.BodyErr(err, h.log, w) {
		return
	}
	h.log.Debug("request body decoded", slog.Any("request", r))
	defer r.Body.Close()

	userInfo := &models.SignInPayload{}
	err = userInfo.UnmarshalJSON(body)
	if err != nil {
		h.log.Error("failed to unmarshal request body", sl.Err(err))
		resp.JSONStatus(w, http.StatusTooManyRequests)

		return
	}

	profileAndCookie, err := h.client.SignIn(r.Context(), &gen.SignInRequest{
		Login:    userInfo.Login,
		Password: userInfo.Password,
	})
	if err != nil {
		h.log.Error("failed to signin", sl.Err(err))
		resp.JSON(w, http.StatusBadRequest, resp.Err("invalid login or password"))

		return
	}

	h.log.Debug("got profile", slog.Any("profile", profileAndCookie.Profile.Id))

	expiresTime, err := time.Parse(customTimeFormat, profileAndCookie.Expires)
	if err != nil {
		h.log.Error("failed to parse time from auth signin", sl.Err(err))
		resp.JSONStatus(w, http.StatusTooManyRequests)

		return
	}

	idUuid, err := uuid.FromString(profileAndCookie.Profile.Id)
	if err != nil {
		h.log.Error("failed to make uuid from string in uuid.FromString", sl.Err(err))
		resp.JSONStatus(w, http.StatusTooManyRequests)

		return
	}
	profile := &models.Profile{
		Id:          idUuid,
		Login:       profileAndCookie.Profile.Login,
		Description: profileAndCookie.Profile.Description,
		ImgSrc:      profileAndCookie.Profile.ImgSrc,
		Phone:       profileAndCookie.Profile.Phone,
	}

	http.SetCookie(w, authmw.MakeTokenCookie(profileAndCookie.Token, expiresTime))
	resp.JSON(w, http.StatusOK, profile)
}

// @Summary	SignUp
// @Tags Auth
// @Description	Create Account
// @Accept json
// @Produce json
// @Param input body models.SignUpPayload true "SignUpPayload"
// @Param X-CSRF-Token header string true "X-CSRF-Token"
// @Param Cookie header string true "Cookie"
// @Success	200 {object} models.Profile "Profile"
// @Failure	400	{object} responser.response	"error messege"
// @Failure	429
// @Router	/api/auth/signup [post]
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	h.log = h.log.With(
		slog.String("op", sl.GFN()),
		slog.String("request_id", r.Header.Get(logmw.RequestIDCtx)),
	)

	body, err := io.ReadAll(r.Body)
	if resp.BodyErr(err, h.log, w) {
		return
	}
	h.log.Debug("request body decoded", slog.Any("request", r))

	userInfo := &models.SignUpPayload{}
	err = userInfo.UnmarshalJSON(body)
	if err != nil {
		h.log.Error("failed to unmarshal request body", sl.Err(err))
		resp.JSONStatus(w, http.StatusTooManyRequests)

		return
	}

	profileAndCookie, err := h.client.SignUp(r.Context(), &gen.SignUpRequest{
		Login:    userInfo.Login,
		Password: userInfo.Password,
		Phone:    userInfo.Phone,
	})
	if err != nil {
		h.log.Error("failed to signup", sl.Err(err))
		resp.JSON(w, http.StatusBadRequest, resp.Err("invalid login or password"))

		return
	}

	h.log.Debug("got profile", slog.Any("profile", profileAndCookie.Profile.Id))

	expiresTime, err := time.Parse(customTimeFormat, profileAndCookie.Expires)
	if err != nil {
		h.log.Error("failed to parse time from auth signup", sl.Err(err))
		resp.JSONStatus(w, http.StatusTooManyRequests)

		return
	}

	idUuid, err := uuid.FromString(profileAndCookie.Profile.Id)
	if err != nil {
		h.log.Error("failed to make uuid from string in uuid.FromString", sl.Err(err))
		resp.JSONStatus(w, http.StatusTooManyRequests)

		return
	}
	profile := &models.Profile{
		Id:          idUuid,
		Login:       profileAndCookie.Profile.Login,
		Description: profileAndCookie.Profile.Description,
		ImgSrc:      profileAndCookie.Profile.ImgSrc,
		Phone:       profileAndCookie.Profile.Phone,
	}

	http.SetCookie(w, authmw.MakeTokenCookie(profileAndCookie.Token, expiresTime))
	resp.JSON(w, http.StatusOK, profile)

}

// @Summary	Logout
// @Tags Auth
// @Description	Logout from Account
// @Accept json
// @Produce json
// @Success	200
// @Failure	401
// @Router	/api/auth/logout [get]
func (h *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, authmw.MakeTokenCookie("", time.Now().UTC().AddDate(0, 0, -1)))
	h.log.Info("logout")
	resp.JSONStatus(w, http.StatusOK)
}

// @Summary	CheckAuth
// @Tags Auth
// @Description	Check user is logged in
// @Accept json
// @Produce json
// @Param X-CSRF-Token header string true "X-CSRF-Token"
// @Param Cookie header string true "Cookie"
// @Success	200	{object} models.Profile "Profile"
// @Failure	401
// @Failure	429
// @security AuthKey
// @Router	/api/auth/check_auth [get]
func (h *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
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

	h.log.Debug("check auth success", "id", id)

	profile, err := h.client.CheckAuth(r.Context(), &gen.CheckAuthRequst{
		ID: id.String(),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			h.log.Error("failed cast grpc error", sl.Err(err))
			resp.JSONStatus(w, http.StatusTooManyRequests)
			return
		}
		if st.Code() == codes.NotFound {
			h.log.Warn("profile not found", slog.Any("grpc status", st))
			resp.JSONStatus(w, http.StatusUnauthorized)
			return
		}

		h.log.Error("failed to CheckAuth", sl.Err(err))
		resp.JSONStatus(w, http.StatusTooManyRequests)

		return
	}

	idUuid, err := uuid.FromString(profile.Profile.Id)
	if err != nil {
		h.log.Error("failed to make uuid from string in uuid.FromString", sl.Err(err))
		resp.JSONStatus(w, http.StatusTooManyRequests)

		return
	}
	profileModel := models.Profile{
		Id:          idUuid,
		Login:       profile.Profile.Login,
		Description: profile.Profile.Description,
		ImgSrc:      profile.Profile.ImgSrc,
		Phone:       profile.Profile.Phone,
	}

	h.log.Debug("got profile", slog.Any("profile", profile.Profile.Id))
	resp.JSON(w, http.StatusOK, &profileModel)
}
