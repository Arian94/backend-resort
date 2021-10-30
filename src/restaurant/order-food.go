package restaurant

import (
	runDatabases "Resort/src/database"
	"Resort/src/message_broker"
	"Resort/src/middleware"
	"Resort/src/models"
	signupLogin "Resort/src/signup-login"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/go-playground/validator.v9"
)

const (
	CANCELED        = "Canceled"
	ABSENT_CUSTOMER = "Absent Customer"
	RECEIVED        = "Received"
	IN_PROGRESS     = "In Progress"
	SENT            = "Sent"
	DELIVERED       = "Delivered"
)

func OrderFoods(c *gin.Context) {
	db := runDatabases.MongoDb
	ctx := *runDatabases.MongoCtxPtr
	var customerOrder *models.FoodOrder

	// Call BindJSON to bind the received JSON/BSON to struct
	if err := c.BindJSON(&customerOrder); err != nil {
		log.Printf("Cannot be bound: %v", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(*customerOrder); err != nil {
		log.Printf("Validataion error: %v", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	var emailFromToken string
	var id interface{}

	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		token, _ := signupLogin.JWTAuthService().ValidateToken(authHeader[len(middleware.BEARER_SCHEMA)+1:])
		claims := token.Claims.(jwt.MapClaims)
		emailFromToken = fmt.Sprintf("%v", claims["email"])

		userCollection := db.Database("resort").Collection("users")
		if _, err := userCollection.UpdateOne(
			ctx,
			bson.M{"profile.email": emailFromToken},
			bson.M{"$push": bson.M{
				"restaurant": bson.M{
					"orders":     customerOrder.Orders,
					"receiver":   customerOrder.CustomerForm.FullName,
					"totalPrice": customerOrder.TotalPrice,
					"orderDate":  time.Now().Format(time.RFC3339),
					"orderState": RECEIVED,
				},
			},
			},
		); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		foodOrderCollection := db.Database("resort").Collection("food_orders")
		result, err := foodOrderCollection.InsertOne(
			ctx,
			bson.M{
				"email":               emailFromToken,
				"receiver":            customerOrder.CustomerForm.FullName,
				"receiverPhoneNumber": customerOrder.CustomerForm.PhoneNumber,
				"orders":              customerOrder.Orders,
				"totalPrice":          customerOrder.TotalPrice,
				"address":             customerOrder.CustomerForm.Address,
				"orderDate":           time.Now().Format(time.RFC3339),
				"orderState":          RECEIVED,
			},
		)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		id = result.InsertedID
	} else {
		foodOrderCollection := db.Database("resort").Collection("food_orders")
		result, err := foodOrderCollection.InsertOne(
			ctx,
			bson.M{
				"receiver":            customerOrder.CustomerForm.FullName,
				"receiverPhoneNumber": customerOrder.CustomerForm.PhoneNumber,
				"orders":              customerOrder.Orders,
				"totalPrice":          customerOrder.TotalPrice,
				"address":             customerOrder.CustomerForm.Address,
				"orderDate":           time.Now().Format(time.RFC3339),
				"orderState":          RECEIVED,
			},
		)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		id = result.InsertedID
	}

	c.Status(http.StatusOK)

	log.Println("New food order is sent to queue")
	foodOrder := struct {
		Id                  interface{} `json:"_id"`
		Email               string      `json:"email"`
		Receiver            string      `json:"receiver"`
		ReceiverPhoneNumber string      `json:"receiverPhoneNumber"`
		Orders              interface{} `json:"orders"`
		TotalPrice          int64       `json:"totalPrice"`
		Address             string      `json:"address"`
		OrderDate           string      `json:"orderDate"`
		OrderState          string      `json:"orderState"`
	}{
		id,
		emailFromToken,
		customerOrder.CustomerForm.FullName,
		customerOrder.CustomerForm.PhoneNumber,
		customerOrder.Orders,
		customerOrder.TotalPrice,
		customerOrder.CustomerForm.Address,
		time.Now().Format(time.RFC3339),
		RECEIVED,
	}

	message_broker.Producer(message_broker.FOOD_ORDER_QUEUE_NAME, foodOrder)
}
