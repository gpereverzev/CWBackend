package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"cashWise/models"
	"cashWise/repo"
	"cashWise/service"

	"github.com/gorilla/mux"
)

// AddTransaction - додає нову транзакцію
func AddTransaction(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userID") // Отримання userID з параметрів запиту
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	var newTransaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&newTransaction); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newTransaction.UserID = userID

	if err := repo.AddTransaction(newTransaction); err != nil {
		http.Error(w, fmt.Sprintf("Error adding transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction added successfully"})
}

// EditTransaction - редагує існуючу транзакцію
func EditTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	transactionIDStr := params["transactionID"]
	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	userIDStr := r.URL.Query().Get("userID") // Отримання userID з параметрів запиту
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	var updatedTransaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&updatedTransaction); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	updatedTransaction.UserID = userID

	if err := repo.EditTransaction(transactionID, updatedTransaction); err != nil {
		http.Error(w, fmt.Sprintf("Error editing transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction updated successfully"})
}

// DeleteTransaction - видаляє транзакцію
func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	transactionIDStr := params["transactionID"]
	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	userIDStr := r.URL.Query().Get("userID") // Отримання userID з параметрів запиту
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Викликаємо репозиторій для видалення транзакції за userID та transactionID
	if err := repo.DeleteTransaction(userID, transactionID); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction deleted successfully"})
}

// FilterByDate - обробник запиту для фільтрації транзакцій за заданим періодом (день, тиждень, місяць)
func FilterByDate(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userID") // Отримання userID з параметрів запиту
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	period := r.URL.Query().Get("period") // day, week, month
	dateStr := r.URL.Query().Get("date")  // дата у форматі YYYY-MM-DD

	if period == "" || dateStr == "" {
		http.Error(w, "Period and date are required", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	var startDate, endDate time.Time
	switch period {
	case "day":
		startDate = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		weekday := date.Weekday()
		offset := (int(weekday) - int(time.Monday) + 7) % 7             // Розрахунок зсуву до понеділка
		startDate = date.Add(-time.Duration(offset+7) * 24 * time.Hour) // На тиждень назад
		endDate = startDate.Add(7 * 24 * time.Hour)                     // Додаємо 7 днів для кінця тижня
	case "month":
		startDate = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
		endDate = startDate.AddDate(0, 1, -1) // Останній день місяця
	default:
		http.Error(w, "Invalid period, must be day, week, or month", http.StatusBadRequest)
		return
	}

	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	transactions, err := repo.GetTransactionsByDateAndUserID(startDateStr, endDateStr, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching transactions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding transactions: %v", err), http.StatusInternalServerError)
	}
}

// GetAllTransactions - отримує всі транзакції для користувача
func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userID") // Отримання userID з параметрів запиту
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	transactions, err := repo.GetAllTransactions(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching transactions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// Обробник для отримання загальної суми витрат для конкретного користувача
func GetTotalExpense(w http.ResponseWriter, r *http.Request) {
	// Отримання userID з параметрів запиту
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "Missing userID parameter", http.StatusBadRequest)
		return
	}

	// Перетворення userID на ціле число
	parsedUserID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	// Викликаємо функцію для обчислення загальної суми витрат
	totalExpense, err := repo.CalculateTotalExpense(parsedUserID)
	if err != nil {
		http.Error(w, "Error calculating total expense", http.StatusInternalServerError)
		return
	}

	// Повертаємо загальну суму витрат
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"totalExpense": totalExpense})
}

func GetTotalIncome(w http.ResponseWriter, r *http.Request) {
	// Отримання userID з параметрів запиту
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "Missing userID parameter", http.StatusBadRequest)
		return
	}

	// Перетворення userID на ціле число
	parsedUserID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	// Викликаємо функцію для обчислення загальної суми доходу
	// Тут ми використовуємо фільтр для доходу, так як категорія "income" передається в CalculateTotalIncome
	totalIncome, err := repo.CalculateTotalIncome(parsedUserID)
	if err != nil {
		// Виведення більш детальної помилки, якщо є проблеми з підключенням до бази даних або іншими аспектами
		http.Error(w, fmt.Sprintf("Error calculating total income: %v", err), http.StatusInternalServerError)
		return
	}

	// Повертаємо загальну суму доходу
	w.Header().Set("Content-Type", "application/json")
	// Відправляємо JSON-об'єкт з результатом
	json.NewEncoder(w).Encode(map[string]float64{"totalIncome": totalIncome})
}

// CalculateBalance - розраховує баланс користувача (доходи - витрати)
func CalculateBalance(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userID") // Отримуємо userID з параметрів запиту
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Отримуємо всі транзакції користувача з репозиторію
	transactions, err := repo.GetAllTransactionsByUser(userID)
	if err != nil {
		http.Error(w, "Error fetching transactions", http.StatusInternalServerError)
		return
	}

	totalIncome := 0.0
	totalExpense := 0.0

	// Обчислюємо суму доходів і витрат
	for _, transaction := range transactions {
		if transaction.Type == "income" {
			totalIncome += transaction.Amount
		} else if transaction.Type == "expense" {
			totalExpense += transaction.Amount
		}
	}

	// Розраховуємо баланс
	balance := totalIncome - totalExpense

	// Повертаємо баланс
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}

// GetTransactionsByUserIDAndTypeHandler - хендлер для отримання транзакцій по userID і типу
func GetTransactionsByUserIDAndTypeHandler(w http.ResponseWriter, r *http.Request) {
	// Отримуємо userID і тип транзакції з параметрів запиту
	userIDStr := r.URL.Query().Get("userID")
	transactionType := r.URL.Query().Get("type")

	if userIDStr == "" || transactionType == "" {
		http.Error(w, "Missing userID or type parameter", http.StatusBadRequest)
		return
	}

	// Перетворення userID на ціле число
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	// Викликаємо репозиторій для отримання транзакцій
	transactions, err := repo.GetTransactionsByUserIDAndType(userID, transactionType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching transactions for userID %d and type %s: %v", userID, transactionType, err), http.StatusInternalServerError)
		return
	}

	// Повертаємо знайдені транзакції у форматі JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding transactions to JSON: %v", err), http.StatusInternalServerError)
	}
}

// GetTransactionWithCategoryHandler - обробляє запит на отримання транзакції з деталями категорії
func GetTransactionWithCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Отримуємо transactionID з query параметрів
	transactionIDStr := r.URL.Query().Get("transactionID")
	if transactionIDStr == "" {
		http.Error(w, "Missing transactionID", http.StatusBadRequest)
		return
	}

	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		http.Error(w, "Invalid transactionID", http.StatusBadRequest)
		return
	}

	// Отримуємо дані
	data, err := service.GetTransactionWithCategory(transactionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving data: %v", err), http.StatusNotFound)
		return
	}

	// Повертаємо JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetTransactionsHandler - хендлер для отримання всіх транзакцій за userID
func GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	// Отримуємо параметр userID з query параметрів
	userIDStr := r.URL.Query().Get("userID")
	if userIDStr == "" {
		http.Error(w, "Missing userID", http.StatusBadRequest)
		return
	}

	// Перетворюємо userID на ціле число
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	// Викликаємо сервіс для отримання транзакцій
	transactions, err := service.GetTransactionsByUserID(userID)
	if err != nil {
		http.Error(w, "Error fetching transactions", http.StatusInternalServerError)
		return
	}

	// Відправляємо результат у вигляді JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
