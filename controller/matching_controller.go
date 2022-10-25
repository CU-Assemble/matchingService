package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
		activityId, _ := primitive.ObjectIDFromHex(c.Param("activityId"))
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
			ActivityId:  activityId,
			Participant: []string{userId.UserId},
		}

		result, err := matchingCollection.InsertOne(ctx, newMatching)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result.InsertedID}})

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
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var userId models.UserId
		activityId, _ := primitive.ObjectIDFromHex(c.Param("activityId"))
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&userId); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error binding", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&userId); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error validate", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		fmt.Println("attending user at act ID", activityId)
		filter := bson.D{{"activityId", activityId}}
		change := bson.M{"$push": bson.M{"participant": userId.UserId}}

		result, err := matchingCollection.UpdateOne(ctx, filter, change)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func LeaveActivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var userId models.UserId
		activityId, _ := primitive.ObjectIDFromHex(c.Param("activityId"))
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&userId); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error binding", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&userId); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error validate", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		fmt.Println("attending user at act ID", activityId)
		filter := bson.D{{"activityId", activityId}}
		change := bson.M{"$pull": bson.M{"participant": userId.UserId}}

		result, err := matchingCollection.UpdateOne(ctx, filter, change)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetMatching() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		id, _ := primitive.ObjectIDFromHex(c.Param("matchingId"))
		defer cancel()
		fmt.Println(c.Param("matchingId"))
		var matching models.Matching
		var matchingDetail MatchingFullDetail
		result := matchingCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&matching)
		if result != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": result.Error()}})
			return
		}
		matchingDetail.MatchingId = matching.ID.String()
		var myClient = &http.Client{Timeout: 10 * time.Second}
		resp, err := myClient.Get("http://localhost:8000/activity/" + matching.ActivityId.Hex())
		fmt.Println("http://localhost:8000/activity/" + matching.ActivityId.Hex())
		if err != nil {
			return
		}
		defer resp.Body.Close()
		act := models.Activity{}
		b, err := io.ReadAll(resp.Body)
		fmt.Println(string(b))
		json.Unmarshal(b, &act)

		matchingDetail.Activity.ActivityId = act.ActivityId
		matchingDetail.Activity.Name = act.Name
		matchingDetail.Activity.Description = act.Description
		matchingDetail.Activity.ActivityType = []string{}
		matchingDetail.Activity.ActivityType = append(matchingDetail.Activity.ActivityType, act.ActivityType...)
		fmt.Println(act.ActivityType, "111")
		matchingDetail.Activity.ImageProfile = act.ImageProfile
		matchingDetail.Activity.OwnerId = act.OwnerId
		matchingDetail.Activity.Location = act.Location
		matchingDetail.Activity.MaxParticipant = act.MaxParticipant
		matchingDetail.Activity.Date = act.Date
		matchingDetail.Activity.Duration = act.Duration
		matchingDetail.Activity.ChatId = act.ChatId

		for i, par := range matching.Participant {
			fmt.Print(i)
			resp, err := myClient.Get("http://localhost:3000/user/" + par)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			user := models.UserFull{}
			json.NewDecoder(resp.Body).Decode(&user)
			fmt.Println(string(user.Detail.Name))

			userDetail := UserDetail{
				StudentId: user.Detail.StudentId,
				Name:      user.Detail.Name,
				Nickname:  user.Detail.Nickname,
				Faculty:   user.Detail.Faculty,
				Tel:       user.Detail.Tel,
				Email:     user.Detail.Email,
			}

			matchingDetail.Users = append(matchingDetail.Users, userDetail)
		}
		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": matchingDetail}})
	}
}

func GetMatchingByActivity() gin.HandlerFunc {
	return func(c *gin.Context) {

		//filter matching by actId
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		activityId, _ := primitive.ObjectIDFromHex(c.Param("activityId"))
		defer cancel()
		fmt.Println(c.Param("activityId"))
		var matching models.Matching
		var matchingDetail MatchingDetail
		result := matchingCollection.FindOne(ctx, bson.M{"activityId": activityId}).Decode(&matching)
		if result != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": result.Error()}})
			return
		}
		matchingDetail.MatchingId = matching.ID.String()
		matchingDetail.ActivityId = matching.ActivityId.String()
		for i, par := range matching.Participant {
			fmt.Print(i)
			var myClient = &http.Client{Timeout: 10 * time.Second}
			resp, err := myClient.Get("http://localhost:3000/user/" + par)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			user := models.UserFull{}
			json.NewDecoder(resp.Body).Decode(&user)
			fmt.Println(string(user.Detail.Name))

			userDetail := UserDetail{
				StudentId: user.Detail.StudentId,
				Name:      user.Detail.Name,
				Nickname:  user.Detail.Nickname,
				Faculty:   user.Detail.Faculty,
				Tel:       user.Detail.Tel,
				Email:     user.Detail.Email,
			}

			matchingDetail.Users = append(matchingDetail.Users, userDetail)
		}
		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": matchingDetail}})
	}
}

type Activity struct {
	ActivityId     string   `json:"ActivityId"`
	Name           string   `json:"Name" validate:"required"`
	Description    string   `json: "Description"`
	ActivityType   []string `json: "ActivityType"`
	ImageProfile   string   `json: "ImageProfile"`
	OwnerId        string   `json: "OwnerId" validate:"required"`
	Location       string   `json: "Location" validate:"required"`
	MaxParticipant int      `json: "MaxParticipant" validate:"required"`
	Date           string   `json: "Date"`
	Duration       float32  `json: "Duration"`
	ChatId         string   `json: "ChatId"`
}

type MatchingFullDetail struct {
	Activity   Activity `json:"Activity"`
	MatchingId string
	Users      []UserDetail `json:"ParticipantId"`
}

type MatchingDetail struct {
	ActivityId string
	MatchingId string
	Users      []UserDetail `json:"ParticipantId"`
}

type UserDetail struct {
	StudentId string
	Name      string
	Nickname  string
	Faculty   string
	Tel       string
	Email     string
}

func GetMatchingByParticaipant() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		//userId, _ := primitive.ObjectIDFromHex(c.Param("userId"))
		defer cancel()
		fmt.Println(c.Param("userId"))
		// var matchings[] models.Matching
		// var matchingDetail[] MatchingDetail
		cursor, err := matchingCollection.Find(ctx, bson.M{"participant": c.Param("userId")})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		var results []models.Matching
		if err = cursor.All(ctx, &results); err != nil {
			panic(err)
		}
		if err := cursor.Close(ctx); err != nil {
			panic(err)
		}
		fmt.Println(results)
		var matchingFullDetails []MatchingFullDetail
		for a, matching := range results {
			fmt.Println(a)
			var matchingDetail MatchingFullDetail
			matchingDetail.MatchingId = matching.ID.String()
			var myClient = &http.Client{Timeout: 10 * time.Second}
			resp, err := myClient.Get("http://localhost:8000/activity/" + matching.ActivityId.Hex())
			fmt.Println("http://localhost:8000/activity/" + matching.ActivityId.Hex())
			if err != nil {
				return
			}
			defer resp.Body.Close()
			act := models.Activity{}
			b, err := io.ReadAll(resp.Body)
			fmt.Println(string(b))
			json.Unmarshal(b, &act)

			matchingDetail.Activity.ActivityId = act.ActivityId
			matchingDetail.Activity.Name = act.Name
			matchingDetail.Activity.Description = act.Description
			matchingDetail.Activity.ActivityType = []string{}
			matchingDetail.Activity.ActivityType = append(matchingDetail.Activity.ActivityType, act.ActivityType...)
			matchingDetail.Activity.ImageProfile = act.ImageProfile
			matchingDetail.Activity.OwnerId = act.OwnerId
			matchingDetail.Activity.Location = act.Location
			matchingDetail.Activity.MaxParticipant = act.MaxParticipant
			matchingDetail.Activity.Date = act.Date
			matchingDetail.Activity.Duration = act.Duration
			matchingDetail.Activity.ChatId = act.ChatId

			for i, par := range matching.Participant {
				fmt.Print(i)
				var myClient = &http.Client{Timeout: 10 * time.Second}
				resp, err := myClient.Get("http://localhost:3000/user/" + par)
				if err != nil {
					return
				}
				defer resp.Body.Close()
				user := models.UserFull{}
				json.NewDecoder(resp.Body).Decode(&user)
				fmt.Println(string(user.Detail.Name))

				userDetail := UserDetail{
					StudentId: user.Detail.StudentId,
					Name:      user.Detail.Name,
					Nickname:  user.Detail.Nickname,
					Faculty:   user.Detail.Faculty,
					Tel:       user.Detail.Tel,
					Email:     user.Detail.Email,
				}

				matchingDetail.Users = append(matchingDetail.Users, userDetail)
			}
			matchingFullDetails = append(matchingFullDetails, matchingDetail)
		}
		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": matchingFullDetails}})
	}
}
