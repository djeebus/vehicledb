package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
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
	// setup database
	db.OpenDatabase("testing.sqlite")
	defer func() { logError(db.CloseDatabase()) }()

	// setup schema
	schema, err := graph.GenerateSchema()
	if err != nil {
		log.Fatal("Failed to se up graph", err)
	}

	// setup global server
	mux := api.NewHandler(schema)
	server = httptest.NewServer(mux)
	defer server.Close()

	// set up global client
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		log.Fatal("Failed to set up cookie jar", err)
	}
	client = &http.Client{Jar: jar}

	// run the test
	result := m.Run()

	// quit
	os.Exit(result)
}

func logError(err error) {
	if err != nil {
		log.Println("Error closing db: ", err)
	}
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
	// create user
	createUserRequest := api.CreateUserRequest{
		EmailAddress: "joe@djeebus.net",
		Password:     "Password1",
	}
	var createUserResponse db.User
	makeApiRequest(t, "POST", "/v1/users/", &createUserRequest, &createUserResponse)

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

	// update vehicle
	updateVehicleRequest := api.UpdateVehicleRequest{
		Year: &db.NullYear{Year: 2011, Valid: true},
		Make: &db.NullString{String: "BMW", Valid: true},
		Model: &db.NullString{String: "550i", Valid: true},
	}
	var updateVehicleResponse db.Vehicle
	makeApiRequest(t, "PATCH", fmt.Sprintf("/v1/vehicles/%d", getVehicleResponse.VehicleID), &updateVehicleRequest, &updateVehicleResponse)
	if updateVehicleResponse.Year != updateVehicleRequest.Year.Year {
		t.Fatalf("Failed to set year: %d != %d", updateVehicleResponse.Year, updateVehicleRequest.Year.Year)
	}
	if updateVehicleResponse.Make != updateVehicleRequest.Make.String {
		t.Fatalf("Failed to set make: %s != %s", updateVehicleResponse.Make, updateVehicleRequest.Make.String)
	}
	if updateVehicleResponse.Model != updateVehicleRequest.Model.String {
		t.Fatalf("Failed to set model: %s != %s", updateVehicleResponse.Model, updateVehicleRequest.Model.String)
	}

	// delete vehicle
	var deleteVehicleResponse db.Vehicle
	makeApiRequest(t, "DELETE", fmt.Sprintf("/v1/vehicles/%d", getVehicleResponse.VehicleID), nil, &deleteVehicleResponse)
	if createVehicleResponse.VehicleID != deleteVehicleResponse.VehicleID {
		t.Fatalf("Failed to set year: %d != %d", createVehicleRequest.Year, deleteVehicleResponse.Year)
	}
	if updateVehicleResponse.Year != deleteVehicleResponse.Year {
		t.Fatalf("Failed to set year: %d != %d", createVehicleRequest.Year, deleteVehicleResponse.Year)
	}
	if updateVehicleResponse.Make != deleteVehicleResponse.Make {
		t.Fatalf("Failed to set make: %s != %s", createVehicleRequest.Make, deleteVehicleResponse.Make)
	}
	if updateVehicleResponse.Model != deleteVehicleResponse.Model {
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

	cookies := client.Jar.Cookies(req.URL)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
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
