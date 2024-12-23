// handlers/transactionHandler.go
package handlers

import (
	"cashWise/models"
	"cashWise/repo"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// AddTransaction - додає нову транзакцію
func AddTransaction(w http.ResponseWriter, r *http.Request) {
	var newTransaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&newTransaction); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

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

	var updatedTransaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&updatedTransaction); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

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

	if err := repo.DeleteTransaction(transactionID); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction deleted successfully"})
}

// FilterByDate - обробник запиту для фільтрації транзакцій за заданим періодом (день, тиждень, місяць)
func FilterByDate(w http.ResponseWriter, r *http.Request) {
	// Отримання параметрів з запиту
	period := r.URL.Query().Get("period") // day, week, month
	dateStr := r.URL.Query().Get("date")  // дата у форматі YYYY-MM-DD

	// Перевірка наявності параметрів
	if period == "" || dateStr == "" {
		http.Error(w, "Period and date are required", http.StatusBadRequest)
		return
	}

	// Парсинг дати
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// Отримання діапазону дат залежно від періоду
	var startDate, endDate time.Time
	switch period {
	case "day":
		// Для дня використовуємо саме цю дату
		startDate = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endDate = startDate.Add(24 * time.Hour) // кінець того ж дня
	case "week":
		// Для тижня знаходимо початок тижня (понеділок)
		weekday := date.Weekday()
		offset := (int(weekday) - int(time.Monday) + 7) % 7 // кількість днів до понеділка
		startDate = date.Add(-time.Duration(offset) * 24 * time.Hour)
		endDate = startDate.Add(7 * 24 * time.Hour) // кінець тижня
	case "month":
		// Для місяця використовуємо перший і останній день місяця
		startDate = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
		endDate = startDate.AddDate(0, 1, 0) // додаємо 1 місяць
	default:
		http.Error(w, "Invalid period, must be day, week, or month", http.StatusBadRequest)
		return
	}

	// Переведення startDate та endDate в формат рядків "YYYY-MM-DD"
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	// Викликаємо репозиторій для отримання транзакцій
	transactions, err := repo.GetTransactionsByDate(startDateStr, endDateStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching transactions: %v", err), http.StatusInternalServerError)
		return
	}

	// Відправка результату у вигляді JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding transactions: %v", err), http.StatusInternalServerError)
	}
}

// CalculateBalance - обчислює баланс бюджету
func CalculateBalance(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	balance, err := repo.CalculateBalance(startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calculating balance: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}
