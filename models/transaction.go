package models

type Transaction struct {
	TransactionID int     `bson:"transactionID" json:"transactionID"`
	UserID        int     `bson:"userID" json:"userID"`
	CategoryID    int     `bson:"categoryID" json:"categoryID"`
	Type          string  `bson:"type" json:"type"`
	Amount        float64 `bson:"amount" json:"amount"`
	Date          string  `bson:"date" json:"date"`
	Description   string  `bson:"description" json:"description"`
}
