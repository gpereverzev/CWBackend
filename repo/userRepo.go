// repo/userRepo.go
package repo

import (
	"cashWise/db"
	"cashWise/models"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func init() {
	userCollection = db.Client.Database("cashWiseDB").Collection("Users")
}

// Функція для отримання користувача за email
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// Функція для створення нового користувача
func CreateUser(user models.User) (models.User, error) {
	// Вставка нового користувача
	_, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// Функція для оновлення даних користувача
func UserUpdate(userID int, updatedUser models.User) error {
	// Оновлення даних користувача за userID
	update := bson.M{
		"$set": bson.M{},
	}

	// Перевіряємо і додаємо лише ті поля, які не є пустими
	if updatedUser.FullName != "" {
		update["$set"].(bson.M)["fullName"] = updatedUser.FullName
	}
	if updatedUser.Email != "" {
		update["$set"].(bson.M)["email"] = updatedUser.Email
	}
	if updatedUser.Password != "" {
		update["$set"].(bson.M)["password"] = updatedUser.Password
	}
	if updatedUser.ProfilePicture != "" {
		update["$set"].(bson.M)["profilePicture"] = updatedUser.ProfilePicture
	}
	if updatedUser.Role != "" {
		update["$set"].(bson.M)["role"] = updatedUser.Role
	}

	// Виконуємо оновлення, якщо хоча б одне поле було вказано
	if len(update["$set"].(bson.M)) > 0 {
		_, err := userCollection.UpdateOne(context.TODO(), bson.M{"userID": userID}, update)
		if err != nil {
			return err
		}
	}

	return nil
}

// Функція для видалення користувача за userID
func UserDelete(userID int) error {
	// Видаляємо користувача з колекції
	_, err := userCollection.DeleteOne(context.TODO(), bson.M{"userID": userID})
	if err != nil {
		return err
	}
	return nil
}

// GetUserByID - Отримання користувача за його ID з бази даних
func GetUserByID(id string) (*models.User, error) {
	// Перетворюємо ID у BSON ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Пошук користувача за ID
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = userCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}

	return &user, nil
}
