package main

import (
	"time"

	runDatabases "Resort/src/database"
	hotel "Resort/src/hotel"
	"Resort/src/message_broker"
	"Resort/src/middleware"
	"Resort/src/restaurant"
	signupLogin "Resort/src/signup-login"
	userInfo "Resort/src/user-info"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func main() {
	runDatabases.RunMySql()
	runDatabases.RunMongoDB()
	runDatabases.RunPostgres()
	runDatabases.RunRedis()
	restaurant.GetCourierIndex()
	message_broker.InitializeRabbitMq()

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
	router.POST("/loginCourier", signupLogin.LoginCourier)
	router.GET("/userInfo", middleware.AuthorizeJWT, userInfo.GetUserInfo)
	router.GET("/foodList", restaurant.GetFoodPrices)
	router.GET("/bookedRooms", hotel.ReservedRooms)
	router.GET("/customerFoodOrder", restaurant.CustomerFoodOrder)
	router.POST("/orderFoods", middleware.AuthorizeOptionalJWT(), restaurant.OrderFoods)
	router.POST("/hotelRooms", hotel.CheckRoomAvailability)
	router.PATCH("/reserveRooms", middleware.AuthorizeJWT, hotel.CheckAndReserveRooms)
	router.PATCH("/updateUserInfo", middleware.AuthorizeJWT, userInfo.UpdateUserInfo)
	router.PATCH("/updateRoomMark", hotel.UpdateRoomMark)
	router.PATCH("/updateFoodOrderState", restaurant.UpdateFoodOrderState)

	router.GET("/bookingws", middleware.AuthorizeOptionalJWT(), hotel.BookingsWebSocket)
	router.GET("/foodws", middleware.AuthorizeOptionalJWT(), restaurant.FoodOrderWebSocket)
	router.GET("/courierws", middleware.AuthorizeOptionalJWT(), restaurant.CourierWebSocket)

	router.Run("localhost:8080")
}
