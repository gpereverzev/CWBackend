package repo

import (
	"context"
	"errors"
	"fmt"
	"log"

	"cashWise/db"
	"cashWise/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetNextGoalID - отримує наступний ID для фінансової цілі
func GetNextGoalID() (int, error) {
	// Оновлюємо лічильник в колекції Counters
	filter := bson.M{"_id": "goalID"} // ми будемо використовувати "goalID" як ідентифікатор лічильника
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
		return 0, fmt.Errorf("failed to get next goalID: %v", err)
	}
	return result.Seq, nil
}

// CreateGoal - створює нову фінансову ціль з автоінкрементом
func CreateGoal(goal models.Goal) error {
	collection := db.GetGoalCollection()

	// Отримуємо наступний доступний ID для цілі
	goalID, err := GetNextGoalID()
	if err != nil {
		log.Printf("Error getting next goalID: %v", err)
		return fmt.Errorf("error getting next goalID: %v", err)
	}

	// Призначаємо цей ID фінансовій цілі
	goal.GoalID = goalID
	goal.Status = "in progress"

	// Вставляємо нову ціль в колекцію
	_, err = collection.InsertOne(context.TODO(), goal)
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

// // UpdateProgress - оновлює статус прогресу фінансової цілі
// func UpdateProgress(goalID int, addedAmount float64) error {
// 	collection := db.GetGoalCollection()
// 	filter := bson.M{"goalID": goalID}

// 	// Знайти поточний стан цілі
// 	var goal models.Goal
// 	err := collection.FindOne(context.TODO(), filter).Decode(&goal)
// 	if err != nil {
// 		log.Printf("Error fetching goal: %v", err)
// 		return fmt.Errorf("error fetching goal: %v", err)
// 	}

// 	// Оновити поточну суму
// 	newCurrentAmount := goal.CurrentAmount + addedAmount
// 	update := bson.M{"$set": bson.M{"currentAmount": newCurrentAmount}}

// 	if newCurrentAmount >= goal.TargetAmount {
// 		update["$set"].(bson.M)["status"] = "completed"
// 	}

// 	_, err = collection.UpdateOne(context.TODO(), filter, update)
// 	if err != nil {
// 		log.Printf("Error updating progress: %v", err)
// 		return fmt.Errorf("error updating progress: %v", err)
// 	}

// 	return nil
// }

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

// GetGoalsByUserID - отримує всі цілі для конкретного користувача за його userID
func GetGoalsByUserID(userID int) ([]models.Goal, error) {
	collection := db.GetGoalCollection()
	filter := bson.M{"userID": userID} // Фільтруємо за userID

	var goals []models.Goal
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Error fetching goals for userID %d: %v", userID, err)
		return nil, fmt.Errorf("error fetching goals for userID %d: %v", userID, err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var goal models.Goal
		if err := cursor.Decode(&goal); err != nil {
			log.Printf("Error decoding goal: %v", err)
			return nil, fmt.Errorf("error decoding goal: %v", err)
		}
		goals = append(goals, goal)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return goals, nil
}
