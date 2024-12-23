package models

type Category struct {
	CategoryID  int    `bson:"categoryID" json:"categoryID"`
	BudgetID    int    `bson:"budgetID" json:"budgetID"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}
