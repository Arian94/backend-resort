package userInfo

import (
	runDatabases "Resort/src/database"
	"Resort/src/hotel"
	"Resort/src/middleware"
	"Resort/src/restaurant"
	signupLogin "Resort/src/signup-login"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/go-playground/validator.v9"
)

type AllUserInfo struct {
	signupLogin.UserSignupInfo `mapstructure:",squash"`
	Hotel                      []hotel.ClientRequest `json:"hotel" bson:"hotel" mapstructure:"hotel"`
	Restaurant                 []struct {
		restaurant.OrdersAndPrice `mapstructure:",squash"`
		FullName                  string `json:"fullname" bson:"receiver" mapstructure:"receiver" validate:"required"`
		OrderedDate               string `json:"orderedDate" mapstructure:"ordered_date"`
	} `json:"restaurant"`
}

func GetUserInfo(c *gin.Context) {
	token, _ := signupLogin.JWTAuthService().ValidateToken(c.GetHeader("Authorization")[len(middleware.BEARER_SCHEMA)+1:])
	claims := token.Claims.(jwt.MapClaims)
	emailFromToken := fmt.Sprintf("%v", claims["email"])

	time.Sleep(1 * time.Second)
	db := runDatabases.MongoDb
	ctx := *runDatabases.MongoCtxPtr

	var unstructuredUserInfo bson.M

	collection := db.Database("resort").Collection("users")

	if err := collection.FindOne(ctx, bson.M{"email": emailFromToken}).Decode(&unstructuredUserInfo); err != nil {
		log.Printf("Get user info Mongo Database error: %v", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	var structuredUserInfo AllUserInfo
	mapstructure.Decode(unstructuredUserInfo, &structuredUserInfo)

	c.JSON(http.StatusOK, structuredUserInfo)
}

func UpdateUserInfo(c *gin.Context) {
	token, _ := signupLogin.JWTAuthService().ValidateToken(c.GetHeader("Authorization")[len(middleware.BEARER_SCHEMA)+1:])
	claims := token.Claims.(jwt.MapClaims)
	emailFromToken := fmt.Sprintf("%v", claims["email"])

	var generalUserInfo *signupLogin.GeneralUserInfo

	// Call BindJSON to bind the received JSON/BSON to struct
	if err := c.BindJSON(&generalUserInfo); err != nil {
		log.Printf("Cannot be bound: %v", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(*generalUserInfo); err != nil {
		log.Printf("Validataion error: %v", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	log.Printf("Get user info Mongo Database error: %v", generalUserInfo)
	db := runDatabases.MongoDb
	ctx := *runDatabases.MongoCtxPtr

	collection := db.Database("resort").Collection("users")

	if _, err := collection.UpdateOne(ctx, bson.M{"email": emailFromToken},
		bson.M{"$set": bson.M{
			"firstname":    generalUserInfo.Firstname,
			"lastname":     generalUserInfo.Lastname,
			"phone_number": generalUserInfo.PhoneNumber,
			"profile":      bson.D{{Key: "address", Value: generalUserInfo.Profile.Address}},
		}}); err != nil {
		log.Printf("Get user info Mongo Database error: %v", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}
