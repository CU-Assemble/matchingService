package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"matchingService/configs"
	"matchingService/models"
	"matchingService/responses"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var matchingCollection *mongo.Collection = configs.GetCollection(configs.DB, "matching")
var validate = validator.New()

func CreateMatching() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var userId models.UserId
		activityId,_ := primitive.ObjectIDFromHex(c.Param("activityId"))
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&userId); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&userId); validationErr != nil {
		    c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
		    return
		}
		fmt.Println("Create Matching")
		newMatching := models.MatchingCreate{
			ActivityId	: activityId,
			Participant : []string{userId.UserId},
		}
	
		result, err := matchingCollection.InsertOne(ctx, newMatching)
		if err != nil {
		    c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		    return
		}
	
		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})

	}
}

func DeleteMatching() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		objId, _ := primitive.ObjectIDFromHex(c.Param("matchingId"))
		defer cancel()

		result, err := matchingCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}})
	}
}

func AttendActivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// activityId := c.Param("activityId")
		// var activity models.Activity
		// defer cancel()

		// objId, _ := primitive.ObjectIDFromHex(activityId)

		// //validate the request body
		// if err := c.BindJSON(&activity); err != nil {
		// 	c.JSON(http.StatusBadRequest, responses.ActivityResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		// 	return
		// }

		// //use the validator library to validate required fields
		// if validationErr := validate.Struct(&activity); validationErr != nil {
		// 	c.JSON(http.StatusBadRequest, responses.ActivityResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
		// 	return
		// }

		// update := bson.M{
		//     "Name":      		activity.Name,
		// 	"Description": 	activity.Description,
		// 	"ImageProfile":	activity.ImageProfile,
		// 	"Type":			activity.Type,
		// 	"OwnerId": 		activity.OwnerId,
		// 	"Location":		activity.Location,
		// 	"MaxParticipant":	activity.MaxParticipant,
		// 	"Participant":    activity.Participant,
		// 	"Date":			activity.Date,
		// 	"Duration":		activity.Duration,
		// 	"ChatId":			activity.ChatId,
		// }

		// result, err := activityCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, responses.ActivityResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		// 	return
		// }

		// //get updated user details
		// var updatedActivity models.Activity
		// if result.MatchedCount == 1 {
		// 	err := activityCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedActivity)

		// 	if err != nil {
		// 		c.JSON(http.StatusInternalServerError, responses.ActivityResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		// 		return
		// 	}
		// }

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "updatedActivity"}})
	}
}

func LeaveActivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// activityId := c.Param("activityId")
		// defer cancel()

		// objId, _ := primitive.ObjectIDFromHex(activityId)

		// result, err := activityCollection.DeleteOne(ctx, bson.M{"_id": objId})

		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, responses.ActivityResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		// 	return
		// }

		// if result.DeletedCount < 1 {
		// 	c.JSON(http.StatusNotFound,
		// 		responses.ActivityResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}},
		// 	)
		// 	return
		// }

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}})
	}
}

func GetMatching() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// var activitys []models.Activity
		// defer cancel()

		// results, err := activityCollection.Find(ctx, bson.M{})

		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, responses.ActivityResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		// 	return
		// }

		// //reading from the db in an optimal way
		// defer results.Close(ctx)
		// for results.Next(ctx) {
		// 	var oneActivity models.Activity
		// 	if err = results.Decode(&oneActivity); err != nil {
		// 		c.JSON(http.StatusInternalServerError, responses.ActivityResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		// 	}

		// 	activitys = append(activitys, oneActivity)
		// }

		c.JSON(http.StatusOK,
			responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "activitys"}},
		)
	}
}
