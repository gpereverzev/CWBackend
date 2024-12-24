package repo

import (
	"cashWise/db"
	"cashWise/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddTransaction - додає нову транзакцію
func AddTransaction(transaction models.Transaction) error {
	collection := db.GetTransactionCollection()

	_, err := collection.InsertOne(context.TODO(), transaction)
	return err
}

// EditTransaction - редагує існуючу транзакцію
func EditTransaction(transactionID int, transaction models.Transaction) error {
	collection := db.GetTransactionCollection()

	// Створюємо фільтр для пошуку за transactionID
	filter := bson.M{"transactionID": transactionID}

	// Перевіряємо чи існує транзакція з таким ID
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return mongo.ErrNoDocuments // Транзакція з таким ID не знайдена
	}

	// Створюємо мапу для оновлення, враховуючи лише ті поля, які змінюються
	update := bson.M{"$set": bson.M{}}

	// Додаємо лише змінені поля в update
	if transaction.Amount != 0 {
		update["$set"].(bson.M)["amount"] = transaction.Amount
	}
	if transaction.CategoryID != 0 {
		update["$set"].(bson.M)["category"] = transaction.CategoryID
	}
	if transaction.Description != "" {
		update["$set"].(bson.M)["description"] = transaction.Description
	}

	// Оновлюємо тільки змінені поля
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	return err
}

// DeleteTransaction - видаляє транзакцію для конкретного користувача
func DeleteTransaction(userID int, transactionID int) error {
	collection := db.GetTransactionCollection()

	// Фільтруємо транзакцію за userID та transactionID
	filter := bson.M{
		"userID":        userID,
		"transactionID": transactionID,
	}

	// Перевіряємо чи існує транзакція з таким ID
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return mongo.ErrNoDocuments // Транзакція з таким ID не знайдена
	}

	// Видаляємо транзакцію
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

// GetTransactionsByDateAndUserID - отримує транзакції за період та userID
func GetTransactionsByDateAndUserID(startDate, endDate string, userID int) ([]models.Transaction, error) {
	collection := db.GetTransactionCollection()

	filter := bson.M{
		"date": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
		"userID": userID,
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var transactions []models.Transaction
	for cursor.Next(context.TODO()) {
		var transaction models.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetAllTransactions - отримує всі транзакції для користувача
func GetAllTransactions(userID int) ([]models.Transaction, error) {
	collection := db.GetTransactionCollection()

	filter := bson.M{"userID": userID}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var transactions []models.Transaction
	for cursor.Next(context.TODO()) {
		var transaction models.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func FilterExpenseTransactions() bson.M {
	// Повертаємо фільтр, який шукає транзакції типу "expense"
	return bson.M{
		"type": "expense", // Тип транзакції має бути "expense"
	}
}

// Функція для обчислення загальної суми витрат з фільтрацією за userID
func CalculateTotalExpense(userID int) (float64, error) {
	// Отримання колекції транзакцій
	collection := db.GetTransactionCollection()

	// Фільтр для вибору транзакцій з типом "expense" і належністю конкретному користувачу
	filter := FilterExpenseTransactions()
	filter["userID"] = userID // Додаємо фільтрацію по userID

	// Контекст із таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Пошук відповідних документів
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error finding expense transactions: %v", err)
		return 0, err
	}
	defer cursor.Close(ctx)

	// Обчислення загальної суми
	var totalExpense float64
	for cursor.Next(ctx) {
		var transaction models.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			log.Printf("Error decoding transaction: %v", err)
			return 0, err
		}
		totalExpense += transaction.Amount
	}

	// Перевірка на помилки перебору курсора
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return 0, err
	}

	return totalExpense, nil
}

func FilterIncomeTransactions() bson.M {
	// Повертаємо фільтр, який шукає транзакції типу "income"
	return bson.M{
		"type": "income", // Тип транзакції має бути "expense"
	}
}

// CalculateTotalAmount - обчислює загальну суму для певного користувача та категорії
func CalculateTotalIncome(userID int, category string) (float64, error) {
	collection := db.GetTransactionCollection()

	// Отримуємо фільтр для доходів або витрат
	filter := FilterIncomeTransactions()

	// Пошук транзакцій за фільтром
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.TODO())

	var totalAmount float64
	for cursor.Next(context.TODO()) {
		var transaction models.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return 0, err
		}
		totalAmount += transaction.Amount
	}

	// Перевірка на помилки перебору курсора
	if err := cursor.Err(); err != nil {
		return 0, err
	}

	return totalAmount, nil
}

// GetAllTransactionsByUser - отримує всі транзакції для користувача
func GetAllTransactionsByUser(userID int) ([]models.Transaction, error) {
	collection := db.GetTransactionCollection()

	// Фільтр для отримання всіх транзакцій користувача
	filter := bson.M{
		"userID": userID,
	}

	// Отримуємо всі транзакції
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var transactions []models.Transaction
	for cursor.Next(context.TODO()) {
		var transaction models.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
