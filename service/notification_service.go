package service

import (
	"cashWise/models"
	"cashWise/repo"
	"encoding/json"
	"fmt"
)

func CheckGoalTransactionsAndNotify(userID int) (bool, error) {
	// 1. Отримуємо транзакції типу "goal" для цього користувача за поточний місяць
	transactions, err := repo.GetTransactionsByUserIDAndType(userID, "goal")
	if err != nil {
		return false, fmt.Errorf("Error fetching transactions: %v", err)
	}

	// Якщо транзакцій немає
	if len(transactions) == 0 {
		// 2. Перевіряємо налаштування користувача на дозволеність сповіщень
		settings, err := repo.GetSettingsByUserID(userID)
		if err != nil {
			return false, fmt.Errorf("Error fetching settings: %v", err)
		}

		// Якщо сповіщення увімкнено, повертаємо, що потрібно надіслати повідомлення
		if settings.Notifications {
			return true, nil // Повідомлення має бути надіслано
		}
	}

	// Якщо є транзакції або сповіщення вимкнені — не надсилаємо повідомлення
	return false, nil
}

// Створення JSON повідомлення про досягнення голу
func CreateGoalNotification(user models.User, goal string) ([]byte, error) {
	// Створюємо повідомлення у форматі JSON
	notification := map[string]interface{}{
		"userID":   user.UserID,
		"fullName": user.FullName,
		"email":    user.Email,
		"role":     user.Role,
		"goal":     goal,
		"message":  fmt.Sprintf("Congratulations, you've reached your goal: %s!", goal),
	}

	// Конвертуємо в JSON
	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON: %v", err)
	}

	return notificationJSON, nil
}
