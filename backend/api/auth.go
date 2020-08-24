package api

import (
	"errors"
	"net/http"
	"vehicledb/auth"
)

type UserHandlerFunc func(user *auth.ClaimsUser, w http.ResponseWriter, r *http.Request)

func RequireAuth(f UserHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth")
		if err != nil {
			w.WriteHeader(401)
			renderJson(w, "must pass cookie")
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
			return
		}

		f(user, w, r)
	}
}
