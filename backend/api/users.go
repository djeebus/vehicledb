package api

import (
	"net/http"
	"vehicledb/auth"
	"vehicledb/db"
)

type CreateUserRequest struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
}

var createUserSchema = `{
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

func createUser(w http.ResponseWriter, request *http.Request) {
	var createUserRequest CreateUserRequest
	err := validateSchemaBuildModel(request, createUserSchema, &createUserRequest)
	if err != nil {
		renderError(w, err)
		return
	}

	user, err := db.CreateUser(createUserRequest.EmailAddress, createUserRequest.Password)
	if err != nil {
		renderError(w, err)
		return
	}

	token, err := auth.CreateToken(user)
	if err != nil {
		renderError(w, err)
		return
	}

	cookie := http.Cookie{
		Path: "/",
		Name: "auth",
		Value: token,
	}
	http.SetCookie(w, &cookie)

	renderJson(w, user)
}

func getUser(claimsUser *auth.ClaimsUser, writer http.ResponseWriter, request *http.Request) {
	user, err := db.GetUser(claimsUser.UserID)
	if err != nil {
		renderError(writer, err)
		return
	}

	renderJson(writer, user)
}

type UpdateUserRequest struct {
	EmailAddress string `json:"email_address"`
}

var updateUserSchema = `{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"email_address": {"type": "string"}
	}
}`

func updateUser(claimsUser *auth.ClaimsUser, w http.ResponseWriter, request *http.Request) {
	var updateUserRequest UpdateUserRequest
	err := validateSchemaBuildModel(request, updateUserSchema, &updateUserRequest)
	if err != nil {
		renderError(w, err)
		return
	}

	user, err := db.UpdateUser(claimsUser.UserID, updateUserRequest.EmailAddress)
	if err != nil {
		renderError(w, err)
		return
	}

	renderJson(w, user)
}

func deleteUser(claimsUser *auth.ClaimsUser, w http.ResponseWriter, request *http.Request) {
	user, err := db.GetUser(claimsUser.UserID)
	if err != nil {
		renderError(w, err)
		return
	}

	err = db.DeleteUser(claimsUser.UserID)
	if err != nil {
		renderError(w, err)
		return
	}

	renderJson(w, user)
}
