package api

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"net/http"
)

func NewHandler(schema *graphql.Schema, allowedOrigins []string) http.Handler {
	router := mux.NewRouter()

	// user routes
	AddMappedMethods(
		router.Path("/v1/users/me"),
		map[string]http.HandlerFunc{
			"GET":    RequireAuth(getUser),
			"PATCH":  RequireAuth(updateUser),
			"DELETE": RequireAuth(deleteUser),
		})

	AddMappedMethods(
		router.Path("/v1/users/"),
		map[string]http.HandlerFunc{
			"POST": createUser,
		})

	// api token routes
	AddMappedMethods(
		router.Path("/v1/tokens/"),
		map[string]http.HandlerFunc{
			"GET":  listTokens,
			"POST": createToken,
		})

	AddMappedMethods(
		router.Path("/v1/tokens/{tokenId}"),
		map[string]http.HandlerFunc{
			"GET":    getToken,
			"PATCH":  updateToken,
			"DELETE": deleteToken,
		})

	// sessionsRoute routes
	AddMappedMethods(
		router.Path("/v1/session"),
		map[string]http.HandlerFunc{
			"GET":    validateSession,
			"POST":   login,
			"DELETE": logout,
		})

	// vehicle routes
	AddMappedMethods(
		router.Path("/v1/vehicles/"),
		map[string]http.HandlerFunc{
			"GET":  RequireAuth(listVehicles),
			"POST": RequireAuth(createVehicle),
		},
	)

	vehicleRoute := router.Path("/v1/vehicles/{vehicleId}")
	AddMappedMethods(vehicleRoute, map[string]http.HandlerFunc{
		"GET":    getVehicle,
		"PATCH":  updateVehicle,
		"DELETE": deleteVehicle,
	})

	// maintenance schedule routes
	scheduleRoute := router.Path("/v1/vehicles/{vehicleId}/schedule/")
	AddMappedMethods(scheduleRoute, map[string]http.HandlerFunc{
		"GET":  listScheduledItems,
		"POST": createScheduledItem,
	})

	scheduleItemRoute := router.Path("/v1/vehicles/{vehicleId}/schedule/{scheduleItemId}")
	AddMappedMethods(scheduleItemRoute, map[string]http.HandlerFunc{
		"GET":    getScheduledItem,
		"PATCH":  updateScheduledItem,
		"DELETE": deleteScheduledItem,
	})

	// ical or rss routes
	router.Path("/v1/vehicles/{vehicleId}/schedule.rss").Methods("GET").HandlerFunc(generateRssFeed)

	// graphql route
	graphQlHandler := handler.New(&handler.Config{
		Schema:   schema,
		Pretty:   true,
		GraphiQL: true,
	})
	router.Path("/v1/graphql").Handler(graphQlHandler)

	// cors
	corsWrapper := handlers.CORS(
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "PUT", "DELETE"}),
		handlers.AllowCredentials(),
	)
	corsWrapped := corsWrapper(router)

	return corsWrapped
}
