package api

import (
	"fmt"
	"github.com/gorilla/mux"
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

	header := w.Header()
	header.Add("Set-Cookie", fmt.Sprintf("auth=%s", token))

	renderJson(w, user)
}

func getUser(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	userId, err := db.ParseRowID(vars["userId"])
	if err != nil {
		renderError(writer, err)
		return
	}

	user, err := db.GetUser(userId)
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

func updateUser(w http.ResponseWriter, request *http.Request) {
	var updateUserRequest UpdateUserRequest
	err := validateSchemaBuildModel(request, updateUserSchema, &updateUserRequest)
	if err != nil {
		renderError(w, err)
		return
	}

	vars := mux.Vars(request)
	userId, err := db.ParseRowID(vars["userId"])
	if err != nil {
		renderError(w, err)
		return
	}

	user, err := db.UpdateUser(userId, updateUserRequest.EmailAddress)
	if err != nil {
		renderError(w, err)
		return
	}

	renderJson(w, user)
}

func deleteUser(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	userId, err := db.ParseRowID(vars["userId"])
	if err != nil {
		renderError(w, err)
		return
	}

	user, err := db.GetUser(userId)
	if err != nil {
		renderError(w, err)
		return
	}

	err = db.DeleteUser(userId)
	if err != nil {
		renderError(w, err)
		return
	}

	renderJson(w, user)
}
