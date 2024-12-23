package routes

import (
	"cashWise/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

// Створення та налаштування роутів
func InitializeRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler).Methods("GET")

	// Роут для реєстрації нового користувача
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")

	// Роут для логіну користувача
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	// Роут для зміни даних користувача
	r.HandleFunc("/update-user/{userID}", handlers.UpdateUser).Methods("PUT")

	// Маршрут для видалення акаунту
	r.HandleFunc("/delete-user/{userID}", handlers.DeleteUser).Methods("DELETE")

	// Маршрути для категорій
	r.HandleFunc("/create-category", handlers.CreateCategory).Methods("POST")
	r.HandleFunc("/edit-category/{categoryID}", handlers.EditCategory).Methods("PUT")
	r.HandleFunc("/delete-category/{categoryID}", handlers.DeleteCategory).Methods("DELETE")
	r.HandleFunc("/categories/{budgetID}", handlers.GetCategoriesByBudgetID).Methods("GET")

	// main.go
	r.HandleFunc("/transaction", handlers.AddTransaction).Methods("POST")
	r.HandleFunc("/transaction/{transactionID}", handlers.EditTransaction).Methods("PUT")
	r.HandleFunc("/transaction/{transactionID}", handlers.DeleteTransaction).Methods("DELETE")
	r.HandleFunc("/transactions/filter", handlers.FilterByDate).Methods("GET")
	r.HandleFunc("/transactions/balance", handlers.CalculateBalance).Methods("GET")

	r.HandleFunc("/budgets", handlers.CreateBudget).Methods("POST")

	// Отримання бюджету за ID
	r.HandleFunc("/budgets/{budgetID}", handlers.GetBudgetByID).Methods("GET")
	// Оновлення бюджету
	r.HandleFunc("/budgets/{budgetID}", handlers.EditBudget).Methods("PUT")
	// Видалення бюджету
	r.HandleFunc("/budgets/{budgetID}", handlers.DeleteBudget).Methods("DELETE")
	// Перевірка ліміту бюджету
	r.HandleFunc("/budgets/{budgetID}/check-limit", handlers.CheckLimit).Methods("GET")

	r.HandleFunc("/goals", handlers.CreateGoalHandler).Methods("POST")                        // Створення нової фінансової цілі
	r.HandleFunc("/goals/{goalID}", handlers.EditGoalHandler).Methods("PUT")                  // Оновлення існуючої цілі
	r.HandleFunc("/goals/{goalID}", handlers.DeleteGoalHandler).Methods("DELETE")             // Видалення цілі
	r.HandleFunc("/goals/{goalID}/progress", handlers.UpdateProgressHandler).Methods("PATCH") // Оновлення прогресу
	r.HandleFunc("/goals/{goalID}/reminder", handlers.SendReminderHandler).Methods("POST")    // Надсилання нагадування
	return r
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
