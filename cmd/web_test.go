package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"vehicledb/api"
	"vehicledb/db"
	"vehicledb/graph"
)

var server *httptest.Server
var client *http.Client

func TestMain(m *testing.M) {
	schema, err := graph.GenerateSchema()
	if err != nil {
		log.Fatal("Failed to se up graph", err)
	}

	// setup global server
	mux := api.NewRouter(schema)
	server = httptest.NewServer(mux)
	defer server.Close()

	// set up global client
	client = &http.Client{}

	db.OpenDatabase("testing.sqlite")
	defer db.CloseDatabase()

	// call flag.Parse() here if TestMain uses flags
	result := m.Run()

	os.Exit(result)
}

func TestUserRoundTrip(t *testing.T) {
	// create user
	createUserRequest := api.CreateUserRequest{
		EmailAddress: "joe@djeebus.net",
		Password:     "Password1",
	}
	var createUserResponse db.User
	makeApiRequest(t, "POST", "/v1/users/", &createUserRequest, &createUserResponse)

	// get user
	var getUserResponse db.User
	makeApiRequest(t, "GET", fmt.Sprintf("/v1/users/%d/", createUserResponse.UserId), nil, &getUserResponse)
	if getUserResponse.UserId != createUserResponse.UserId {
		t.Fatalf("Returned user id != created user id")
	}
	if getUserResponse.EmailAddress != createUserResponse.EmailAddress {
		t.Fatalf("Returned email address != created email address")
	}

	// update user
	updateUserRequest := api.UpdateUserRequest{
		EmailAddress: "joe@eventray.com",
	}
	var updateUserResponse db.User
	makeApiRequest(t, "PATCH", fmt.Sprintf("/v1/users/%d/", createUserResponse.UserId), &updateUserRequest, &updateUserResponse)
	if updateUserResponse.UserId != createUserResponse.UserId {
		t.Fatalf("Returned user id != created user id")
	}
	if updateUserResponse.EmailAddress != updateUserRequest.EmailAddress {
		t.Fatalf("Returned email address != updated email address")
	}

	// delete user
	var deleteUserResponse db.User
	makeApiRequest(t, "DELETE", fmt.Sprintf("/v1/users/%d/", createUserResponse.UserId), nil, &deleteUserResponse)
}

func TestVehicleRoundTrip(t *testing.T) {
	// create vehicle
	createVehicleRequest := api.CreateVehicleRequest{
		Year:  2017,
		Make:  "Chevy",
		Model: "SS",
	}
	var createVehicleResponse db.Vehicle
	makeApiRequest(t, "POST", "/v1/vehicles/", &createVehicleRequest, &createVehicleResponse)
	if createVehicleRequest.Year != createVehicleResponse.Year {
		t.Fatalf("Failed to set year: %d != %d", createVehicleRequest.Year, createVehicleResponse.Year)
	}
	if createVehicleRequest.Make != createVehicleResponse.Make {
		t.Fatalf("Failed to set make: %s != %s", createVehicleRequest.Make, createVehicleResponse.Make)
	}
	if createVehicleRequest.Model != createVehicleResponse.Model {
		t.Fatalf("Failed to set model: %s != %s", createVehicleRequest.Model, createVehicleResponse.Model)
	}

	// get vehicle
	var getVehicleResponse db.Vehicle
	makeApiRequest(t, "GET", fmt.Sprintf("/v1/vehicles/%d", createVehicleResponse.VehicleID), nil, &getVehicleResponse)
	makeApiRequest(t, "POST", "/v1/vehicles/", &createVehicleRequest, &createVehicleResponse)
	if createVehicleResponse.VehicleID != getVehicleResponse.VehicleID {
		t.Fatalf("Failed to set year: %d != %d", createVehicleRequest.Year, getVehicleResponse.Year)
	}
	if createVehicleResponse.Year != getVehicleResponse.Year {
		t.Fatalf("Failed to set year: %d != %d", createVehicleRequest.Year, getVehicleResponse.Year)
	}
	if createVehicleResponse.Make != getVehicleResponse.Make {
		t.Fatalf("Failed to set make: %s != %s", createVehicleRequest.Make, getVehicleResponse.Make)
	}
	if createVehicleResponse.Model != getVehicleResponse.Model {
		t.Fatalf("Failed to set model: %s != %s", createVehicleRequest.Model, getVehicleResponse.Model)
	}

	// list vehicles

	// delete vehicle
	var deleteVehicleResponse db.Vehicle
	makeApiRequest(t, "DELETE", fmt.Sprintf("/v1/vehicles/%d", getVehicleResponse.VehicleID), nil, &deleteVehicleResponse)
	if createVehicleResponse.VehicleID != deleteVehicleResponse.VehicleID {
		t.Fatalf("Failed to set year: %d != %d", createVehicleRequest.Year, deleteVehicleResponse.Year)
	}
	if createVehicleResponse.Year != deleteVehicleResponse.Year {
		t.Fatalf("Failed to set year: %d != %d", createVehicleRequest.Year, deleteVehicleResponse.Year)
	}
	if createVehicleResponse.Make != deleteVehicleResponse.Make {
		t.Fatalf("Failed to set make: %s != %s", createVehicleRequest.Make, deleteVehicleResponse.Make)
	}
	if createVehicleResponse.Model != deleteVehicleResponse.Model {
		t.Fatalf("Failed to set model: %s != %s", createVehicleRequest.Model, deleteVehicleResponse.Model)
	}
}

func jsonToReader(model interface{}) (*bytes.Reader, error) {
	buf, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(buf)
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func makeApiRequest(t *testing.T, method string, path string, requestBody interface{}, responseBody interface{}) {
	var err error
	t.Logf("API request: %s %s\n", method, path)

	var requestBytes io.Reader
	if requestBody != nil {
		requestBytes, err = jsonToReader(requestBody)
		if err != nil {
			t.Fatalf("API request: request body: %v", err)
		}
	} else {
		requestBytes = nil
	}

	req, err := http.NewRequest(method, server.URL+path, requestBytes)
	if err != nil {
		t.Fatalf("API request: request: %v", err)
	}

	response, err := client.Do(req)
	if err != nil {
		t.Fatalf("API request: submit: %v", err)
	}

	t.Logf("API request: %s %s [%d]", method, path, response.StatusCode)

	if response.StatusCode >= 400 {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("API request failed [%d]: <failed to read body, %v>", response.StatusCode, err)
		}

		t.Fatalf("API request failed [%d]: %s", response.StatusCode, string(bodyBytes))
	}

	if responseBody != nil {
		body, err := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(body, responseBody)
		if err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}
	}
}
