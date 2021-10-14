package restaurant

import (
	runDatabases "Resort/src/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func CustomerFoodOrder(c *gin.Context) {
	db := runDatabases.MongoDb
	ctx := *runDatabases.MongoCtxPtr
	collection := db.Database("resort").Collection("food_orders")

	if result, err := collection.Find(ctx, bson.M{}); err != nil {
		log.Printf("Get CustomerFoodOrder Mongo Database error: %v", err)
		c.JSON(http.StatusInternalServerError, nil)
	} else {
		var decodedResult []bson.M
		result.All(ctx, &decodedResult)

		c.JSON(http.StatusOK, decodedResult)
	}
}
