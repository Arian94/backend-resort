package main

import (
	"time"

	runDatabases "Resort/src/database"
	hotel "Resort/src/hotel"
	"Resort/src/middleware"
	signupLogin "Resort/src/signup-login"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

// var Db *sql.DB

func main() {
	// runDatabase()
	runDatabases.RunMySql()
	runDatabases.RunMongoDB()

	router := gin.Default()

	// router.Use(middleware.AuthorizeJWT())

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	router.POST("/signup", signupLogin.Signup)
	router.POST("/hotelRooms", hotel.CheckRoomAvailability)
	router.POST("/reserveRooms", middleware.AuthorizeJWT(), hotel.CheckAndReserveRooms)

	router.Run("localhost:8080")
}

// func runDatabase() {
// 	// Capture connection properties.
// 	cfg := mysql.Config{
// 		User:   "arian",
// 		Passwd: "123",
// 		Net:    "tcp",
// 		Addr:   "127.0.0.1:3306",
// 		DBName: "resort",
// 	}

// 	// Get a database handle.
// 	var err error
// 	Db, err = sql.Open("mysql", cfg.FormatDSN())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	pingErr := Db.Ping()
// 	if pingErr != nil {
// 		log.Fatal(pingErr)
// 	}
// 	log.Println("Connected!")
// }
