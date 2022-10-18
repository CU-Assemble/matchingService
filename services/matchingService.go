package services

import (
	// "context"
	"net/http"
	// "time"

	"matchingService/configs"
	"matchingService/models"
	"matchingService/responses"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var matchingCollection *mongo.Collection = configs.GetCollection(configs.DB, "matching")

func CreateMatching(id string)(string,error){}

func DeleteMatching(id string)(string,error){}

func AttendActivity(id string,userId string)(string,error){}

func LeaveActivity(id string,userId string)(string,error){}

func GetMatching(id string)(string,error){}


