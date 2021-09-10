package hotel

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	runDatabases "Resort/src/database"
	"Resort/src/models"

	"github.com/gin-gonic/gin"
)

type HttpResponse struct {
	status  int
	result  bool
	message string
}

// Table names are declared below as constants to prevent SQL injection as they will be concatenated to the query.
const (
	SINGLE_ROOM = "single_room"
	DOUBLE_ROOM = "double_room"
	TRIPLE_ROOM = "triple_room"
	TWIN_ROOM   = "twin_room"
)
const (
	SUBTYPE_STANDARD      = "Standard"
	SUBTYPE_STANDARD_PLUS = "Standard_Plus"
	SUBTYPE_DELUXE        = "Deluxe"
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

func CheckRoomAvailability(c *gin.Context) {
	var clientRequest models.ClientRequest

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&clientRequest); err != nil {
		log.Printf("Cannot be bound: %v", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	db := runDatabases.MysqlDb
	response, _ := checkRoom(&clientRequest, db)
	c.JSON(response.status, gin.H{"response": response.result})
}

func checkRoom(clientRequest *models.ClientRequest, db *sql.DB) (HttpResponse, error) {
	format := "2006-01-02"
	startDate, startDateError := time.Parse(format, clientRequest.StartDate)
	endDate, endDateError := time.Parse(format, clientRequest.EndDate)
	if startDateError != nil || endDateError != nil {
		log.Printf("Wrong time format: %v or %v", startDate, endDate)
		return HttpResponse{status: http.StatusBadRequest}, gin.Error{}
	}

	today, _ := time.Parse(format, time.Now().Format(format)) // in order to get rid of time and only use date, time.Now() is stringified and then parsed
	oneMonthLater := today.AddDate(0, 0, 29)

	x := startDate.Before(today) || startDate.After(oneMonthLater) || endDate.Before(today) || endDate.After(oneMonthLater) ||
		endDate.Before(startDate)
	if x {
		log.Printf("Wrong time bound \n today: %v \n one month later: %v \n start: %v \n end: %v",
			today, oneMonthLater, startDate, endDate)
		return HttpResponse{status: http.StatusBadRequest}, gin.Error{}
	}

	log.Printf("Number of rooms: %v", clientRequest.NumberOfRooms)
	if clientRequest.NumberOfRooms < 1 || clientRequest.NumberOfRooms > 3 {
		return HttpResponse{status: http.StatusBadRequest}, gin.Error{}
	}

	requestedTableName := clientRequest.RoomType // example: "single_room"
	requestedGenericSubtype := clientRequest.GenericSubtype

	var rows *sql.Rows

	x0 := requestedTableName == SINGLE_ROOM || requestedTableName == DOUBLE_ROOM || requestedTableName == TRIPLE_ROOM || requestedTableName == TWIN_ROOM
	x1 := requestedGenericSubtype == SUBTYPE_STANDARD || requestedGenericSubtype == SUBTYPE_STANDARD_PLUS || requestedGenericSubtype == SUBTYPE_DELUXE
	if x0 && x1 {
		var err error
		query := "SELECT room_subtype, number_of_rooms FROM " + requestedTableName
		rows, err = db.Query(query+" WHERE ((start_date BETWEEN ? AND ?) OR (end_date BETWEEN ? AND ?)) AND room_subtype = ?",
			clientRequest.StartDate, clientRequest.EndDate, clientRequest.StartDate, clientRequest.EndDate, clientRequest.GenericSubtype)
		if err != nil {
			log.Printf("Database error: %v", err)
			// context.JSON(http.StatusInternalServerError, nil)
			return HttpResponse{status: http.StatusInternalServerError}, err
		}
	} else {
		log.Printf("Wrong table name or room subtype: %v, %v", requestedTableName, clientRequest.GenericSubtype)
		return HttpResponse{status: http.StatusBadRequest, message: "Invalid input."}, gin.Error{}
	}

	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	var reservedRows []models.NumberAndGenericSubtype
	for rows.Next() {
		var reservedRow models.NumberAndGenericSubtype
		if err := rows.Scan(
			&reservedRow.GenericSubtype,
			&reservedRow.NumberOfRooms,
		); err != nil {
			log.Printf("Row Scanning error: %v", err)
			return HttpResponse{status: http.StatusBadGateway}, err
		}
		reservedRows = append(reservedRows, reservedRow)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error on closing rows: %v", err)
		return HttpResponse{status: http.StatusConflict}, err
	}

	if reservedRows == nil {
		log.Printf("No occupied room, name: %v, cap: %v", clientRequest.RoomType, clientRequest.GenericSubtype)
		return HttpResponse{status: http.StatusOK, result: true}, nil // true means there is/are available room(s)
	} else {
		log.Printf("All occupied rooms, rows: %v", reservedRows)
		isFound := findRoomIfAvailable(&reservedRows, clientRequest, &requestedTableName)
		return HttpResponse{status: http.StatusOK, result: isFound}, nil // false means there is/are NOT available room(s)
	}
}

func findRoomIfAvailable(reservedRows *[]models.NumberAndGenericSubtype, clientRequest *models.ClientRequest, requestedTableName *string) bool {
	var numberOfOccupiedRooms byte = 0

	roomSubtypeKey := "MAX_" + strings.ToUpper(*requestedTableName) + "_" + strings.ToUpper(clientRequest.GenericSubtype)
	maximumRoomCapacity := roomTypeCapacity[roomSubtypeKey]
	log.Printf("%v capacity: %v", roomSubtypeKey, maximumRoomCapacity)
	log.Println("Requested number of rooms:", clientRequest.NumberOfRooms)

	for i := 0; i < len(*reservedRows); i++ {
		numberOfOccupiedRooms = (*reservedRows)[i].NumberOfRooms + numberOfOccupiedRooms
	}
	log.Println("Database sum of occupied rooms:", numberOfOccupiedRooms)

	if (maximumRoomCapacity - numberOfOccupiedRooms) >= clientRequest.NumberOfRooms {
		return true
	} else {
		log.Println("Fully booked")
		return false
	}
}
