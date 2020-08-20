package api

import "net/http"

type CreateVehicleRequest struct {
	Year uint16
	Make string
	Model string
}

func listVehicles(writer http.ResponseWriter, request *http.Request) {

}

func createVehicle(writer http.ResponseWriter, request *http.Request) {

}

func getVehicle(writer http.ResponseWriter, request *http.Request) {

}

func updateVehicle(writer http.ResponseWriter, request *http.Request) {

}

func deleteVehicle(writer http.ResponseWriter, request *http.Request) {

}
