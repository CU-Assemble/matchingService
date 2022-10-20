package main

import (
	"os"

	"matchingService/configs"
	controllers "matchingService/controller"

	"github.com/gin-gonic/gin"
)

func main() {

	PORT := os.Getenv("PORT")

	r := gin.Default()
	configs.ConnectDB()
	
	r.POST("/matching/:activityId", controllers.CreateMatching())
	r.DELETE("/matching/:matchingId", controllers.DeleteMatching())
	r.PUT("/matching/attend/:matchingId", controllers.AttendActivity())
	r.PUT("/matching/leave/:matchingId", controllers.LeaveActivity())
	r.GET("/matching/:matchingId", controllers.GetMatching())

	r.Run("localhost:"+PORT)

}