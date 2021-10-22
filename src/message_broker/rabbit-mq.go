package message_broker

import (
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

var MqChannel *amqp.Channel
var Queue amqp.Queue
var FoodOrderQueue amqp.Queue

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

	Queue, err = MqChannel.QueueDeclare(
		"bookings", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	FailOnError(err, "Failed to declare the bookingQueue queue")

	FoodOrderQueue, err = MqChannel.QueueDeclare(
		"foodOrders", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	FailOnError(err, "Failed to declare the foodOrderQueue queue")
}

func BookingProducer(newBooking []byte) {
	err := MqChannel.Publish(
		"",         // exchange
		Queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        newBooking,
		})
	FailOnError(err, "Failed to publish a message")
}

func FoodOrderProducer(order []byte) {
	err := MqChannel.Publish(
		"",                  // exchange
		FoodOrderQueue.Name, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        order,
		})
	FailOnError(err, "Failed to publish a message")

}
