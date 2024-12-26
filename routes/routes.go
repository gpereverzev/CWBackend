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
	r.HandleFunc("/update-user", handlers.UpdateUser).Methods("PUT")

	// Маршрут для видалення акаунту
	r.HandleFunc("/delete-user/{userID}", handlers.DeleteUser).Methods("DELETE")

	// Роут для отримання даних користувача по його айдi
	r.HandleFunc("/getUserByID/{id}", handlers.GetUserByID).Methods("GET")
	r.HandleFunc("/userByEmail", handlers.GetUserByEmailHandler).Methods("GET")

	r.HandleFunc("/settings/toggle-dark-theme", handlers.HandleToggleDarkTheme).Methods("POST")
	r.HandleFunc("/settings/toggle-terms-condition", handlers.ToggleTermsConditionHandler).Methods("POST")
	r.HandleFunc("/settings/notification", handlers.ToggleNotificationsHandler).Methods("POST")

	r.HandleFunc("/getUserAndSettings", handlers.GetUserAndSettingsHandler).Methods("GET")

	// Маршрути для категорій
	r.HandleFunc("/create-category", handlers.CreateCategory).Methods("POST")
	r.HandleFunc("/edit-category/{categoryID}", handlers.EditCategory).Methods("PUT")
	r.HandleFunc("/delete-category/{categoryID}", handlers.DeleteCategory).Methods("DELETE")
	r.HandleFunc("/categories", handlers.GetAllCategoriesByUserID).Methods("GET")
	r.HandleFunc("/category", handlers.GetCategoryByName).Methods("GET")
	r.HandleFunc("/get/category", handlers.GetCategoryByID).Methods("GET")

	// main.go

	r.HandleFunc("/transaction/getTotalExpense", handlers.GetTotalExpense).Methods("GET") // Отримати загальну суму витрат для userID
	r.HandleFunc("/transaction/getTotalIncome", handlers.GetTotalIncome).Methods("GET")   // Отримати загальну суму доходів для userID

	r.HandleFunc("/transaction", handlers.AddTransaction).Methods("POST")                             // Додати транзакцію
	r.HandleFunc("/transaction/{transactionID}", handlers.EditTransaction).Methods("PUT")             // Редагувати транзакцію
	r.HandleFunc("/transaction/{transactionID}", handlers.DeleteTransaction).Methods("DELETE")        // Видалити транзакцію
	r.HandleFunc("/transactions/filter", handlers.FilterByDate).Methods("GET")                        // Фільтрувати транзакції за датою
	r.HandleFunc("/transactions/balance", handlers.CalculateBalance).Methods("GET")                   // Розрахувати баланс
	r.HandleFunc("/transactions", handlers.GetAllTransactions).Methods("GET")                         // Отримати всі транзакції
	r.HandleFunc("/transactions-goal", handlers.GetTransactionsByUserIDAndTypeHandler).Methods("GET") // Отримати всі транзакції за userID та type
	r.HandleFunc("/transaction/details", handlers.GetTransactionWithCategoryHandler).Methods("GET")
	r.HandleFunc("/transactionsIcon", handlers.GetTransactionsHandler).Methods("GET")

	r.HandleFunc("/budgets", handlers.CreateBudget).Methods("POST") // Створити бюджет

	// Отримання бюджету за ID
	r.HandleFunc("/budgets/{budgetID}", handlers.GetBudgetByID).Methods("GET")
	// Оновлення бюджету
	r.HandleFunc("/budgets/{budgetID}", handlers.EditBudget).Methods("PUT")
	// Видалення бюджету
	r.HandleFunc("/budgets/{budgetID}", handlers.DeleteBudget).Methods("DELETE")
	// Перевірка ліміту бюджету
	r.HandleFunc("/budgets/{budgetID}/check-limit", handlers.CheckLimit).Methods("GET")

	//Получение текущих накоплений
	r.HandleFunc("/goals", handlers.CreateGoalHandler).Methods("POST")            // Створення нової фінансової цілі
	r.HandleFunc("/goals/{goalID}", handlers.EditGoalHandler).Methods("PUT")      // Оновлення існуючої цілі
	r.HandleFunc("/goals/{goalID}", handlers.DeleteGoalHandler).Methods("DELETE") // Видалення цілі
	//r.HandleFunc("/goals/{goalID}/progress", handlers.UpdateProgressHandler).Methods("PATCH") // Оновлення прогресу
	r.HandleFunc("/goals/{goalID}/reminder", handlers.SendReminderHandler).Methods("POST") // Надсилання нагадування
	r.HandleFunc("/goals", handlers.GetGoalsByUserIDHandler).Methods("GET")                // Отримання всіх цілей юзера

	// Роут для нагадування про транзакції "goal"
	r.HandleFunc("/goal-reminder", handlers.ToggleGoalReminderHandler).Methods("GET")

	return r
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
