package restaurant

import (
	runDatabases "Resort/src/database"
	"Resort/src/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
)

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
	usersCollection := db.Database("resort").Collection("users")

	objId, _ := primitive.ObjectIDFromHex(foodOrder.Id)

	if _, err := foodOrderCollection.UpdateByID(ctx, objId, bson.M{"$set": bson.M{"orderState": foodOrder.OrderState}}); err != nil {
		log.Printf("UpdateFoodOrderState UpdateById error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

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
