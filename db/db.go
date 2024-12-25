// db/db.go
package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var categoryCollection *mongo.Collection
var transactionCollection *mongo.Collection
var budgetCollection *mongo.Collection
var goalCollection *mongo.Collection
var countersCollection *mongo.Collection

func init() {
	// Створення параметрів підключення
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Підключення до MongoDB
	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}

	// Перевірка підключення
	err = Client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Could not ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB")

	// Ініціалізація колекції після успішного підключення
	countersCollection = Client.Database("cashWiseDB").Collection("counters")
	categoryCollection = Client.Database("cashWiseDB").Collection("Category")
	transactionCollection = Client.Database("cashWiseDB").Collection("Transactions")
	budgetCollection = Client.Database("cashWiseDB").Collection("Budgets")
	goalCollection = Client.Database("cashWiseDB").Collection("Goals")
}

// GetCategoryCollection - повертає колекцію категорій
func GetCategoryCollection() *mongo.Collection {
	if categoryCollection == nil {
		categoryCollection = Client.Database("cashWiseDB").Collection("Category")
	}
	return categoryCollection
}

// GetCategoryCollection - повертає колекцію категорій
func GetTransactionCollection() *mongo.Collection {
	if transactionCollection == nil {
		transactionCollection = Client.Database("cashWiseDB").Collection("Transactions")
	}
	return transactionCollection
}

func GetBudgetCollection() *mongo.Collection {
	if budgetCollection == nil {
		budgetCollection = Client.Database("cashWiseDB").Collection("Budgets")
	}
	return budgetCollection
}

func GetGoalCollection() *mongo.Collection {
	if goalCollection == nil {
		goalCollection = Client.Database("cashWiseDB").Collection("Goals")
	}
	return goalCollection
}

func GetCountersCollection() *mongo.Collection {
	if countersCollection == nil {
		countersCollection = Client.Database("cashWiseDB").Collection("counters")
	}
	return countersCollection
}
