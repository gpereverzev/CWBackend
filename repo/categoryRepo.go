package repo

import (
	"cashWise/db"
	"cashWise/models"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

// AddCategory - додає нову категорію
func AddCategory(category models.Category) error {
	collection := db.GetCategoryCollection()
	// Додаємо категорію через InsertOne
	_, err := collection.InsertOne(context.TODO(), category)
	if err != nil {
		log.Printf("Error inserting category: %v", err)
		return err
	}
	log.Println("Category added successfully")
	return nil
}

// GetCategories - отримує список всіх категорій для бюджету
func GetCategories(budgetID int) ([]models.Category, error) {
	collection := db.GetCategoryCollection()
	filter := bson.M{"budgetID": budgetID} // Фільтруємо по budgetID
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Error finding categories for budgetID %v: %v", budgetID, err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var categories []models.Category
	for cursor.Next(context.TODO()) {
		var category models.Category
		if err := cursor.Decode(&category); err != nil {
			log.Printf("Error decoding category: %v", err)
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return categories, nil
}

// GetCategoryByID - отримує категорію за її ID та budgetID
func GetCategoryByID(budgetID int, categoryID int) (*models.Category, error) {
	collection := db.GetCategoryCollection()
	filter := bson.M{"budgetID": budgetID, "categoryID": categoryID}
	var category models.Category
	err := collection.FindOne(context.TODO(), filter).Decode(&category)
	if err != nil {
		log.Printf("Error finding category by ID: %v", err)
		return nil, err
	}
	return &category, nil
}

// UpdateCategory - оновлює категорію
func UpdateCategory(budgetID int, categoryID int, updatedCategory models.Category) error {
	collection := db.GetCategoryCollection()

	// Створюємо карту для оновлення
	update := bson.M{
		"$set": bson.M{},
	}

	// Додаємо лише непусті поля
	if updatedCategory.Name != "" {
		update["$set"].(bson.M)["name"] = updatedCategory.Name
	}
	if updatedCategory.Description != "" {
		update["$set"].(bson.M)["description"] = updatedCategory.Description
	}

	// Виконання оновлення
	if len(update["$set"].(bson.M)) > 0 {
		filter := bson.M{"budgetID": budgetID, "categoryID": categoryID}
		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return fmt.Errorf("Error updating category: %v", err)
		}
	}

	return nil
}

// DeleteCategory - видаляє категорію за її ID та budgetID
func DeleteCategory(budgetID int, categoryID int) error {
	collection := db.GetCategoryCollection()

	// Фільтр для видалення категорії по budgetID та categoryID
	filter := bson.M{"budgetID": budgetID, "categoryID": categoryID}

	// Видалення категорії з колекції
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Error deleting category: %v", err)
		return err
	}

	log.Println("Category deleted successfully")
	return nil
}
