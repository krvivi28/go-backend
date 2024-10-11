package project

import (
	"GOLANG/project/handlers"
	"fmt"
	"log"
	"net/http"

	muxHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Project() {
	router := mux.NewRouter()

	handlers.UserRoutes(router)
	handlers.GeoRoutes(router)

	corsHeaders := muxHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	corsOrigins := muxHandlers.AllowedOrigins([]string{"*"})
	corsMethods := muxHandlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
	loggedRouter := muxHandlers.LoggingHandler(log.Writer(), router)
	corsRouter := muxHandlers.CORS(corsHeaders, corsOrigins, corsMethods)(loggedRouter)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", corsRouter); err != nil {
		fmt.Println("Failed to start server:", err)
	}

	fmt.Println("Server started on :8080")
}
