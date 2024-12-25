package repo

import (
	"cashWise/db"
	"cashWise/models"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getNextSequence(sequenceName string) (int, error) {
	collection := db.GetCountersCollection()
	filter := bson.M{"_id": sequenceName}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var result struct {
		Seq int `bson:"seq"`
	}

	// Оновлюємо документ, збільшуючи seq на 1
	err := collection.FindOneAndUpdate(context.TODO(), filter, update, options).Decode(&result)
	if err == mongo.ErrNoDocuments {
		// Якщо запису не існує, створюємо його з початковим значенням seq = 1
		log.Println("Document not found, creating new counter")
		_, err := collection.InsertOne(context.TODO(), bson.M{"_id": sequenceName, "seq": 1})
		if err != nil {
			return 0, fmt.Errorf("error initializing counter: %v", err)
		}
		return 1, nil
	} else if err != nil {
		log.Printf("Error incrementing counter for %s: %v", sequenceName, err)
		return 0, fmt.Errorf("error incrementing counter: %v", err)
	}

	// Перевірка значення
	log.Printf("New sequence value for %s: %d", sequenceName, result.Seq)

	return result.Seq, nil
}

// AddCategory - додає нову категорію з автоінкрементом і іконкою
func AddCategory(category models.Category) error {
	collection := db.GetCategoryCollection()

	// Отримуємо новий CategoryID
	categoryID, err := getNextSequence("categoryID")
	if err != nil {
		log.Printf("Error getting next sequence for categoryID: %v", err)
		return err
	}
	category.CategoryID = categoryID

	// Перевірка наявності іконки
	if category.Icon == "" {
		log.Println("Category icon is empty, setting default icon")
		category.Icon = "default-icon.png" // Можна вказати стандартну іконку, якщо вона не надана
	}

	// Додаємо категорію через InsertOne
	_, err = collection.InsertOne(context.TODO(), category)
	if err != nil {
		log.Printf("Error inserting category: %v", err)
		return err
	}
	log.Println("Category added successfully with ID and Icon")
	return nil
}

// GetCategories - отримує список всіх категорій для користувача
func GetCategories(userID int) ([]models.Category, error) {
	collection := db.GetCategoryCollection()
	filter := bson.M{"userID": userID} // Фільтруємо по userID
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Error finding categories for userID %v: %v", userID, err)
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

// GetCategoryByID - отримує категорію за її ID та userID
func GetCategoryByID(userID int, categoryID int) (*models.Category, error) {
	collection := db.GetCategoryCollection()
	filter := bson.M{"userID": userID, "categoryID": categoryID}
	var category models.Category
	err := collection.FindOne(context.TODO(), filter).Decode(&category)
	if err != nil {
		log.Printf("Error finding category by ID: %v", err)
		return nil, err
	}
	return &category, nil
}

// UpdateCategory - оновлює категорію
func UpdateCategory(userID int, categoryID int, updatedCategory models.Category) error {
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
	if updatedCategory.Icon != "" {
		update["$set"].(bson.M)["icon"] = updatedCategory.Icon
	}

	// Виконання оновлення
	if len(update["$set"].(bson.M)) > 0 {
		filter := bson.M{"userID": userID, "categoryID": categoryID}
		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return fmt.Errorf("Error updating category: %v", err)
		}
	}

	return nil
}

// DeleteCategory - видаляє категорію за її ID та userID
func DeleteCategory(userID int, categoryID int) error {
	collection := db.GetCategoryCollection()

	// Фільтр для видалення категорії по userID та categoryID
	filter := bson.M{"userID": userID, "categoryID": categoryID}

	// Видалення категорії з колекції
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Error deleting category: %v", err)
		return err
	}

	log.Println("Category deleted successfully")
	return nil
}

// GetAllCategoriesByUserIDLogic - логіка отримання категорій з бази даних
func GetAllCategoriesByUserIDLogic(userID int) ([]models.Category, error) {
	collection := db.GetCategoryCollection()

	// Формуємо фільтр для пошуку категорій за userID
	filter := bson.M{"userID": userID}

	// Виконуємо запит до колекції
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Error finding categories for userID %d: %v", userID, err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var categories []models.Category
	// Прочитуємо всі документи зі знайдених категорій
	for cursor.Next(context.TODO()) {
		var category models.Category
		if err := cursor.Decode(&category); err != nil {
			log.Printf("Error decoding category: %v", err)
			continue
		}
		categories = append(categories, category)
	}

	// Перевіряємо наявність помилок при курсорі
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	// Повертаємо всі категорії для конкретного користувача
	return categories, nil
}

// GetCategoriesByName - отримує всі категорії за назвою
func GetCategoriesByName(categoryName string) ([]models.Category, error) {
	collection := db.GetCategoryCollection()

	// Формуємо фільтр для пошуку категорій за їх ім'ям
	filter := bson.M{"name": categoryName}

	// Виконуємо запит до колекції
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Error finding categories by name %s: %v", categoryName, err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var categories []models.Category
	// Читаємо всі категорії, що відповідають фільтру
	for cursor.Next(context.TODO()) {
		var category models.Category
		if err := cursor.Decode(&category); err != nil {
			log.Printf("Error decoding category: %v", err)
			continue
		}
		categories = append(categories, category)
	}

	// Перевіряємо наявність помилок при курсорі
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	// Повертаємо всі знайдені категорії
	return categories, nil
}
