// main.go
package main

import (
	"cashWise/routes"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	// Ініціалізація маршрутизатора з маршрутизацією для mux
	r := routes.InitializeRoutes()

	// Налаштування CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"https://example.com", "https://another-domain.com"}), // Список дозволених доменів
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),           // Дозволені методи
		handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Accept", "Authorization"}), // Дозволені заголовки
		handlers.AllowCredentials(), // Дозволити передачу cookies
	)(r)

	// Запуск сервера на порту 8080
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
