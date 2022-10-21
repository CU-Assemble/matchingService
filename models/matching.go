package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type UserId struct {
	UserId string `json: "userId"`
}

type MatchingCreate struct {
	ActivityId  primitive.ObjectID `bson:"activityId" json:"activityId"`
	Participant []string           `json: "participant"`
}

type Matching struct {
	ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	ActivityId  primitive.ObjectID `bson:"activityId" json:"activityId"`
	Participant []string           `json: "participant"`
}

type UserFull struct {
	Detail User `json:"user"`
}
type User struct {
	// gorm.Model
	StudentId string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
	Nickname  string
	Faculty   string
	Tel       string
	Email     string
	Password  string
}

type response struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type Activity struct {
	ActivityId     string    `json:"ActivityId"`
	Name           string    `json:"Name"`
	Description    string    `json:"Description"`
	ActivityType   []string  `json:"ActivityType"`
	ImageProfile   string    `json:"ImageProfile"`
	OwnerId        string    `json:"OwnerId"`
	Location       string    `json:"Location"`
	MaxParticipant int       `json:"MaxParticipant"`
	Participant    string    `json:"Participant"`
	Date           time.Time `json:"Date"`
	Duration       float32   `json:"Duration"`
	ChatId         string    `json:"ChatId"`
}
