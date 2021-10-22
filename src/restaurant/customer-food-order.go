package restaurant

import (
	runDatabases "Resort/src/database"
	"Resort/src/message_broker"
	"Resort/src/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var isWebSocketOpen bool = false

func FoodOrderWebSocket(c *gin.Context) {
	//Upgrade get request to webSocket protocol
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("error get connection")
		log.Fatal(err)
	}

	isWebSocketOpen = true

	consumerName := "foodOrderConsumer" + c.Request.Header.Get("Sec-Websocket-Key")

	msgs, err := message_broker.MqChannel.Consume(
		message_broker.FoodOrderQueue.Name, // queue
		consumerName,                       // consumer
		true,                               // auto-ack
		false,                              // exclusive
		false,                              // no-local
		false,                              // no-wait
		nil,                                // args
	)
	message_broker.FailOnError(err, "Failed to register a consumer for foodOrder")

	go func() {
		for d := range msgs {
			var jsonify interface{}
			json.Unmarshal(d.Body, &jsonify)
			err = ws.WriteJSON(jsonify)
			// err = ws.WriteJSON(string(d.Body))
			if err != nil {
				log.Println("error write json for foodOrder")
			}
		}
	}()

	go func() {
		code, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("error read message for foodOrder", err)
		}
		log.Printf("Received a message for foodOrder: %s", message)
		log.Printf("Received a code for foodOrder: %v", code)

		if code == -1 {
			message_broker.MqChannel.Cancel(consumerName, false) // client page is closed, no need to keep filling the queue since on returning back to the page, updated data will be requested by http protocol.
		}
	}()

	// log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
}

func CustomerFoodOrder(c *gin.Context) {
	db := runDatabases.MongoDb
	ctx := *runDatabases.MongoCtxPtr
	collection := db.Database("resort").Collection("food_orders")

	if result, err := collection.Find(ctx, bson.M{}); err != nil {
		log.Printf("Get CustomerFoodOrder Mongo Database error: %v", err)
		c.Status(http.StatusInternalServerError)
	} else {
		var decodedResult []bson.M
		result.All(ctx, &decodedResult)

		c.JSON(http.StatusOK, decodedResult)
	}
}

func UpdateFoodOrderState(c *gin.Context) {
	var foodOrder *models.UpdateFoodOrder
	// Call BindJSON to bind the received JSON/BSON to struct
	if err := c.BindJSON(&foodOrder); err != nil {
		log.Printf("Cannot be bound: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(*foodOrder); err != nil {
		log.Printf("Validataion error: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}

	db := runDatabases.MongoDb
	ctx := *runDatabases.MongoCtxPtr
	foodOrderCollection := db.Database("resort").Collection("food_orders")

	objId, _ := primitive.ObjectIDFromHex(foodOrder.Id)

	if _, err := foodOrderCollection.UpdateByID(ctx, objId, bson.M{"$set": bson.M{"orderState": foodOrder.OrderState}}); err != nil {
		log.Printf("UpdateFoodOrderState UpdateById error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if foodOrder.Email != "" { // email being null means the customer is guest.
		usersCollection := db.Database("resort").Collection("users")
		if _, err := usersCollection.UpdateOne(
			ctx,
			bson.M{
				"profile.email": foodOrder.Email,
				"restaurant": bson.M{
					"$elemMatch": bson.M{"orderDate": foodOrder.OrderDate},
				},
			},
			bson.M{
				"$set": bson.M{
					"restaurant.$.orderState": foodOrder.OrderState,
				},
			},
		); err != nil {
			log.Printf("UpdateFoodOrderState UpdateOne error: %v", err)
			c.Status(http.StatusInternalServerError)
		}
	}
}

func GetCourierIndex() {
	reply, _ := runDatabases.RedisDb.Do("get", "courierIndex")
	if reply == nil {
		if _, err := runDatabases.RedisDb.Do("set", "courierIndex", "0"); err != nil {
			log.Fatal("Something occured in resetting courierIndex:", err)
		}
		log.Println("Courier index has been set to zero.")
	} else {
		log.Printf("Current courier index: %s", reply)
	}
}
