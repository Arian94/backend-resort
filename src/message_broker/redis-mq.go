package message_broker

import (
	"log"

	"github.com/streadway/amqp"
)

// var redisBookingProducer rmq.Queue
// var ConsumerQ = NewConsumer(0)

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
var ch *amqp.Channel
var q amqp.Queue

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func InitializeRabbitMq() {
	// connection, err := rmq.OpenConnection("producer", "tcp", "127.0.0.1:6379", 1, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// if redisBookingProducer, err = connection.OpenQueue("bookings"); err != nil {
	// 	panic(err)
	// } else {
	// 	log.Println("redis connected")
	// }
	// var err error
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()

	ch, err = conn.Channel()
	FailOnError(err, "Failed to open a channel")
	// defer ch.Close()

	q, err = ch.QueueDeclare(
		"bookings", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	FailOnError(err, "Failed to declare a queue")

}

func Producer(newBooking []byte) {
	// body := "Hello World!"
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        newBooking,
		})
	FailOnError(err, "Failed to publish a message")
	// log.Printf(" [x] Sent %s", newBooking)
	// var before time.Time
	// for i := 0; i < numDeliveries; i++ {
	// delivery := fmt.Sprintf("delivery %d", i)
	// if err := redisBookingProducer.Publish(newBooking); err != nil {
	// 	log.Printf("failed to publish: %s", err)
	// }

	// if i%batchSize == 0 {
	// 	duration := time.Since(time.Time{})
	// 	before = time.Now()
	// 	perSecond := time.Second / (duration / batchSize)
	// 	log.Printf("produced %d %s %d", i, delivery, perSecond)
	// }
}

// func Receiver() {
// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// failOnError(err, "Failed to connect to RabbitMQ")
// // defer conn.Close()

// ch, err := conn.Channel()
// failOnError(err, "Failed to open a channel")
// // defer ch.Close()

// q, err := ch.QueueDeclare(
// 	"hello", // name
// 	false,   // durable
// 	false,   // delete when unused
// 	false,   // exclusive
// 	false,   // no-wait
// 	nil,     // arguments
// )
// failOnError(err, "Failed to declare a queue")

// msgs, err := ch.Consume(
// 	q.Name, // queue
// 	"",     // consumer
// 	true,   // auto-ack
// 	false,  // exclusive
// 	false,  // no-local
// 	false,  // no-wait
// 	nil,    // args
// )
// failOnError(err, "Failed to register a consumer")

// // return msgs

// forever := make(chan bool)

// go func() {
// 	for d := range msgs {
// 		// return d.Body
// 		log.Printf("Received a message: %s", d.Body)
// 	}
// }()

// log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
// <-forever
// }

// func InitConsumer() {
// errChan := make(chan error, 10)
// go logErrors(errChan)

// connection, err := rmq.OpenConnection("consumer", "tcp", "localhost:6379", 1, errChan)
// if err != nil {
// 	panic(err)
// }

// queue, err := connection.OpenQueue("bookings")
// if err != nil {
// 	panic(err)
// }

// if err := queue.StartConsuming(prefetchLimit, pollDuration); err != nil {
// 	panic(err)
// }

// // bookingConsumer := &Consumer{}

// // for i := 0; i < numConsumers; i++ {
// // name := fmt.Sprintf("consumer %d", i)
// if _, err := queue.AddConsumer("name", ConsumerQ); err != nil {
// 	panic(err)
// }
// // }

// log.Println("Bookings Consumer Connected!")

// signals := make(chan os.Signal, 1)
// signal.Notify(signals, syscall.SIGINT)
// defer signal.Stop(signals)

// <-signals // wait for signal
// go func() {
// 	<-signals // hard exit on second signal (in case shutdown gets stuck)
// 	os.Exit(1)
// }()

// }

// type Consumer struct {
// 	name   string
// 	count  int
// 	before time.Time
// }

// func NewConsumer(tag int) *Consumer {
// 	return &Consumer{
// 		name:   fmt.Sprintf("consumer%d", tag),
// 		count:  0,
// 		before: time.Now(),
// 	}
// }

// func (consumer *Consumer) Consume(delivery rmq.Delivery) {
// 	payload := delivery.Payload()
// 	// debugf("start consume %s", payload)
// 	// time.Sleep(consumeDuration)

// 	// return payload

// 	consumer.count++
// 	if consumer.count%reportBatchSize == 0 {
// 		duration := time.Since(consumer.before)
// 		consumer.before = time.Now()
// 		perSecond := time.Second / (duration / reportBatchSize)
// 		log.Printf("%s consumed %d %s %d", consumer.name, consumer.count, payload, perSecond)
// 	}

// if consumer.count%reportBatchSize > 0 {
// 	if err := delivery.Ack(); err != nil {
// 		debugf("failed to ack %s: %s", payload, err)
// 	} else {
// 		debugf("acked %s", payload)
// 	}
// } else { // reject one per batch
// 	if err := delivery.Reject(); err != nil {
// 		debugf("failed to reject %s: %s", payload, err)
// 	} else {
// 		debugf("rejected %s", payload)
// 	}
// }
// }

// func logErrors(errChan <-chan error) {
// 	for err := range errChan {
// 		switch err := err.(type) {
// 		case *rmq.HeartbeatError:
// 			if err.Count == rmq.HeartbeatErrorLimit {
// 				log.Print("heartbeat error (limit): ", err)
// 			} else {
// 				log.Print("heartbeat error: ", err)
// 			}
// 		case *rmq.ConsumeError:
// 			log.Print("consume error: ", err)
// 		case *rmq.DeliveryError:
// 			log.Print("delivery error: ", err.Delivery, err)
// 		default:
// 			log.Print("other error: ", err)
// 		}
// 	}
// }
