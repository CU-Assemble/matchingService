package services

import (
	// "context"
	"context"
	"fmt"
	"time"

	// "time"

	"matchingService/configs"
	"matchingService/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var matchingCollection *mongo.Collection = configs.GetCollection(configs.DB, "matching")

func CreateMatching(id string)(string,error){

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	actId, _ := primitive.ObjectIDFromHex(id)

	newMatching := models.Matching{
		ActivityId	: actId,
		Participant : []string{},
	}

	result, err := matchingCollection.InsertOne(ctx, newMatching)
	if err != nil {
		return "Error Inserting data",err
	}

	ID := fmt.Sprintf("%v", result.InsertedID)

	return ID,nil

}

func DeleteMatching(id string)(string,error){
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(id)

	result, err := matchingCollection.DeleteOne(ctx, bson.M{"_id": objId})


	if err != nil {
		return "401 Error while deleting "+id,err
	}

	if result.DeletedCount < 1 {
		return "404 Matching id not found",nil
	}

	return "200 Sucessed",nil
}

// func AttendActivity(id string,userId string)(string,error){}

// func LeaveActivity(id string,userId string)(string,error){}

// func GetMatching(id string)(string,error){}


