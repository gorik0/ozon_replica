package csrfmw

import (
	"log/slog"
	"net/http"
	"ozon_replic/internal/pkg/utils/jwter"
	"ozon_replic/internal/pkg/utils/logger/sl"
	resp "ozon_replic/internal/pkg/utils/responser"

	"github.com/gorilla/mux"
)

func contains(vals []string, s string) bool {
	for _, v := range vals {
		if v == s {
			return true
		}
	}

	return false
}

var safeMethods = []string{"GET", "HEAD", "OPTIONS", "TRACE"}

const HeaderName = "X-CSRF-Token"

func New(log *slog.Logger, jwtCORS jwter.JWTer) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler { // TODO: del
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if contains(safeMethods, r.Method) {
				token, _, err := jwtCORS.EncodeCSRFToken(r.UserAgent())
				if err != nil {
					log.Error("error happened in Auther.GenerateToken", sl.Err(err))
					resp.JSONStatus(w, http.StatusUnauthorized)

					return
				}
				w.Header().Set(HeaderName, token)

				return
			}

			token := r.Header.Get(HeaderName)

			println("TOKEN :::: ", token)
			if token == "" {
				log.Error("miss csrf jwt")
				resp.JSONStatus(w, http.StatusForbidden)

				return
			}
			//log.Debug("CSRF MW get csrf token", "token", token)

			UserAgent, err := jwtCORS.DecodeCSRFToken(token)

			if err != nil {
				log.Error("jws token is invalid csrf", sl.Err(err))
				resp.JSONStatus(w, http.StatusForbidden)

				return
			}
			if r.UserAgent() != UserAgent {
				log.Error("UserAgent from token does not match request UserAgent", "UserAgent", UserAgent)
				resp.JSONStatus(w, http.StatusForbidden)

				return
			}
			next.ServeHTTP(w, r)

		})
	}
}
