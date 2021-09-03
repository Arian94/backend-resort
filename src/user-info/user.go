package userInfo

import (
	runDatabases "Resort/src/database"
	"Resort/src/hotel"
	"Resort/src/middleware"
	signupLogin "Resort/src/signup-login"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
)

type AllUserInfo struct {
	signupLogin.UserSignupInfo `mapstructure:",squash"`
	Hotel                      []hotel.ClientRequest `json:"hotel" bson:"hotel" mapstructure:"hotel"`
}

func GetUserInfo(c *gin.Context) {
	token, _ := signupLogin.JWTAuthService().ValidateToken(c.GetHeader("Authorization")[len(middleware.BEARER_SCHEMA)+1:])
	claims := token.Claims.(jwt.MapClaims)
	emailFromToken := fmt.Sprintf("%v", claims["email"])

	time.Sleep(3 * time.Second)
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
	// log.Println("beguo", gooz.Firstname, err)
	// log.Println("beguo", userInfo)

	// var convertedDt AllUserInfo
	// dbByte, _ := bson.Marshal(userInfo)
	// if err := bson.Unmarshal(dbByte, &convertedDt); err != nil {
	// 	log.Println("beguo", err)
	// }

	// var convertedDt1 AllUserInfo
	// dbByte1, _ := json.Marshal(userInfo)
	// _ = json.Unmarshal(dbByte1, &convertedDt1)

	// log.Println("email ro nega", middleware.EmailFromToken)
	// log.Println("az mongo", userInfo)
	// log.Println("az gooz", convertedDt)

	c.JSON(http.StatusOK, structuredUserInfo)
}
