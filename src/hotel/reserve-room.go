package hotel

import (
	runDatabases "Resort/src/database"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"Resort/src/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Fullname struct {
	Firstname string `bson:"firstname"`
	Lastname  string `bson:"lastname"`
}

func CheckAndReserveRooms(c *gin.Context) {
	mysqlDb := runDatabases.MysqlDb

	var clientRequests *[]ClientRequest
	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&clientRequests); err != nil {
		log.Printf("CheckAndReserveRooms Cannot be bound: %v", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	length := len(*clientRequests)
	if length > 3 || length < 1 {
		log.Printf("Exceeded array length")
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if err := duplicateInArray(clientRequests); err != nil {
		log.Printf("%v", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	for _, clientRequest := range *clientRequests {
		if result, err := checkRoom(&clientRequest, mysqlDb); err != nil {
			log.Println("Check room:", err)
			c.JSON(result.status, gin.H{"response": result.result})
			return
		} else if !result.result {
			c.JSON(result.status, gin.H{"response": result.result})
			return
		}
	}

	for _, clientRequest := range *clientRequests {
		if result, err := reserveRoom(&clientRequest, mysqlDb); err != nil {
			log.Println("Reserve room:", err)
			c.JSON(result.status, gin.H{"response": result.result})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"response": true})
}

func indexOf(word string, data []string) int {
	log.Println("kalame:", word, data)
	for k, v := range data {
		if word == v {
			return k + 1
		}
	}
	return -1
}

func duplicateInArray(arrayRequest *[]ClientRequest) error {
	var roomTypeArray [3]string
	length := len(*arrayRequest)

	for i := 0; i < length; i++ {
		roomTypeArray[i] = (*arrayRequest)[i].RoomType
	}
	log.Println("kalame:", arrayRequest, roomTypeArray)

	for i := 0; i < length; i++ {
		var foundroomTypeIndex int
		if i+1 < length {
			foundroomTypeIndex = indexOf(roomTypeArray[i], roomTypeArray[i+1:])
			log.Println("PEIDA:", foundroomTypeIndex)

			if (foundroomTypeIndex != -1) && ((*arrayRequest)[i].GenericSubtype == (*arrayRequest)[foundroomTypeIndex].GenericSubtype) {
				log.Println("akar:", (*arrayRequest)[i].GenericSubtype, (*arrayRequest)[foundroomTypeIndex].GenericSubtype)
				x := ((*arrayRequest)[i].StartDate == (*arrayRequest)[foundroomTypeIndex].StartDate) || ((*arrayRequest)[i].EndDate == (*arrayRequest)[foundroomTypeIndex].EndDate) ||
					((*arrayRequest)[i].StartDate == (*arrayRequest)[foundroomTypeIndex].EndDate) || ((*arrayRequest)[i].EndDate == (*arrayRequest)[foundroomTypeIndex].StartDate)
				if x {
					return errors.New("duplicate element")
				}
			}
		}
	}

	return nil
}

func reserveRoom(clientRequest *ClientRequest, mysqlDb *sql.DB) (httpResponse, error) {
	fullname, err := findUser(middleware.EmailFromToken)
	if err != nil {
		log.Printf("Mongo Database error for reserving room: %v", err)
		return httpResponse{status: http.StatusInternalServerError}, err
	}

	query := "INSERT INTO " + clientRequest.RoomType
	if _, err := mysqlDb.Query(query+" (number_of_rooms, room_subtype, fullname, email, start_date, end_date) VALUES (?, ?, ?, ?, ?, ?)",
		clientRequest.NumberOfRooms, clientRequest.GenericSubtype, fullname, middleware.EmailFromToken, clientRequest.StartDate, clientRequest.EndDate); err != nil {
		log.Printf("Mysql Database error for reserving room: %v", err)
		return httpResponse{status: http.StatusInternalServerError}, err
	}

	if err := insertRoomsToUser(middleware.EmailFromToken, clientRequest); err != nil {
		log.Printf("Mongo Database insertRoomsToUser error for reserving room: %v", err)
		return httpResponse{status: http.StatusInternalServerError}, err
	}

	return httpResponse{status: http.StatusOK}, nil
}

func findUser(email string) (string, error) {
	mongoDb := runDatabases.MongoDb
	ctx := *runDatabases.MongoCtxPtr
	collection := mongoDb.Database("resort").Collection("users")

	var fullname Fullname
	if err := collection.FindOne(ctx, bson.M{"email": email}, &options.FindOneOptions{Projection: bson.M{"firstname": 1, "lastname": 1}}).Decode(&fullname); err != nil {
		log.Printf("findUser Mongo Database error: %v", err)
		return "", err
	}

	return fmt.Sprintf("%v %v", fullname.Firstname, fullname.Lastname), nil
}

func insertRoomsToUser(email string, roomArray *ClientRequest) error {
	mongoDb := runDatabases.MongoDb
	ctx := *runDatabases.MongoCtxPtr
	collection := mongoDb.Database("resort").Collection("users")

	if _, err := collection.UpdateOne(
		ctx,
		bson.M{"email": email},
		bson.M{"$push": bson.M{
			"hotel": bson.D{
				{Key: "room_number", Value: nil},
				{Key: "room_type", Value: roomArray.RoomType},
				{Key: "number_of_rooms", Value: roomArray.NumberOfRooms},
				{Key: "room_subtype", Value: roomArray.GenericSubtype},
				{Key: "start_date", Value: roomArray.StartDate},
				{Key: "end_date", Value: roomArray.EndDate},
			},
		},
		},
	); err != nil {
		return err
	}

	return nil
}
