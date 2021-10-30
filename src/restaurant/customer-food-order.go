package restaurant

import (
	runDatabases "Resort/src/database"
	"Resort/src/message_broker"
	"Resort/src/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/go-playground/validator.v9"
)

var (
	upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// isWebSocketOpen bool = false
	CourierIndex interface{}
)

func FoodOrderWebSocket(c *gin.Context) {
	//Upgrade get request to webSocket protocol
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("error get connection")
		log.Fatal(err)
	}

	authHeader := c.GetHeader("Sec-Websocket-Protocol")
	token, _ := jwt.Parse(authHeader, nil)
	claims := token.Claims.(jwt.MapClaims)
	emailFromToken := fmt.Sprintf("%v", claims["email"])
	log.Println("Restaurant Admin Email:", emailFromToken)

	// isWebSocketOpen = true

	FoodOrderQueue, err := message_broker.MqChannel.QueueDeclare(
		message_broker.FOOD_ORDER_QUEUE_NAME, // name
		false,                                // durable
		false,                                // delete when unused
		false,                                // exclusive
		false,                                // no-wait
		nil,                                  // arguments
	)
	message_broker.FailOnError(err, "Failed to declare the foodOrderQueue queue")

	message_broker.MqChannel.Cancel(emailFromToken, false) // for dev purposes

	msgs, err := message_broker.MqChannel.Consume(
		FoodOrderQueue.Name, // queue
		emailFromToken,      // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	message_broker.FailOnError(err, "Failed to register a consumer for foodOrder")

	log.Println("man injam")

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
			message_broker.MqChannel.Cancel(emailFromToken, false) // client page is closed, no need to keep filling the queue since on returning back to the page, updated data will be requested by http protocol.
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

	foodOrderCollResult := foodOrderCollection.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"orderState": foodOrder.OrderState}})
	if foodOrderCollResult.Err() == mongo.ErrNoDocuments {
		log.Printf("UpdateFoodOrderState FindOneAndUpdate error: %s", foodOrderCollResult.Err())
		c.Status(http.StatusBadRequest)
		return
	}

	var res bson.M
	foodOrderCollResult.Decode(&res)
	jsonRes, _ := json.Marshal(res)
	var courierInfo struct {
		Fullname    string `json:"fullname"`
		PhoneNumber string `json:"phonenumber"`
		Email       string `json:"email"`
	}
	if foodOrder.OrderState == SENT {
		if err := runDatabases.PostgresDb.QueryRow("UPDATE courier SET orders_list = $1 WHERE id = $2 RETURNING fullname, phoneNumber, email", jsonRes, "3").
			Scan(&courierInfo.Fullname, &courierInfo.PhoneNumber, &courierInfo.Email); err != nil {
			log.Println("Postgres updateFoodOrder error:", err)
			c.Status(http.StatusInternalServerError)
			return
		} else {
			log.Println("Postgres Courier Info:", courierInfo)
			// notify the courier
			log.Println("wtf", res)
			message_broker.Producer(courierInfo.Email, res)
			c.JSON(http.StatusOK, courierInfo)
		}
	} else if foodOrder.OrderState == DELIVERED {
		// notify the cooking admin
		if _, err := runDatabases.PostgresDb.
			Exec("UPDATE courier SET orders_list = jsonb_set(orders_list, '{orderState}', to_jsonb($1::text), true)  WHERE id = $2", DELIVERED, "3"); err != nil {
			log.Println("Postgres updateFoodOrder error:", err)
			c.Status(http.StatusInternalServerError)
			return
		} else {
			log.Println("Postgres Courier Order has been delivered.")
			// notify the courier
			message_broker.Producer(message_broker.FOOD_ORDER_QUEUE_NAME, []string{foodOrder.Id, foodOrder.OrderState})
		}
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
					"restaurant.$.orderState":  foodOrder.OrderState,
					"restaurant.$.courierInfo": courierInfo,
				},
			},
		); err != nil {
			log.Printf("UpdateFoodOrderState UpdateOne error: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}
}

func GetCourierIndex() {
	CourierIndex, _ = runDatabases.RedisDb.Do("get", "courierIndex")
	if CourierIndex == nil {
		if _, err := runDatabases.RedisDb.Do("set", "courierIndex", 1); err != nil {
			log.Fatal("Something occured in resetting courierIndex:", err)
		}
		log.Println("Courier index has been set to zero.")
	} else {
		log.Printf("Current courier index: %s", CourierIndex)
	}
}

func CourierWebSocket(c *gin.Context) {
	authHeader := c.GetHeader("Sec-Websocket-Protocol")
	token, _ := jwt.Parse(authHeader, nil)
	claims := token.Claims.(jwt.MapClaims)
	emailFromToken := fmt.Sprintf("%v", claims["email"])
	log.Println("Courier Email:", emailFromToken)

	// Upgrade get request to webSocket protocol
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("error get connection")
		log.Fatal(err)
	}

	_, err = message_broker.MqChannel.QueueDeclare(
		emailFromToken, // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	message_broker.FailOnError(err, "Failed to declare the courier queue")

	message_broker.MqChannel.Cancel(emailFromToken, false) // for dev purposes

	msgs, err := message_broker.MqChannel.Consume(
		emailFromToken, // queue
		emailFromToken, // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	message_broker.FailOnError(err, "Failed to register a consumer for courier")

	go func() {
		for d := range msgs {
			var jsonify interface{}
			json.Unmarshal(d.Body, &jsonify)
			log.Println("msg is:", jsonify)
			err = ws.WriteJSON(jsonify)
			// err = ws.WriteJSON(string(d.Body))
			if err != nil {
				log.Println("Error write json for courier")
			}
		}
	}()

	go func() {
		for {
			code, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("Error read message for courier", err)
			}
			log.Printf("Received a message for courier: %s", message)
			log.Printf("Received a code for courier: %v", code)

			if code == -1 {
				message_broker.MqChannel.Cancel(emailFromToken, false) // courier app is closed, no need to keep filling the queue since on returning back to the page, updated data will be requested by http protocol.
			}
		}
	}()

	// log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
}
