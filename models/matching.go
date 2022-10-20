package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserId struct{
	UserId			string `json: "userId"`
}

type MatchingCreate struct{
	ActivityId			primitive.ObjectID `bson:"activityId" json:"activityId"`
	Participant			[]string `json: "participant"`
} 

type Matching struct{
	ID    				primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	ActivityId			primitive.ObjectID `bson:"activityId" json:"activityId"`
	Participant			[]string `json: "participant"`
} 


type response struct {
    Status  int                    `json:"status"`
    Message string                 `json:"message"`
    Data    map[string]interface{} `json:"data"`
}