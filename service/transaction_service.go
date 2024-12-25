package service

import (
	"cashWise/db"
	"cashWise/models"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

// GetTransactionWithCategory - отримує транзакцію з деталями категорії та повертає вивід без вкладеного об'єкта category
func GetTransactionWithCategory(transactionID int) (map[string]interface{}, error) {
	transactionCollection := db.GetTransactionCollection()
	categoryCollection := db.GetCategoryCollection()

	// Отримуємо транзакцію
	var transaction models.Transaction
	err := transactionCollection.FindOne(context.TODO(), bson.M{"transactionID": transactionID}).Decode(&transaction)
	if err != nil {
		log.Printf("Error retrieving transaction: %v", err)
		return nil, errors.New("transaction not found")
	}

	// Отримуємо категорію, пов'язану з транзакцією
	var category models.Category
	err = categoryCollection.FindOne(context.TODO(), bson.M{"categoryID": transaction.CategoryID}).Decode(&category)
	if err != nil {
		log.Printf("Error retrieving category: %v", err)
		return nil, errors.New("category not found")
	}

	// Формуємо фінальний JSON
	result := map[string]interface{}{
		"transactionID": transaction.TransactionID,
		"userID":        transaction.UserID,
		"categoryID":    transaction.CategoryID,
		"type":          transaction.Type,
		"amount":        transaction.Amount,
		"date":          transaction.Date,
		"description":   transaction.Description,
		"icon":          category.Icon, // Додаємо icon з категорії
	}

	return result, nil
}

// GetTransactionsByUserID - отримує всі транзакції для заданого userID і повертає їх у потрібному форматі
func GetTransactionsByUserID(userID int) ([]map[string]interface{}, error) {
	transactionCollection := db.GetTransactionCollection()
	categoryCollection := db.GetCategoryCollection()

	// Отримуємо всі транзакції для користувача за userID
	cursor, err := transactionCollection.Find(context.TODO(), bson.M{"userID": userID})
	if err != nil {
		log.Printf("Error retrieving transactions: %v", err)
		return nil, errors.New("transactions not found")
	}
	defer cursor.Close(context.TODO())

	var transactions []map[string]interface{}

	// Обробляємо кожну транзакцію
	for cursor.Next(context.TODO()) {
		var transaction models.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			log.Printf("Error decoding transaction: %v", err)
			continue
		}

		// Отримуємо категорію для поточної транзакції
		var category models.Category
		err := categoryCollection.FindOne(context.TODO(), bson.M{"categoryID": transaction.CategoryID}).Decode(&category)
		if err != nil {
			log.Printf("Error retrieving category: %v", err)
			continue
		}

		// Перевірка наявності іконки
		if category.Icon == "" {
			log.Printf("Category icon is empty for categoryID: %d", transaction.CategoryID)
		}

		// Форматуємо дату у правильний формат
		formattedDate := transaction.Date

		// Формуємо об'єкт для результату
		transactionData := map[string]interface{}{
			"transactionID": transaction.TransactionID,
			"userID":        transaction.UserID,
			"categoryID":    transaction.CategoryID,
			"type":          transaction.Type,
			"amount":        transaction.Amount,
			"date":          formattedDate, // Форматуємо дату
			"description":   transaction.Description,
			"icon":          category.Icon, // Додаємо icon з категорії
		}

		// Додаємо транзакцію в результат
		transactions = append(transactions, transactionData)
	}

	// Перевіряємо, чи є помилка після завершення обробки курсору
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return transactions, nil
}
