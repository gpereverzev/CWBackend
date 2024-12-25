package models

type Goal struct {
	GoalID       int     `bson:"goalID" json:"goalID"`
	UserID       int     `bson:"userID" json:"userID"`
	Name         string  `bson:"name" json:"name"`
	TargetAmount float64 `bson:"targetAmount" json:"targetAmount"`
	Deadline     string  `bson:"deadline" json:"deadline"`
	Status       string  `bson:"status" json:"status"`
}
