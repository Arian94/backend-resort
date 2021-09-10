package main

import (
	"time"

	runDatabases "Resort/src/database"
	hotel "Resort/src/hotel"
	"Resort/src/middleware"
	"Resort/src/restaurant"
	signupLogin "Resort/src/signup-login"
	userInfo "Resort/src/user-info"

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
		Methods:         "GET,  POST, PUT, PATCH, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	router.POST("/signup", signupLogin.Signup)
	router.POST("/login", signupLogin.Login)
	router.GET("/userInfo", middleware.AuthorizeJWT, userInfo.GetUserInfo)
	router.GET("/foodList", restaurant.GetFoodPrices)
	router.POST("/orderFoods", middleware.AuthorizeOptionalJWT(), restaurant.OrderFoods)
	router.POST("/hotelRooms", hotel.CheckRoomAvailability)
	router.PATCH("/reserveRooms", middleware.AuthorizeJWT, hotel.CheckAndReserveRooms)
	router.PATCH("/updateUserInfo", middleware.AuthorizeJWT, userInfo.UpdateUserInfo)

	router.Run("localhost:8080")
}
