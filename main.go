package main

import (
	"os"

	"matchingService/configs"
	controllers "matchingService/controller"

	"github.com/gin-gonic/gin"

	"fmt"
	"log"
	"net/http"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
)

func serviceRegistryWithConsul() {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Println(err)
	}

	serviceID := "matching-service1"
	port, _ := strconv.Atoi(getPort()[1:len(getPort())])
	address := getHostname()

	registration := &consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    "matching-service",
		Port:    port,
		Address: getHostname(),
		Check: &consulapi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%v/check", address, port),
			Interval: "10s",
			Timeout:  "30s",
		},
	}

	regiErr := consul.Agent().ServiceRegister(registration)

	if regiErr != nil {
		log.Printf("Failed to register service: %s:%v ", address, port)
	} else {
		log.Printf("successfully register service: %s:%v", address, port)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func check(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Consul check")
}

func getPort() (port string) {
	port = os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	port = ":" + port
	return
}

func getHostname() (hostname string) {
	hostname, _ = os.Hostname()
	return
}

func main() {
	serviceRegistryWithConsul()

	PORT := os.Getenv("PORT")

	r := gin.Default()
	configs.ConnectDB()

	r.POST("/matching/:activityId", controllers.CreateMatching())
	r.DELETE("/matching/:matchingId", controllers.DeleteMatching())
	r.POST("/attendActivity/:activityId", controllers.AttendActivity())
	r.POST("/leaveActivity/:activityId", controllers.LeaveActivity())
	r.GET("/matching/:matchingId", controllers.GetMatching())
	r.GET("/getMatchingByActivity/:activityId", controllers.GetMatchingByActivity())
	r.GET("/getActivitiesByParticipant/:userId", controllers.GetMatchingByParticaipant())
	r.GET("/check", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"check": "ok",
		})

	})
	r.Run("172.31.86.56:" + PORT)
}
