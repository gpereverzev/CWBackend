package handlers

import (
	"cashWise/models"
	"cashWise/repo"
	"cashWise/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Define a custom type for the context key
type ContextKey string

// Declare key constants using the custom type
const (
	RoleKey ContextKey = "role"
)

// Реєстрація нового користувача
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Перевірка, чи вже існує користувач з таким email
	existingUser, _ := repo.GetUserByEmail(user.Email)
	if existingUser.Email != "" {
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	// Хешування пароля
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	// Додавання користувача в БД
	result, err := repo.CreateUser(user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// Авторизація користувача
func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Обработка preflight-запроса
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Парсинг тела запроса
	var user models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Получение пользователя из базы данных
	userFromDB, err := repo.GetUserByEmail(user.Email)
	if err != nil || !utils.CheckPasswordHash(user.Password, userFromDB.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Формирование JSON-ответа
	response := map[string]interface{}{
		"message": "Login successful",
		"userID":  userFromDB.UserID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Перевірка ролі користувача
func CheckRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Отримуємо роль із контексту запиту
			role, ok := r.Context().Value(RoleKey).(string)
			if !ok || role == "" {
				http.Error(w, "Role not found", http.StatusUnauthorized)
				return
			}
			// Перевіряємо чи роль співпадає з необхідною
			if role != requiredRole {
				http.Error(w, "Access Denied", http.StatusForbidden)
				return
			}
			// Якщо все ок - викликаємо наступний хендлер
			next.ServeHTTP(w, r)
		})
	}
}

// Мідлвар для "підробного" користувача (для тестування)
func MockUser(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Додаємо роль користувача до контексту запиту
			ctx := context.WithValue(r.Context(), RoleKey, role)
			// Викликаємо наступний хендлер з оновленим контекстом
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Отримуємо роль з query параметрів
	role := r.URL.Query().Get("role")
	fmt.Println("Role from query:", role) // Додаємо для дебагу
	if role == "" {
		http.Error(w, "Role not found", http.StatusUnauthorized)
		return
	}

	// Тільки користувачі з роллю "admin" можуть оновлювати дані інших користувачів
	if role != "admin" {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	// Отримуємо userID з параметрів запиту та конвертуємо його в int
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	userID, err := strconv.Atoi(userIDStr) // Перетворюємо рядок в ціле число
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Отримуємо дані користувача, якому потрібно оновити інформацію
	var updatedUser models.User
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Оновлення користувача в БД
	err = repo.UserUpdate(userID, updatedUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not update user: %v", err), http.StatusInternalServerError)
		return
	}

	// Повертаємо успішну відповідь
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User data updated"})
}

// DeleteUser - Видалення акаунту користувача
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Отримуємо роль з query параметрів
	role := r.URL.Query().Get("role")
	fmt.Println("Role from query:", role) // Додаємо для дебагу
	if role == "" {
		http.Error(w, "Role not found", http.StatusUnauthorized)
		return
	}

	// Тільки користувачі з роллю "admin" можуть видаляти дані інших користувачів
	if role != "admin" {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	// Отримуємо userID з параметрів запиту та конвертуємо його в int
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	userID, err := strconv.Atoi(userIDStr) // Перетворюємо рядок в ціле число
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Видалення користувача з БД
	err = repo.UserDelete(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete user: %v", err), http.StatusInternalServerError)
		return
	}

	// Повертаємо успішну відповідь
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User account deleted"})
}

// GetUserByID - Отримання користувача за його ID
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Отримуємо userID з параметрів запиту
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	// Отримуємо користувача з бази даних
	user, err := repo.GetUserByID(userIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not retrieve user: %v", err), http.StatusInternalServerError)
		return
	}

	// Якщо користувач не знайдений
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Повертаємо успішну відповідь із даними користувача
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
