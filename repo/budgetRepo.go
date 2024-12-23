package repo

import (
	"cashWise/db"
	"cashWise/models"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

// createBudget - створює новий бюджет у базі даних
func CreateBudget(newBudget models.Budget) error {
	collection := db.GetBudgetCollection()
	_, err := collection.InsertOne(context.TODO(), newBudget)
	if err != nil {
		log.Printf("Error creating budget: %v", err)
		return fmt.Errorf("error creating budget: %v", err)
	}
	return nil
}

// EditBudget - оновлює бюджет у базі даних
func EditBudget(updatedBudget models.Budget) error {
	collection := db.GetBudgetCollection()
	filter := bson.M{"budgetID": updatedBudget.BudgetID}

	// Динамічне формування об'єкта оновлення
	update := bson.M{
		"$set": bson.M{},
	}

	if updatedBudget.Name != "" {
		update["$set"].(bson.M)["name"] = updatedBudget.Name
	}
	if updatedBudget.InitialBalance != 0 {
		update["$set"].(bson.M)["initialBalance"] = updatedBudget.InitialBalance
	}
	if updatedBudget.Limit != 0 {
		update["$set"].(bson.M)["limit"] = updatedBudget.Limit
	}
	if updatedBudget.Period != "" {
		update["$set"].(bson.M)["period"] = updatedBudget.Period
	}

	// Виконуємо оновлення лише якщо є щось для оновлення
	if len(update["$set"].(bson.M)) == 0 {
		return fmt.Errorf("no fields to update")
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Error updating budget: %v", err)
		return fmt.Errorf("error updating budget: %v", err)
	}

	return nil
}

// deleteBudget - видаляє бюджет з бази даних
func DeleteBudget(budgetID int) error {
	collection := db.GetBudgetCollection()
	filter := bson.M{"budgetID": budgetID}

	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Error deleting budget: %v", err)
		return fmt.Errorf("error deleting budget: %v", err)
	}
	return nil
}

// getBudgetByID - отримує бюджет за ID
func GetBudgetByID(budgetID int) (models.Budget, error) {
	var budget models.Budget
	collection := db.GetBudgetCollection()
	filter := bson.M{"budgetID": budgetID}

	err := collection.FindOne(context.TODO(), filter).Decode(&budget)
	if err != nil {
		log.Printf("Error finding budget: %v", err)
		return budget, fmt.Errorf("error finding budget: %v", err)
	}
	return budget, nil
}

// checkLimit - перевіряє ліміт бюджету і повертає повідомлення
func CheckLimit(budgetID int) (string, error) {
	budget, err := GetBudgetByID(budgetID)
	if err != nil {
		return "", err
	}

	if budget.InitialBalance > budget.Limit {
		return "Limit exceeded!", nil
	}
	return "Limit is within range.", nil
}
