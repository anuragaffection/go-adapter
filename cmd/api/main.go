package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"c2c.in/api/internal/database/adapters"
	"c2c.in/api/internal/httphandlers"
)

func main() {

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}


	db := adapters.Connect("ctwoc-database")

	mux:= http.NewServeMux()

	httphandlers.NewUnitHttpHandler(db).RegisterServiceWithMux(mux)
	httphandlers.NewTopicHttpHandler(db).RegisterServiceWithMux(mux)
	httphandlers.NewPathwayHttpHandler(db).RegisterServiceWithMux(mux)
	httphandlers.NewModuleHttpHandler(db).RegisterServiceWithMux(mux)


	fmt.Printf("Server is running on the http://localhost:%d\n", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)

	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
}
