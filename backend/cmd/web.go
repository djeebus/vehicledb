package cmd

import (
	"log"
	"net/http"
	"time"
	"vehicledb/api"
	"vehicledb/db"
	"vehicledb/graph"
)

func RunApiServer() {
	db.OpenDatabase("db.sqlite")

	schema, err := graph.GenerateSchema()
	if err != nil {
		log.Fatal("Failed to generate schema", err)
	}
	handler := api.NewHandler(schema)

	addr := "127.0.0.1:8000"

	srv := &http.Server{
		Handler: handler,
		Addr:    addr,

		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	log.Println("Serving on " + addr)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}