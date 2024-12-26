package handlers

import (
	"cashWise/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func ToggleGoalReminderHandler(w http.ResponseWriter, r *http.Request) {
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

	// Викликаємо функцію для перевірки транзакцій і формуємо відповідь
	response, err := service.CheckGoalTransactionsAndNotify(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking transactions: %v", err), http.StatusInternalServerError)
		return
	}

	// Повертаємо JSON-відповідь
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
