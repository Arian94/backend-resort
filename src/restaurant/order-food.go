package restaurant

import (
	runDatabases "Resort/src/database"
	"Resort/src/middleware"
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

type foodOrder struct {
	Orders []struct {
		Name          string `json:"name" bson:"name" validate:"required"`
		NumberOfMeals byte   `json:"numberOfMeals" bson:"number_of_meals" validate:"required"`
	} `json:"orders"`
	Customer struct {
		FullName    string `json:"fullName" bson:"full_name" validate:"required"`
		Address     string `json:"address" bson:"address" validate:"required"`
		PhoneNumber string `json:"phoneNumber" bson:"phone_number" validate:"required"`
	} `json:"customer"`
	TotalPrice int64 `json:"totalPrice" validate:"required"`
}

func OrderFoods(c *gin.Context) {
	db := runDatabases.MongoDb
	ctx := *runDatabases.MongoCtxPtr
	var customerOrder *foodOrder

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

	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		token, _ := signupLogin.JWTAuthService().ValidateToken(authHeader[len(middleware.BEARER_SCHEMA)+1:])
		claims := token.Claims.(jwt.MapClaims)
		emailFromToken := fmt.Sprintf("%v", claims["email"])

		userCollection := db.Database("resort").Collection("users")
		if _, err := userCollection.UpdateOne(
			ctx,
			bson.M{"email": emailFromToken},
			bson.M{"$push": bson.M{
				"restaurant": bson.M{
					"orders":       customerOrder.Orders,
					"receiver":     customerOrder.Customer.FullName,
					"total_price":  customerOrder.TotalPrice,
					"ordered_date": time.Now().Format(time.RFC3339),
				},
			},
			},
		); err != nil {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		foodOrderCollection := db.Database("resort").Collection("food_orders")
		if _, err := foodOrderCollection.InsertOne(
			ctx,
			bson.M{"$push": bson.M{
				"email":                 emailFromToken,
				"receiver":              customerOrder.Customer.FullName,
				"receiver_phone_number": customerOrder.Customer.PhoneNumber,
				"orders":                customerOrder.Orders,
				"total_price":           customerOrder.TotalPrice,
				"ordered_date":          time.Now().Format(time.RFC3339),
			},
			},
		); err != nil {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
	} else {
		foodOrderCollection := db.Database("resort").Collection("food_orders")
		if _, err := foodOrderCollection.InsertOne(
			ctx,
			bson.M{"$push": bson.M{
				"receiver":              customerOrder.Customer.FullName,
				"receiver_phone_number": customerOrder.Customer.PhoneNumber,
				"orders":                customerOrder.Orders,
				"total_price":           customerOrder.TotalPrice,
				"ordered_date":          time.Now().Format(time.RFC3339),
			},
			},
		); err != nil {
			c.JSON(http.StatusInternalServerError, nil)
		}
	}

	c.JSON(http.StatusOK, nil)
}
