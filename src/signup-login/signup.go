package signupLogin

import (
	runMysqlDatabase "Resort/src/mysql-database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phoneNumber"`
}

func Signup(c *gin.Context) {
	db := runMysqlDatabase.Db
	var userInfo UserInfo

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&userInfo); err != nil {
		log.Printf("Cannot be bound: %v", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	rows, err := db.Query("SELECT email FROM users WHERE email = ?", userInfo.Email)
	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"response:": "Invalid input."})
		return
	}

	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	var userDatabaseInfo UserInfo
	for rows.Next() {
		// var reservedRow AllReservedData
		if err := rows.Scan(
			// &userDatabaseInfo.Firstname,
			// &userDatabaseInfo.Lastname,
			&userDatabaseInfo.Email,
			// &userDatabaseInfo.Password,
			// &userDatabaseInfo.PhoneNumber,
		); err != nil {
			log.Printf("Row Scanning error: %v", err)
			c.JSON(http.StatusBadGateway, nil)
			return
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error on closing rows: %v", err)
		c.JSON(http.StatusConflict, nil)
		return
	}

	if userDatabaseInfo.Email != "" {
		log.Printf("Email already registered")
		c.JSON(http.StatusBadRequest, gin.H{"response:": "Email already registered"})
		return
	}

	if addUserError := addUserToDatabase(&userInfo); addUserError != nil {
		log.Printf("Insert to database error: %v", addUserError)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	generatedToken := JWTAuthService().GenerateToken(userDatabaseInfo.Email)

	c.JSON(http.StatusOK, generatedToken)
}

func addUserToDatabase(user *UserInfo) error {
	db := runMysqlDatabase.Db

	log.Printf(user.Email, user.Firstname)
	_, err := db.Query(
		"INSERT INTO users (firstname, lastname, email, password, phone_number) VALUES (?, ?, ?, ?, ?);",
		user.Firstname, user.Lastname, user.Email, user.Password, user.PhoneNumber)
	return err
}
