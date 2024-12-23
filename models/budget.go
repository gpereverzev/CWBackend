package models

type Budget struct {
	BudgetID       int     `bson:"budgetID" json:"budgetID"`
	UserID         int     `bson:"userID" json:"userID"`
	Name           string  `bson:"name" json:"name"`
	InitialBalance float64 `bson:"initialBalance" json:"initialBalance"`
	Limit          float64 `bson:"limit" json:"limit"`
	Period         string  `bson:"period" json:"period"`
}
