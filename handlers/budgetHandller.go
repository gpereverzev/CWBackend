package handlers

import (
	"cashWise/models"
	"cashWise/repo"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateBudget - створює новий бюджет
func CreateBudget(w http.ResponseWriter, r *http.Request) {
	var newBudget models.Budget

	// Декодуємо тіло запиту
	if err := json.NewDecoder(r.Body).Decode(&newBudget); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Додаємо бюджет через repo
	if err := repo.CreateBudget(newBudget); err != nil {
		http.Error(w, fmt.Sprintf("Error creating budget: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Budget created successfully"})
}

// GetBudgetByID - обробляє запит на отримання бюджету за ID
func GetBudgetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	budgetIDStr := params["budgetID"]
	budgetID, err := strconv.Atoi(budgetIDStr)
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}

	budget, err := repo.GetBudgetByID(budgetID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving budget: %v", err), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(budget)
}

// EditBudget - обробляє запит на оновлення бюджету
func EditBudget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	budgetIDStr := vars["budgetID"]

	budgetID, err := strconv.Atoi(budgetIDStr)
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}

	var updatedBudget models.Budget
	err = json.NewDecoder(r.Body).Decode(&updatedBudget)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	updatedBudget.BudgetID = budgetID // Встановлюємо ID в оновлений бюджет

	err = repo.EditBudget(updatedBudget)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating budget: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Budget updated successfully"})
}

// DeleteBudget - обробляє запит на видалення бюджету
func DeleteBudget(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	budgetIDStr := params["budgetID"]
	budgetID, err := strconv.Atoi(budgetIDStr)
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}

	err = repo.DeleteBudget(budgetID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting budget: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Budget deleted successfully"})
}

// CheckLimit - обробляє запит на перевірку ліміту бюджету
func CheckLimit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	budgetIDStr := params["budgetID"]
	budgetID, err := strconv.Atoi(budgetIDStr)
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}

	message, err := repo.CheckLimit(budgetID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking limit: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
