package models

type Settings struct {
	UserID    int  `bson:"userID"`
	DarkTheme bool `bson:"darkTheme"`
}
