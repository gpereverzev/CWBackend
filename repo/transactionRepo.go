// repo/transactionRepo.go
package repo

import (
	"cashWise/db"
	"cashWise/models"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// AddTransaction - додає нову транзакцію
func AddTransaction(transaction models.Transaction) error {
	collection := db.GetTransactionCollection()
	_, err := collection.InsertOne(context.TODO(), transaction)
	if err != nil {
		log.Printf("Error inserting transaction: %v", err)
		return err
	}
	return nil
}

func EditTransaction(transactionID int, updatedTransaction models.Transaction) error {
	collection := db.GetTransactionCollection()
	filter := bson.M{"transactionID": transactionID}

	// Створення мапи для оновлення лише тих полів, які змінюються
	updateFields := bson.M{}
	if updatedTransaction.Type != "" {
		updateFields["type"] = updatedTransaction.Type
	}
	if updatedTransaction.Amount != 0 {
		updateFields["amount"] = updatedTransaction.Amount
	}
	if updatedTransaction.Date != "" {
		updateFields["date"] = updatedTransaction.Date
	}
	if updatedTransaction.Description != "" {
		updateFields["description"] = updatedTransaction.Description
	}

	// Якщо хоча б одне поле потрібно оновити, застосуємо операцію $set
	if len(updateFields) > 0 {
		update := bson.M{"$set": updateFields}
		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Printf("Error updating transaction: %v", err)
			return err
		}
	}

	return nil
}

// DeleteTransaction - видаляє транзакцію
func DeleteTransaction(transactionID int) error {
	collection := db.GetTransactionCollection()
	filter := bson.M{"transactionID": transactionID}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Error deleting transaction: %v", err)
		return err
	}
	return nil
}

func GetTransactionsByDate(startDate, endDate string) ([]models.Transaction, error) {
	collection := db.GetTransactionCollection()

	filter := bson.M{
		"date": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Error finding transactions: %v", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var transactions []models.Transaction
	for cursor.Next(context.TODO()) {
		var transaction models.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			log.Printf("Error decoding transaction: %v", err)
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return transactions, nil
}

// CalculateBalance - обчислює баланс бюджету за період
func CalculateBalance(startDate, endDate string) (float64, error) {
	// Перевірка на коректність формату дат
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start date format")
	}
	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		return 0, fmt.Errorf("invalid end date format")
	}

	// Використання рядкових дат для фільтрації
	filter := bson.M{
		"date": bson.M{
			"$gte": startDate, // порівняння по рядку
			"$lte": endDate,   // порівняння по рядку
		},
	}

	// Запит до MongoDB
	collection := db.GetTransactionCollection()
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Error finding transactions: %v", err)
		return 0, err
	}
	defer cursor.Close(context.TODO())

	// Обчислення балансу
	var balance float64
	for cursor.Next(context.TODO()) {
		var transaction models.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			log.Printf("Error decoding transaction: %v", err)
			return 0, err
		}
		// Зміна балансу в залежності від типу транзакції
		if transaction.Type == "expense" {
			balance -= transaction.Amount
		} else if transaction.Type == "income" {
			balance += transaction.Amount
		}
	}

	// Перевірка на помилки під час перебору курсора
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return 0, err
	}

	return balance, nil
}
