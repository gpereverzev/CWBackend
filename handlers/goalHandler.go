package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cashWise/models"
	"cashWise/repo"

	"github.com/gorilla/mux"
)

// CreateGoalHandler - обробник для створення фінансової цілі
func CreateGoalHandler(w http.ResponseWriter, r *http.Request) {
	var goal models.Goal
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := repo.CreateGoal(goal); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(goal)
}

// EditGoalHandler - обробник для оновлення цілі
func EditGoalHandler(w http.ResponseWriter, r *http.Request) {
	var updatedGoal models.Goal
	if err := json.NewDecoder(r.Body).Decode(&updatedGoal); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := repo.EditGoal(updatedGoal); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedGoal)
}

// DeleteGoalHandler - обробник для видалення цілі
func DeleteGoalHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	goalID, err := strconv.Atoi(params["goalID"])
	if err != nil {
		http.Error(w, "Invalid goal ID", http.StatusBadRequest)
		return
	}

	if err := repo.DeleteGoal(goalID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// // UpdateProgressHandler - обробник для оновлення прогресу
// func UpdateProgressHandler(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	goalID, err := strconv.Atoi(params["goalID"])
// 	if err != nil {
// 		http.Error(w, "Invalid goal ID", http.StatusBadRequest)
// 		return
// 	}

// 	var input struct {
// 		AddedAmount float64 `json:"addedAmount"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}

// 	if err := repo.UpdateProgress(goalID, input.AddedAmount); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

// SendReminderHandler - обробник для надсилання нагадувань
func SendReminderHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	goalID, err := strconv.Atoi(params["goalID"])
	if err != nil {
		http.Error(w, "Invalid goal ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Reminder string `json:"reminder"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := repo.SendReminder(goalID, input.Reminder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetGoalsByUserIDHandler - хендлер для отримання всіх цілей конкретного користувача
func GetGoalsByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	// Отримання userID з параметрів запиту
	userIDStr := r.URL.Query().Get("userID")
	if userIDStr == "" {
		http.Error(w, "Missing userID parameter", http.StatusBadRequest)
		return
	}

	// Перетворення userID на ціле число
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	// Викликаємо репозиторій для отримання цілей
	goals, err := repo.GetGoalsByUserID(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching goals for userID %d: %v", userID, err), http.StatusInternalServerError)
		return
	}

	// Повертаємо знайдені цілі у форматі JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(goals); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding goals to JSON: %v", err), http.StatusInternalServerError)
	}
}
