package models

type Category struct {
	CategoryID  int    `bson:"categoryID" json:"categoryID"`
	UserID      int    `bson:"userID" json:"userID"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Icon        string `bson:"icon" json:"icon"`
}
