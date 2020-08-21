package api

import (
	"errors"
	"net/http"
	"vehicledb/auth"
)

func RequireAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth")
		if err != nil {
			renderError(w, err)
			return
		}

		if cookie == nil {
			w.WriteHeader(401)
			renderJson(w, errors.New("must be authenticated"))
			return
		}

		user, err := auth.ValidateToken(cookie.Value)
		if err != nil {
			renderError(w, err)
			return
		}

		if user == nil {
			w.WriteHeader(403)
			renderJson(w, errors.New("invalid jwt"))
		}

		h.ServeHTTP(w, r)
	}
}
