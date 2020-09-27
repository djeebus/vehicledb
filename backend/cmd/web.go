package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"

	"vehicledb/api"
	"vehicledb/db"
	"vehicledb/graph"
)

var (
	rootCmd = &cobra.Command{
		Use: "vehicledb",
		Short: "Run the VehicleDB Backend",
		Run: runApiServer,
	}
	corsOrigins []string
	listen = ""
)

func init() {
	persistentFlags := rootCmd.PersistentFlags()
	persistentFlags.StringArrayVarP(
		&corsOrigins, "cosOrigin", "c", []string{"http://localhost:8080", "http://127.0.0.1:8080"}, "the host name of the frontend",
	)
	persistentFlags.StringVarP(
		&listen, "listen", "l", "127.0.0.1:8000", "the host to listen for requests on",
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runApiServer(cmd *cobra.Command, args []string) {
	db.OpenDatabase("db.sqlite")

	schema, err := graph.GenerateSchema()
	if err != nil {
		log.Fatal("Failed to generate schema", err)
	}
	handler := api.NewHandler(schema, corsOrigins)

	srv := &http.Server{
		Handler: handler,
		Addr:    listen,

		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	log.Println("Serving on " + listen)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
