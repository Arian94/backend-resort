package signupLogin

import (
	runDatabases "Resort/src/database"
	"Resort/src/models"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/go-playground/validator.v9"
)

func Login(c *gin.Context) {
	db := runDatabases.MongoDb
	ctx := *runDatabases.MongoCtxPtr
	var loginRequest *models.LoginInfo

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

	if err := collection.FindOne(
		ctx,
		bson.M{"profile.email": loginRequest.Email, "profile.password": loginRequest.Password}); err.Err() == mongo.ErrNoDocuments { // means nothing found
		log.Printf("Email not found or wrong password: %v", err)
		c.JSON(http.StatusBadRequest, nil)
	} else {
		log.Printf("User info found")
		generatedToken := JWTAuthService().GenerateToken(loginRequest.Email)
		c.JSON(http.StatusOK, generatedToken)
	}
}

func LoginCourier(c *gin.Context) {
	db := runDatabases.PostgresDb

	var loginRequest *models.LoginInfo

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

	if rows := db.QueryRow("SELECT email FROM courier WHERE email = $1 AND password = $2", loginRequest.Email, loginRequest.Password); rows.Scan().Error() == sql.ErrNoRows.Error() {
		log.Println("No courier info found")
		c.Status(http.StatusBadRequest)
	} else {
		log.Println("Courier info found")
		generatedToken := JWTAuthService().GenerateToken(loginRequest.Email)
		c.JSON(http.StatusOK, generatedToken)
	}
}
