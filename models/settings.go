package models

type Settings struct {
	UserID         int  `bson:"userID"`
	DarkTheme      bool `bson:"darkTheme"`
	TermsCondition bool `bson:"termsConditional"`
	Notification   bool `bson:"notification"`
}
