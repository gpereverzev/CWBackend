package repo

import (
	"cashWise/db"
	"cashWise/models"
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetNextTransactionID отримує наступний унікальний ID для транзакції
func GetNextTransactionID() (int, error) {
	// Оновлюємо лічильник в колекції Counters
	filter := bson.M{"_id": "transactionID"} // використаємо "transactionID" як ідентифікатор лічильника
	update := bson.M{
		"$inc": bson.M{"seq": 1}, // інкрементуємо значення на 1
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After) // створити документ, якщо його немає

	var result struct {
		Seq int `bson:"seq"`
	}

	// Повертаємо наступне значення ID
	err := db.GetCountersCollection().FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&result)
	if err != nil {
		return 0, fmt.Errorf("failed to get next transactionID: %v", err)
	}
	return result.Seq, nil
}

// AddTransaction додає нову транзакцію
func AddTransaction(transaction models.Transaction) error {
	collection := db.GetTransactionCollection()

	// Отримуємо новий transactionID
	transactionID, err := GetNextTransactionID()
	if err != nil {
		return fmt.Errorf("failed to get next transaction ID: %v", err)
	}

	// Призначаємо новий ID транзакції
	transaction.TransactionID = transactionID

	_, err = collection.InsertOne(context.TODO(), transaction)
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

func sortTransactionsByDate(transactions []models.Transaction) ([]models.Transaction, error) {
	sort.Slice(transactions, func(i, j int) bool {
		dateI, errI := time.Parse("2006-01-02", transactions[i].Date)
		dateJ, errJ := time.Parse("2006-01-02", transactions[j].Date)
		if errI != nil || errJ != nil {
			log.Printf("Error parsing date: %v, %v", errI, errJ)
			return false
		}
		return dateI.After(dateJ) // Найновіші транзакції спочатку
	})
	return transactions, nil
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

	return sortTransactionsByDate(transactions)
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

	return sortTransactionsByDate(transactions)
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

	return sortTransactionsByDate(transactions)
}

// Функція для фільтрації витрат
func FilterExpenseTransactions(userID int, transactionTypes []string) bson.M {
	// Повертаємо фільтр, який шукає транзакції з одним з типів з переданого списку і належність до певного користувача
	return bson.M{
		"type":   bson.M{"$in": transactionTypes}, // Використовуємо $in для кількох типів транзакцій
		"userID": userID,                          // Фільтруємо по userID
	}
}

// Функція для обчислення загальної суми витрат з фільтрацією за userID
func CalculateTotalExpense(userID int) (float64, error) {
	// Отримання колекції транзакцій
	collection := db.GetTransactionCollection()

	// Фільтр для вибору транзакцій з типом "expense" і належністю конкретному користувачу
	filter := FilterExpenseTransactions(userID, []string{"expense", "goal"})

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

// Функція для фільтрації доходів
func FilterIncomeTransactions(userID int) bson.M {
	// Повертаємо фільтр, який шукає транзакції типу "income" і належність до певного користувача
	return bson.M{
		"type":   "income", // Тип транзакції має бути "income"
		"userID": userID,   // Фільтруємо по userID
	}
}

// Функція для обчислення загальної суми доходів з фільтрацією за userID
func CalculateTotalIncome(userID int) (float64, error) {
	collection := db.GetTransactionCollection()

	// Отримуємо фільтр для доходів або витрат
	filter := FilterIncomeTransactions(userID)

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

// GetTransactionsByUserIDAndType - отримує транзакції для конкретного користувача за userID і типом
func GetTransactionsByUserIDAndType(userID int, transactionType string) ([]models.Transaction, error) {
	collection := db.GetTransactionCollection()

	// Фільтруємо за userID і типом транзакції
	filter := bson.M{
		"userID": userID,
		"type":   transactionType,
	}

	var transactions []models.Transaction
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Error fetching transactions for userID %d and type %s: %v", userID, transactionType, err)
		return nil, fmt.Errorf("error fetching transactions for userID %d and type %s: %v", userID, transactionType, err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var transaction models.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			log.Printf("Error decoding transaction: %v", err)
			return nil, fmt.Errorf("error decoding transaction: %v", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return sortTransactionsByDate(transactions)
}
