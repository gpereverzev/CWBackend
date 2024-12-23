package models

type User struct {
	UserID         int    `bson:"userID" json:"userID"`
	FullName       string `bson:"fullName" json:"fullName"`
	Email          string `bson:"email" json:"email"`
	Password       string `bson:"password" json:"password"`
	ProfilePicture string `bson:"profilePicture" json:"profilePicture"`
	Role           string `bson:"role" json:"role"`
}
