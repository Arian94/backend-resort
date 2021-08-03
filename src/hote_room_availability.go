package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type HotelReservation struct {
	RoomSubtype   string `json:"roomSubtype"`   // 3 possibilites
	NumberOfRooms byte   `json:"numberOfRooms"` // 1 to 3
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
}

type ClientRequest struct {
	RoomType string `json:"roomType"` // table name
	HotelReservation
}

type AllReservedData struct {
	HotelReservation
	ID         uint16 `json:"id"`
	Fullname   string `json:"fullname"`
	Email      string `json:"email"`
	RoomNumber uint8  `json:"roomNumber"`
}

// Table names are declared below as constants to prevent SQL injection as they will be concatenated to the query.
const (
	SINGLE_ROOM = "single_room"
	DOUBLE_ROOM = "double_room"
	TRIPLE_ROOM = "triple_room"
	TWIN_ROOM   = "twin_room"
)

const (
	SUBTYPE_STANDARD      = "standard"
	SUBTYPE_STANDARD_PLUS = "standard_plus"
	SUBTYPE_DELUXE        = "deluxe"
)

// Rooms Capacity
var roomTypeCapacity = map[string]byte{
	"MAX_SINGLE_ROOM_STANDARD":      3,
	"MAX_SINGLE_ROOM_STANDARD_PLUS": 5,
	"MAX_SINGLE_ROOM_DELUXE":        3,
	// SUM 11

	"MAX_DOUBLE_ROOM_STANDARD":      3,
	"MAX_DOUBLE_ROOM_STANDARD_PLUS": 4,
	"MAX_DOUBLE_ROOM_DELUXE":        3,
	//" SUM 10

	"MAX_TRIPLE_ROOM_STANDARD":      4,
	"MAX_TRIPLE_ROOM_STANDARD_PLUS": 3,
	"MAX_TRIPLE_ROOM_DELUXE":        3,
	// SUM 10

	"MAX_TWIN_ROOM_STANDARD":      3,
	"MAX_TWIN_ROOM_STANDARD_PLUS": 4,
	"MAX_TWIN_ROOM_DELUXE":        4,
	// SUM 11
}

func checkRoomAvailability(c *gin.Context) {
	var clientRequest ClientRequest

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&clientRequest); err != nil {
		log.Printf("Cannot be bound: %v", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	startDate, startDateError := time.Parse("2006-01-02", clientRequest.StartDate)
	endDate, endDateError := time.Parse("2006-01-02", clientRequest.EndDate)
	if startDateError != nil || endDateError != nil {
		log.Printf("Wrong time format: %v or %v", startDate, endDate)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	now := time.Now().UTC()
	oneMonthLater := now.AddDate(0, 1, 0)
	if startDate.Before(now) || startDate.After(oneMonthLater) || endDate.Before(now) || endDate.After(oneMonthLater) ||
		endDate.Before(startDate) {
		log.Println("Wrong time bound")
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	log.Printf("numbof rooms %v", clientRequest.NumberOfRooms)
	if clientRequest.NumberOfRooms < 1 || clientRequest.NumberOfRooms > 3 {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	requestedTableName := strings.ToLower(strings.Split(clientRequest.RoomType, " ")[0]) + "_room" // example: "Single Room" will become "single_room"
	requestedSubtype := strings.ToLower(clientRequest.RoomSubtype)

	var rows *sql.Rows
	var err error

	x0 := requestedTableName == SINGLE_ROOM || requestedTableName == DOUBLE_ROOM || requestedTableName == TRIPLE_ROOM || requestedTableName == TWIN_ROOM
	x1 := requestedSubtype == SUBTYPE_STANDARD || requestedSubtype == SUBTYPE_STANDARD_PLUS || requestedSubtype == SUBTYPE_DELUXE
	if x0 && x1 {
		query := "SELECT * FROM " + requestedTableName
		rows, err = db.Query(query+" WHERE ((start_date BETWEEN ? AND ?) OR (end_date BETWEEN ? AND ?)) AND room_subtype = ?",
			clientRequest.StartDate, clientRequest.EndDate, clientRequest.StartDate, clientRequest.EndDate, clientRequest.RoomSubtype)
		if err != nil {
			log.Printf("Database error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"response:": "Invalid input."})
			return
		}
	} else {
		log.Printf("Wrong table name or room subtype: %v, %v", requestedTableName, clientRequest.RoomSubtype)
		c.JSON(http.StatusBadRequest, gin.H{"response:": "Invalid input."})
		return
	}

	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	var reservedRows []AllReservedData
	for rows.Next() {
		var reservedRow AllReservedData
		if err := rows.Scan(
			&reservedRow.ID,
			&reservedRow.RoomNumber,
			&reservedRow.NumberOfRooms,
			&reservedRow.RoomSubtype,
			&reservedRow.Fullname,
			&reservedRow.Email,
			&reservedRow.StartDate,
			&reservedRow.EndDate,
		); err != nil {
			log.Printf("Row Scanning error: %v", err)
			c.JSON(http.StatusBadGateway, nil)
			return
		}
		reservedRows = append(reservedRows, reservedRow)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusConflict, nil)
		return
	}

	if reservedRows == nil {
		log.Printf("No occupied room, name: %v, cap: %v", clientRequest.RoomType, clientRequest.RoomSubtype)
		c.JSON(http.StatusOK, gin.H{"response": true}) // true means there is/are available room(s)
		return
	} else {
		log.Printf("All occupied rooms, rows: %v", reservedRows)
		isFound := findRoomIfAvailable(&reservedRows, &clientRequest, &requestedTableName)
		c.JSON(http.StatusOK, gin.H{"response": isFound}) // false means there is/are NOT available room(s)
	}
}

func findRoomIfAvailable(reservedRows *[]AllReservedData, clientRequest *ClientRequest, requestedTableName *string) bool {
	var numberOfOccupiedRooms byte = 0

	roomSubtypeKey := "MAX_" + strings.ToUpper(*requestedTableName) + "_" + strings.ToUpper(clientRequest.RoomSubtype)
	maximumRoomCapacity := roomTypeCapacity[roomSubtypeKey]
	log.Printf("%v capacity: %v", roomSubtypeKey, maximumRoomCapacity)
	log.Println("Requested number of rooms:", clientRequest.NumberOfRooms)

	for i := 0; i < len(*reservedRows); i++ {
		numberOfOccupiedRooms = (*reservedRows)[i].NumberOfRooms + numberOfOccupiedRooms
	}
	log.Println("Database sum of occupied rooms:", numberOfOccupiedRooms)

	if maximumRoomCapacity-numberOfOccupiedRooms >= clientRequest.NumberOfRooms {
		return true
	} else {
		return false
	}
}
