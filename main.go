package main

import (
	// "fmt"
	"context"
	"log"
	
	"strings"
	"time"

	"matchingService/services"
	// "net/http"

	// "github.com/gin-gonic/gin"

	amqp "github.com/rabbitmq/amqp091-go"
)


func failOnError(err error, msg string) {
	if err != nil {
			log.Panicf("%s: %s", msg, err)
	}
}

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			for d := range msgs {
					n := strings.Fields(string(d.Body))

					serviceType := n[0]
					
					result,err := "none",error(nil)

					switch serviceType{
					case "create":
						result,err = services.CreateMatching(n[1])
					case "delete":
						result,err = services.DeleteMatching(n[1])
					case "attend":
						result,err = services.AttendActivity(n[1],n[2])
					case "leave":
						result,err = services.LeaveActivity(n[1],n[2])
					case "get":
						result,err = services.GetMatching(n[1])
					default:
						result,err = "serviceType error",error(nil)
					}

					failOnError(err, "Failed to do a services.")

					err = ch.PublishWithContext(ctx,
							"",        // exchange
							d.ReplyTo, // routing key
							false,     // mandatory
							false,     // immediate
							amqp.Publishing{
									ContentType:   "text/plain",
									CorrelationId: d.CorrelationId,
									Body:          []byte(result),
							})
					failOnError(err, "Failed to publish a message")

					d.Ack(false)
			}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever

	// r := gin.New()
	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "Hello World!",
	// 	})
	// })
	// r.Run()



	// fmt.Println("Successfully Connected to our RabbitMQ Instance")
}