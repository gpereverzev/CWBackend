package models

import "time"

// Reminder - структура для нагадування для цілі
type Reminder struct {
	GoalID       int       `json:"goalID"`
	Reminder     string    `json:"reminder"`     // Текст нагадування
	ReminderDate time.Time `json:"reminderDate"` // Дата наступного нагадування
	Period       string    `json:"period"`       // Період (day, week, month)
}
