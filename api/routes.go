package api

import (
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"net/http"
)

func NewHandler(schema *graphql.Schema) http.Handler {
	router := mux.NewRouter()

	// user routes
	userRoute := router.Path("/v1/users/{userId}/")
	AddMappedMethods(userRoute, map[string]http.HandlerFunc{
		"GET":    getUser,
		"PATCH":  updateUser,
		"DELETE": deleteUser,
	})

	usersRoute := router.Path("/v1/users/")
	AddMappedMethods(usersRoute, map[string]http.HandlerFunc{
		"POST": createUser,
	})

	// api token routes
	tokensRoute := router.Path("/v1/tokens/")
	AddMappedMethods(tokensRoute, map[string]http.HandlerFunc{
		"GET":  listTokens,
		"POST": createToken,
	})

	tokenRoute := router.Path("/v1/tokens/{tokenId}")
	AddMappedMethods(tokenRoute, map[string]http.HandlerFunc{
		"GET":    getToken,
		"PATCH":  updateToken,
		"DELETE": deleteToken,
	})

	// sessionsRoute routes
	sessionsRoute := router.Path("/v1/sessions/")
	AddMappedMethods(sessionsRoute, map[string]http.HandlerFunc{
		"GET":    validateSession,
		"POST":   login,
		"DELETE": logout,
	})

	// vehicle routes
	vehiclesRoute := router.Path("/v1/vehicles/")
	AddMappedMethods(
		vehiclesRoute,
		map[string]http.HandlerFunc{
			"GET":  listVehicles,
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

	return router
}
