package api

import (
	"fmt"
	"net/http"
	"vehicledb/auth"
	"vehicledb/db"
)

var authCookieName = "auth"

func validateSession(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie(authCookieName)
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

type LoginRequest struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
}

var loginSchema = `{
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "email_address": {"type": "string"},
	"password": {"type": "string"}
  },
  "required": [
    "email_address",
	"password"
  ]
}`

func login(writer http.ResponseWriter, request *http.Request) {
	var loginRequest LoginRequest
	err := validateSchemaBuildModel(request, loginSchema, &loginRequest)
	if err != nil {
		renderError(writer, err)
		return
	}

	user, err := db.FindUserByEmailAddress(loginRequest.EmailAddress)
	if err != nil {
		renderError(writer, err)
		return
	}

	if user == nil {
		renderError(writer, fmt.Errorf("failed to find user with email address %s", loginRequest.EmailAddress))
		return
	}

	if !user.DoesPasswordMatch(loginRequest.Password) {
		renderError(writer, fmt.Errorf("password does not match"))
		return
	}

 	setAuthCookie(writer, user)

	renderJson(writer, user)
}

func logout(writer http.ResponseWriter, request *http.Request) {
	// wipe cookie
	removeAuthCookie(writer)

	writer.WriteHeader(http.StatusNoContent)
}
