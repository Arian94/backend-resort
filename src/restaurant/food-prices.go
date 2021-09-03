package restaurant

import (
	runDatabases "Resort/src/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type foodList struct {
	Type        string `json:"type"`
	Origin      string `json:"origin"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Price       uint8  `json:"price"`
}

func GetFoodPrices(c *gin.Context) {
	db := runDatabases.MysqlDb
	var FoodList []foodList

	rows, err := db.Query("SELECT name, price FROM restaurant")
	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var foodRow foodList
		if err := rows.Scan(
			&foodRow.Name,
			&foodRow.Price,
		); err != nil {
			log.Printf("Row Scanning error: %v", err)
			c.JSON(http.StatusBadGateway, nil)
			return
		}
		FoodList = append(FoodList, foodRow)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error on closing rows: %v", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, FoodList)
}
