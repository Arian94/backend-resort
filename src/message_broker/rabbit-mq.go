package message_broker

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// const (
// 	prefetchLimit = 1000
// 	pollDuration  = 100 * time.Millisecond
// 	numConsumers  = 1

// 	reportBatchSize = 10000
// 	consumeDuration = time.Millisecond
// 	shouldLog       = false
// )

// const (
// 	numDeliveries = 100000000
// 	batchSize     = 10000
// )

var (
	MqChannel *amqp.Channel
)

const (
	BOOKING_QUEUE_NAME    = "bookings"
	FOOD_ORDER_QUEUE_NAME = "foodOrders"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func InitializeRabbitMq() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()

	MqChannel, err = conn.Channel()
	FailOnError(err, "Failed to open a channel")
	// defer ch.Close()

	// BookingQueue, err = MqChannel.QueueDeclare(
	// 	"bookings", // name
	// 	false,      // durable
	// 	false,      // delete when unused
	// 	false,      // exclusive
	// 	false,      // no-wait
	// 	nil,        // arguments
	// )
	// FailOnError(err, "Failed to declare the bookingQueue queue")

	// FoodOrderQueue, err = MqChannel.QueueDeclare(
	// 	"foodOrders", // name
	// 	false,        // durable
	// 	false,        // delete when unused
	// 	false,        // exclusive
	// 	false,        // no-wait
	// 	nil,          // arguments
	// )
	// FailOnError(err, "Failed to declare the foodOrderQueue queue")
}

func Producer(queueName string, message interface{}) {
	msg, _ := json.Marshal(message)

	err := MqChannel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})

	FailOnError(err, "Failed to publish a message")
}
