package main

import (
	"time"

	hotel "Resort/src"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func main() {
	router := gin.Default()

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	router.POST("/hotelRooms", hotel.CheckRoomAvailability)
	// router.GET("/alaki", alaki)
	// router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}
