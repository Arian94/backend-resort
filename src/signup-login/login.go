package signupLogin

import (
	runDatabases "Resort/src/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/go-playground/validator.v9"
)

type LoginInfo struct {
	Email    string `validate:"email,required"`
	Password string `validate:"required"`
}

func Login(c *gin.Context) {
	db := runDatabases.MongoDb
	ctx := *runDatabases.MongoCtxPtr
	var loginRequest *LoginInfo

	// Call BindJSON to bind the received JSON/BSON to struct
	if err := c.BindJSON(&loginRequest); err != nil {
		log.Printf("Cannot be bound: %v", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(*loginRequest); err != nil {
		log.Printf("Validataion error: %v", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	log.Printf("Login request:\nemail: %v, password: %v", loginRequest.Email, loginRequest.Password)

	collection := db.Database("resort").Collection("users")
	var databaseLoginInfo *LoginInfo
	if err := collection.FindOne(ctx, bson.M{"email": loginRequest.Email}).Decode(&databaseLoginInfo); err != nil { // means nothing found
		log.Printf("Email not found: %v", err)
		c.JSON(http.StatusBadRequest, nil)
	} else {
		log.Printf("User info found:\nemail: %v, password: %v", databaseLoginInfo.Email, databaseLoginInfo.Password)
		if databaseLoginInfo.Password != loginRequest.Password {
			log.Printf("Password not matched: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		generatedToken := JWTAuthService().GenerateToken(loginRequest.Email)
		c.JSON(http.StatusOK, generatedToken)
	}
}
