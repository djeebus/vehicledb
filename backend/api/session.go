package api

import (
	"fmt"
	"net/http"
	"vehicledb/auth"
)

func validateSession(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("auth")
	if err == http.ErrNoCookie {
		writer.WriteHeader(401)
		renderJson(writer, map[string]string {"code": "unauthorized"})
		return
	}

	if err != nil {
		renderError(writer, fmt.Errorf("failed to read auth cookie: %v", err))
		return
	}

	user, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		writer.WriteHeader(403)
		renderJson(writer, map[string]string {"code": "forbidden"})
		return
	}

	renderJson(writer, map[string]string{
		"email_address": user.EmailAddress,
		"user_id": user.UserID.String(),
	})
}

func login(writer http.ResponseWriter, request *http.Request) {

}

func logout(writer http.ResponseWriter, request *http.Request) {

}
