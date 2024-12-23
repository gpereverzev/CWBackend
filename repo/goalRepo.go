package repo

import (
	"context"
	"errors"
	"fmt"
	"log"

	"cashWise/db"
	"cashWise/models"

	"go.mongodb.org/mongo-driver/bson"
)

// CreateGoal - створює нову фінансову ціль
func CreateGoal(goal models.Goal) error {
	collection := db.GetGoalCollection()
	goal.Status = "in progress"
	_, err := collection.InsertOne(context.TODO(), goal)
	if err != nil {
		log.Printf("Error creating goal: %v", err)
		return fmt.Errorf("error creating goal: %v", err)
	}
	return nil
}

// EditGoal - оновлює існуючу фінансову ціль
func EditGoal(updatedGoal models.Goal) error {
	collection := db.GetGoalCollection()
	filter := bson.M{"goalID": updatedGoal.GoalID}

	updateFields := bson.M{}
	if updatedGoal.Name != "" {
		updateFields["name"] = updatedGoal.Name
	}
	if updatedGoal.TargetAmount != 0 {
		updateFields["targetAmount"] = updatedGoal.TargetAmount
	}
	if updatedGoal.Deadline != "" {
		updateFields["deadline"] = updatedGoal.Deadline
	}
	if len(updateFields) == 0 {
		return errors.New("no fields to update")
	}

	update := bson.M{"$set": updateFields}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Error editing goal: %v", err)
		return fmt.Errorf("error editing goal: %v", err)
	}
	return nil
}

// DeleteGoal - видаляє фінансову ціль із системи
func DeleteGoal(goalID int) error {
	collection := db.GetGoalCollection()
	filter := bson.M{"goalID": goalID}

	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Error deleting goal: %v", err)
		return fmt.Errorf("error deleting goal: %v", err)
	}
	return nil
}

// UpdateProgress - оновлює статус прогресу фінансової цілі
func UpdateProgress(goalID int, addedAmount float64) error {
	collection := db.GetGoalCollection()
	filter := bson.M{"goalID": goalID}

	// Знайти поточний стан цілі
	var goal models.Goal
	err := collection.FindOne(context.TODO(), filter).Decode(&goal)
	if err != nil {
		log.Printf("Error fetching goal: %v", err)
		return fmt.Errorf("error fetching goal: %v", err)
	}

	// Оновити поточну суму
	newCurrentAmount := goal.CurrentAmount + addedAmount
	update := bson.M{"$set": bson.M{"currentAmount": newCurrentAmount}}

	if newCurrentAmount >= goal.TargetAmount {
		update["$set"].(bson.M)["status"] = "completed"
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Error updating progress: %v", err)
		return fmt.Errorf("error updating progress: %v", err)
	}

	return nil
}

// SendReminder - надсилає нагадування про дедлайн або прогрес
func SendReminder(goalID int, reminder string) error {
	goal, err := GetGoalByID(goalID)
	if err != nil {
		return fmt.Errorf("error fetching goal: %v", err)
	}

	log.Printf("Reminder for Goal '%s' (ID: %d): %s", goal.Name, goalID, reminder)
	return nil
}

// GetGoalByID - отримує ціль за ID
func GetGoalByID(goalID int) (models.Goal, error) {
	collection := db.GetGoalCollection()
	filter := bson.M{"goalID": goalID}

	var goal models.Goal
	err := collection.FindOne(context.TODO(), filter).Decode(&goal)
	if err != nil {
		log.Printf("Error fetching goal by ID: %v", err)
		return goal, fmt.Errorf("error fetching goal by ID: %v", err)
	}

	return goal, nil
}
