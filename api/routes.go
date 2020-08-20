package api

import (
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"net/http"
)

func NewRouter(schema *graphql.Schema) *mux.Router {
	router := mux.NewRouter()

	// user routes
	userRoute := router.Path("/v1/users/{userId}/")
	AddMappedMethods(userRoute, map[string]http.HandlerFunc{
		"GET": getUser,
		"PATCH": updateUser,
		"DELETE": deleteUser,
	})

	usersRoute := router.Path("/v1/users/")
	usersRoute.Methods("POST").HandlerFunc(createUser)

	// api token routes
	tokensRoute := router.Path("/v1/tokens/")
	tokensRoute.Methods("GET").HandlerFunc(listTokens)
	tokensRoute.Methods("POST").HandlerFunc(createToken)

	tokenRoute := router.Path("/v1/tokens/{tokenId}")
	tokenRoute.Methods("GET").HandlerFunc(getToken)
	tokenRoute.Methods("PATCH").HandlerFunc(updateToken)
	tokenRoute.Methods("DELETE").HandlerFunc(deleteToken)

	// sessionsRoute routes
	sessionsRoute := router.Path("/v1/sessions/")
	sessionsRoute.Methods("GET").HandlerFunc(validateSession)
	sessionsRoute.Methods("POST").HandlerFunc(login)
	sessionsRoute.Methods("DELETE").HandlerFunc(logout)

	// vehicle routes
	vehiclesRoute := router.Path("/v1/vehicles/")
	vehiclesRoute.Methods("GET").HandlerFunc(listVehicles)
	vehiclesRoute.Methods("POST").HandlerFunc(createVehicle)

	vehicleRoute := router.Path("/v1/vehicles/{vehicleId}")
	vehicleRoute.Methods("GET").HandlerFunc(getVehicle)
	vehicleRoute.Methods("PATCH").HandlerFunc(updateVehicle)
	vehicleRoute.Methods("DELETE").HandlerFunc(deleteVehicle)

	// maintenance schedule routes
	scheduleRoute := router.Path("/v1/vehicles/{vehicleId}/schedule/")
	scheduleRoute.Methods("GET").HandlerFunc(listScheduledItems)
	scheduleRoute.Methods("POST").HandlerFunc(createScheduledItem)

	scheduleItemRoute := router.Path("/v1/vehicles/{vehicleId}/schedule/{scheduleItemId}")
	scheduleItemRoute.Methods("GET").HandlerFunc(getScheduledItem)
	scheduleItemRoute.Methods("PATCH").HandlerFunc(updateScheduledItem)
	scheduleItemRoute.Methods("DELETE").HandlerFunc(deleteScheduledItem)

	// ical or rss routes
	router.Path("/v1/vehicles/{vehicleId}/schedule.rss").Methods("GET").HandlerFunc(generateRssFeed)

	// graphql route
	graphQlHandler := handler.New(&handler.Config{
		Schema: schema,
		Pretty: true,
		GraphiQL: true,
	})
	router.Path("/v1/graphql").Handler(graphQlHandler)

	return router
}