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

// ToggleTermsConditionHandler - хендлер для перемикання termsCondition
func ToggleTermsConditionHandler(w http.ResponseWriter, r *http.Request) {
	// Отримуємо userID з параметрів запиту
	userIDStr := r.URL.Query().Get("userID")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Перетворюємо userID з рядка в ціле число
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Викликаємо функцію ToggleTermsCondition
	err = db.ToggleTermsCondition(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error toggling terms condition: %v", err), http.StatusInternalServerError)
		return
	}

	// Відповідь на успішне виконання
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Terms and conditions toggled successfully"))
}

// ToggleNotificationsHandler - хендлер для перемикання notifications
func ToggleNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	// Отримуємо userID з параметрів запиту
	userIDStr := r.URL.Query().Get("userID")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Перетворюємо userID з рядка в ціле число
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Викликаємо функцію ToggleNotifications для зміни стану notifications
	err = db.ToggleNotifications(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error toggling notifications: %v", err), http.StatusInternalServerError)
		return
	}

	// Відповідь на успішне виконання
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notifications toggled successfully"))
}
