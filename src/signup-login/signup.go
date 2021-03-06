package signupLogin

import (
	runDatabases "Resort/src/database"
	"Resort/src/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/go-playground/validator.v9"
)

func Signup(c *gin.Context) {
	db := runDatabases.MongoDb
	ctx := *runDatabases.MongoCtxPtr

	var userInfo *models.UserSignupInfo

	// Call BindJSON to bind the received JSON/BSON to struct
	if err := c.BindJSON(&userInfo); err != nil {
		log.Printf("Cannot be bound: %v", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(*userInfo); err != nil {
		log.Printf("Validataion error: %v", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	collection := db.Database("resort").Collection("users")
	if singleResult := collection.FindOne(ctx, bson.M{"email": userInfo.Email}); singleResult.Err() != nil { // means nothing found
		if addUserError := addUserToDatabase(userInfo, collection); addUserError != nil {
			log.Printf("Insert to database error: %v", addUserError)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
	} else {
		log.Printf("Email already registered. %v", singleResult.Err())
		c.AbortWithStatusJSON(http.StatusOK, "Email already registered")
		return
	}

	generatedToken := JWTAuthService().GenerateToken(userInfo.Email)
	c.JSON(http.StatusOK, generatedToken)
}

func addUserToDatabase(user *models.UserSignupInfo, collection *mongo.Collection) error {
	ctx := *runDatabases.MongoCtxPtr

	_, err := collection.InsertOne(ctx, bson.M{
		"profile": bson.D{
			{Key: "firstName", Value: user.FirstName},
			{Key: "lastName", Value: user.LastName},
			{Key: "email", Value: user.Email},
			{Key: "password", Value: user.Password},
			{Key: "phoneNumber", Value: user.PhoneNumber},
		}})

	return err
}
