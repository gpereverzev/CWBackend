// main.go
package main

import (
	"cashWise/routes"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	// Initialize the router with routes
	r := routes.InitializeRoutes()

	// Configure CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),                                                                               // Allow access from any origin
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),                                         // Allow HTTP methods
		handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Accept", "Authorization", "ngrok-skip-browser-warning"}), // Include additional headers
	)(r)

	// Preflight request handling for OPTIONS method
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		r.ServeHTTP(w, r)
	})

	// Start the server on port 8081
	log.Println("Server is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", corsHandler))
}
