package cmd

import (
	"log"
	"net/http"
	"time"
	"vehicledb/api"
	"vehicledb/graph"
)

func RunApiServer() {
	schema, err := graph.GenerateSchema()
	if err != nil {
		log.Fatal("Failed to generate schema", err)
	}
	router := api.NewRouter(schema)

	addr := "127.0.0.1:8000"

	srv := &http.Server{
		Handler: router,
		Addr: addr,

		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	log.Println("Serving on " + addr)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}