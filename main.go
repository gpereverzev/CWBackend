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
		handlers.AllowedOrigins([]string{"*"}), // Разрешить доступ с любых источников
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Accept", "Authorization"}),
	)(r)

	// Запуск сервера на порту 8080
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
