package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"c2c.in/api/internal/database/adapters"
	"c2c.in/api/internal/httphandlers"
)

func main() {

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	db := adapters.Connect("ctwoc-database")

	mux := http.NewServeMux()

	httphandlers.NewUnitHttpHandler(db).RegisterServiceWithMux(mux)
	httphandlers.NewTopicHttpHandler(db).RegisterServiceWithMux(mux)
	httphandlers.NewPathwayHttpHandler(db).RegisterServiceWithMux(mux)
	httphandlers.NewModuleHttpHandler(db).RegisterServiceWithMux(mux)

	fmt.Printf("Server is running on the http://localhost:%s\n", port)

	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
}
