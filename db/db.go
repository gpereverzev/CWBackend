// db/db.go
package db

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var categoryCollection *mongo.Collection
var transactionCollection *mongo.Collection
var budgetCollection *mongo.Collection
var goalCollection *mongo.Collection
var countersCollection *mongo.Collection
var settingCollection *mongo.Collection

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
	settingCollection = Client.Database("cashWiseDB").Collection("Settings")
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

func GetSettingCollection() *mongo.Collection {
	if settingCollection == nil {
		settingCollection = Client.Database("cashWiseDB").Collection("Settings")
	}
	return settingCollection
}

// ToggleDarkTheme - встановлює darkTheme на протилежне значення
func ToggleDarkTheme(userID int) error {
	collection := GetSettingCollection()

	// Знаходимо поточний стан darkTheme
	filter := bson.M{"userID": userID}
	var currentSetting bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&currentSetting)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("Settings for userID %d not found", userID)
			return errors.New("settings not found")
		}
		log.Printf("Error retrieving settings for userID %d: %v", userID, err)
		return err
	}

	// Витягуємо поточне значення darkTheme
	currentDarkTheme, ok := currentSetting["darkTheme"].(bool)
	if !ok {
		log.Printf("Invalid data format for darkTheme for userID %d", userID)
		return errors.New("invalid data format for darkTheme")
	}

	// Перемикаємо значення
	update := bson.M{"$set": bson.M{"darkTheme": !currentDarkTheme}}

	// Оновлюємо документ
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Failed to update darkTheme for userID %d: %v", userID, err)
		return err
	}

	return nil
}
