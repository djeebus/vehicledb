package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"vehicledb/db"
)

func listVehicles(writer http.ResponseWriter, request *http.Request) {

}

type CreateVehicleRequest struct {
	Year  db.Year `json:"year"`
	Make  string  `json:"make"`
	Model string  `json:"model"`
}

var createVehicleSchema = `{
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "year": {"type": "integer"},
	"make": {"type": "string"},
	"model": {"type": "string"}
  },
  "required": [
    "year",
	"make",
	"model"
  ]
}`

func createVehicle(writer http.ResponseWriter, request *http.Request) {
	var createVehicleRequest CreateVehicleRequest
	err := validateSchemaBuildModel(request, createVehicleSchema, &createVehicleRequest)
	if err != nil {
		renderError(writer, err)
		return
	}

	user, err := db.CreateVehicle(createVehicleRequest.Year, createVehicleRequest.Make, createVehicleRequest.Model)
	if err != nil {
		renderError(writer, err)
	} else {
		renderJson(writer, user)
	}
}

func getVehicle(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	vehicleId, err := db.ParseRowID(vars["vehicleId"])
	if err != nil {
		renderError(writer, err)
		return
	}

	vehicle, err := db.GetVehicle(vehicleId)
	if err != nil {
		renderError(writer, err)
		return
	}

	renderJson(writer, vehicle)
}

var updateVehicleSchema = `{
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "year": {"type": "integer"},
	"make": {"type": "string"},
	"model": {"type": "string"}
  }
}`

type UpdateVehicleRequest struct {
	Year  *db.NullYear   `json:"year"`
	Make  *db.NullString `json:"make"`
	Model *db.NullString `json:"model"`
}

func updateVehicle(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	vehicleId, err := db.ParseRowID(vars["vehicleId"])
	if err != nil {
		renderError(writer, err)
		return
	}

	var updateVehicleRequest UpdateVehicleRequest
	err = validateSchemaBuildModel(request, updateVehicleSchema, &updateVehicleRequest)
	if err != nil {
		renderError(writer, err)
		return
	}

	err = db.UpdateVehicle(vehicleId, updateVehicleRequest.Year, updateVehicleRequest.Make, updateVehicleRequest.Model)
	if err != nil {
		renderError(writer, err)
		return
	}

	vehicle, err := db.GetVehicle(vehicleId)

	renderJson(writer, vehicle)
}

func deleteVehicle(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	vehicleId, err := db.ParseRowID(vars["vehicleId"])
	if err != nil {
		renderError(writer, err)
		return
	}

	vehicle, err := db.GetVehicle(vehicleId)
	if err != nil {
		renderError(writer, err)
		return
	}

	err = db.DeleteVehicle(vehicleId)
	if err != nil {
		renderError(writer, err)
		return
	}

	renderJson(writer, vehicle)
}
