package restaurant

import (
	runDatabases "Resort/src/database"
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

	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		token, _ := signupLogin.JWTAuthService().ValidateToken(authHeader[len(middleware.BEARER_SCHEMA)+1:])
		claims := token.Claims.(jwt.MapClaims)
		emailFromToken := fmt.Sprintf("%v", claims["email"])

		userCollection := db.Database("resort").Collection("users")
		if _, err := userCollection.UpdateOne(
			ctx,
			bson.M{"profile.email": emailFromToken},
			bson.M{"$push": bson.M{
				"restaurant": bson.M{
					"orders":      customerOrder.Orders,
					"receiver":    customerOrder.CustomerForm.FullName,
					"totalPrice":  customerOrder.TotalPrice,
					"orderedDate": time.Now().Format(time.RFC3339),
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
			bson.M{
				"email":               emailFromToken,
				"receiver":            customerOrder.CustomerForm.FullName,
				"receiverPhoneNumber": customerOrder.CustomerForm.PhoneNumber,
				"orders":              customerOrder.Orders,
				"totalPrice":          customerOrder.TotalPrice,
				"address":             customerOrder.CustomerForm.Address,
				"orderedDate":         time.Now().Format(time.RFC3339),
			},
		); err != nil {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
	} else {
		foodOrderCollection := db.Database("resort").Collection("food_orders")
		if _, err := foodOrderCollection.InsertOne(
			ctx,
			bson.M{
				"receiver":            customerOrder.CustomerForm.FullName,
				"receiverPhoneNumber": customerOrder.CustomerForm.PhoneNumber,
				"orders":              customerOrder.Orders,
				"totalPrice":          customerOrder.TotalPrice,
				"address":             customerOrder.CustomerForm.Address,
				"orderedDate":         time.Now().Format(time.RFC3339),
			},
		); err != nil {
			c.JSON(http.StatusInternalServerError, nil)
		}
	}

	c.JSON(http.StatusOK, nil)
}
