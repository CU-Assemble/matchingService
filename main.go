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
	r.PUT("/matching/attend/:activityId", controllers.AttendActivity())
	r.PUT("/matching/leave/:activityId", controllers.LeaveActivity())
	r.GET("/matching/:matchingId", controllers.GetMatching())
	r.GET("/getMatchingByActivity/:activityId", controllers.GetMatchingByActivity())
	r.GET("/getActivitiesByParticipant/:userId", controllers.GetMatchingByParticaipant())
	r.Run("localhost:" + PORT)

}
