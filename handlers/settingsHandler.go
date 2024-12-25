package handlers

import (
	"cashWise/db"
	"fmt"
	"net/http"
	"strconv"
)

// HandleToggleDarkTheme - обробляє запит для перемикання darkTheme
func HandleToggleDarkTheme(w http.ResponseWriter, r *http.Request) {
	// Отримуємо userID з query параметрів
	userIDStr := r.URL.Query().Get("userID")
	if userIDStr == "" {
		http.Error(w, "Missing userID", http.StatusBadRequest)
		return
	}

	// Перетворюємо userID на int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	// Викликаємо функцію для перемикання darkTheme
	err = db.ToggleDarkTheme(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to toggle dark theme: %v", err), http.StatusInternalServerError)
		return
	}

	// Повертаємо успішний статус
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Dark theme toggled successfully"))
}
