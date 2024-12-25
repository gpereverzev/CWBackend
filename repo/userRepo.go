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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

func init() {
	userCollection = db.Client.Database("cashWiseDB").Collection("Users")
}

// Функція для отримання наступного значення ID
func GetNextUserID() (int, error) {
	// Оновлюємо лічильник в колекції Counters
	filter := bson.M{"_id": "userID"} // ми будемо використовувати "userID" як ідентифікатор лічильника
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
		return 0, fmt.Errorf("failed to get next userID: %v", err)
	}
	return result.Seq, nil
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

// CreateUser - створює нового користувача і додає налаштування за замовчуванням
func CreateUser(user models.User) (models.User, error) {
	// Отримуємо наступний userID
	userID, err := GetNextUserID()
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get next user ID: %v", err)
	}

	// Встановлюємо отриманий userID
	user.UserID = userID

	// Вставка нового користувача
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to insert user: %v", err)
	}

	// Створення початкових налаштувань для нового користувача
	settings := models.Settings{
		UserID:         userID,
		DarkTheme:      false,
		TermsCondition: false,
		Notification:   false,
	}

	// Додаємо документ у колекцію Settings
	_, err = db.GetSettingCollection().InsertOne(context.TODO(), settings)
	if err != nil {
		// Якщо виникає помилка, видаляємо користувача, щоб уникнути розсинхронізації
		_, deleteErr := userCollection.DeleteOne(context.TODO(), bson.M{"userID": userID})
		if deleteErr != nil {
			return models.User{}, fmt.Errorf(
				"failed to create settings and cleanup user: %v; delete error: %v", err, deleteErr,
			)
		}
		return models.User{}, fmt.Errorf("failed to create settings: %v", err)
	}

	return user, nil
}

// Функція для хешування пароля
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
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
		// Хешуємо пароль перед оновленням
		hashedPassword, err := hashPassword(updatedUser.Password)
		if err != nil {
			return err
		}
		update["$set"].(bson.M)["password"] = hashedPassword
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
func GetUserByID(userID int) (*models.User, error) {
	// Пошук користувача за userID
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"userID": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}

	return &user, nil
}

// GetUserAndSettings - отримує всю інформацію по користувачу та його налаштуванням за userID
func GetUserAndSettings(userID int) (map[string]interface{}, error) {
	// Отримуємо колекції для користувачів та налаштувань
	settingsCollection := db.GetSettingCollection()

	// Знаходимо користувача за userID
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"userID": userID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to retrieve user: %v", err)
	}

	// Знаходимо налаштування користувача
	var settings models.Settings
	err = settingsCollection.FindOne(context.TODO(), bson.M{"userID": userID}).Decode(&settings)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("settings not found")
		}
		return nil, fmt.Errorf("failed to retrieve settings: %v", err)
	}

	// Формуємо результат
	userData := map[string]interface{}{
		"userID":         user.UserID,
		"fullName":       user.FullName,
		"email":          user.Email,
		"profilePicture": user.ProfilePicture,
		"role":           user.Role,
		"darkTheme":      settings.DarkTheme,
		"termsCondition": settings.TermsCondition,
	}

	return userData, nil
}
