package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)


type Activity struct{
	ID    				primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name 				string `json:"activityName" validate:"required"`
	Description 		string  `json: "description"`
	// ImageProfile 		string `json: "imageProfile"`
	OwnerId 			string `json: "ownerId" validate:"required"`
	Location			string `json: "location" validate:"required"`
	MaxParticipant		int `json: "maxParticipant" validate:"required"`
	Participant			[]string `json: "participant"`
	Date				time.Time `json: "date"`
	Duration			float32 `json: "duration"` 
	ChatId   			string `json: "chatId"`
} 

type response struct {
    Status  int                    `json:"status"`
    Message string                 `json:"message"`
    Data    map[string]interface{} `json:"data"`
}