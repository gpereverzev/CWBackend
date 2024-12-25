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

// CreateCategory - створює нову категорію
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory models.Category

	// Декодуємо тіло запиту
	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Додаємо категорію через repo
	if err := repo.AddCategory(newCategory); err != nil {
		http.Error(w, fmt.Sprintf("Error creating category: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Category created successfully"})
}

// GetAllCategoriesByUserID - отримує всі категорії для конкретного користувача за його userID, переданим у query string
func GetAllCategoriesByUserID(w http.ResponseWriter, r *http.Request) {
	// Отримуємо параметр userID з query string
	userIDStr := r.URL.Query().Get("userID")
	if userIDStr == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	// Перетворюємо userID в ціле число
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid userID: %v", err), http.StatusBadRequest)
		return
	}

	// Отримуємо всі категорії для конкретного користувача
	categories, err := repo.GetAllCategoriesByUserIDLogic(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving categories: %v", err), http.StatusInternalServerError)
		return
	}

	// Повертаємо категорії у форматі JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding categories: %v", err), http.StatusInternalServerError)
	}
}

// GetCategoryByID - обробляє запит на отримання категорії за її ID
func GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	budgetIDStr := params["budgetID"]
	categoryIDStr := params["categoryID"]

	budgetID, err := strconv.Atoi(budgetIDStr)
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := repo.GetCategoryByID(budgetID, categoryID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving category: %v", err), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

// EditCategory - обробляє запит на оновлення категорії
func EditCategory(w http.ResponseWriter, r *http.Request) {
	// Отримуємо budgetID та categoryID з параметрів запиту
	vars := mux.Vars(r)
	budgetIDStr := vars["budgetID"]
	categoryIDStr := vars["categoryID"]

	budgetID, err := strconv.Atoi(budgetIDStr)
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Отримуємо оновлену категорію з тіла запиту
	var updatedCategory models.Category
	err = json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Оновлюємо категорію в БД
	err = repo.UpdateCategory(budgetID, categoryID, updatedCategory)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not update category: %v", err), http.StatusInternalServerError)
		return
	}

	// Повертаємо успішну відповідь
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Category updated"})
}

// DeleteCategory - обробляє запит на видалення категорії
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// Отримуємо budgetID та categoryID з параметрів запиту
	vars := mux.Vars(r)
	budgetIDStr := vars["budgetID"]
	categoryIDStr := vars["categoryID"]

	budgetID, err := strconv.Atoi(budgetIDStr)
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Викликаємо функцію для видалення категорії
	err = repo.DeleteCategory(budgetID, categoryID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete category: %v", err), http.StatusInternalServerError)
		return
	}

	// Повертаємо успішну відповідь
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted"})
}

func GetCategoriesByBudgetID(w http.ResponseWriter, r *http.Request) {
	// Отримуємо budgetID з параметрів запиту
	params := mux.Vars(r)
	budgetIDStr := params["budgetID"]
	budgetID, err := strconv.Atoi(budgetIDStr) // Перетворюємо budgetID на ціле число
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}

	// Отримуємо категорії для даного budgetID
	categories, err := repo.GetCategories(budgetID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving categories: %v", err), http.StatusInternalServerError)
		return
	}

	// Повертаємо категорії у відповіді
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

// GetCategoryByName - отримує категорію за її назвою через query string
func GetCategoryByName(w http.ResponseWriter, r *http.Request) {
	// Отримуємо параметр name з query string
	categoryName := r.URL.Query().Get("name")
	if categoryName == "" {
		http.Error(w, "name parameter is required", http.StatusBadRequest)
		return
	}

	// Викликаємо репозиторій для отримання категорії за назвою
	category, err := repo.GetCategoriesByName(categoryName)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error retrieving category: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Повертаємо знайдену категорію у форматі JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(category); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding category: %v", err), http.StatusInternalServerError)
	}
}
