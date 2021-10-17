package hotel

import (
	runDatabases "Resort/src/database"
	"Resort/src/message_broker"
	"Resort/src/middleware"
	"Resort/src/models"
	signupLogin "Resort/src/signup-login"
	"encoding/json"

	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CheckAndReserveRooms(c *gin.Context) {
	token, _ := signupLogin.JWTAuthService().ValidateToken(c.GetHeader("Authorization")[len(middleware.BEARER_SCHEMA)+1:])
	claims := token.Claims.(jwt.MapClaims)
	emailFromToken := fmt.Sprintf("%v", claims["email"])
	mysqlDb := runDatabases.MysqlDb

	var clientRequests *[]models.ClientRequest
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
		if response, err := checkRoom(&clientRequest, mysqlDb); err != nil {
			log.Println("Check room:", err)
			c.JSON(response.status, gin.H{"response": response.result})
			return
		} else if !response.result {
			c.JSON(response.status, gin.H{"response": response.result})
			return
		}
	}

	newBookings := map[string][]models.Bookings{}

	for _, clientRequest := range *clientRequests {
		if response, idFullName, err := reserveRoom(&clientRequest, mysqlDb, emailFromToken); err != nil {
			log.Println("Reserve room:", err)
			c.JSON(response.status, gin.H{"response": response.result})
			return
		} else {
			newBookings[clientRequest.RoomType] = append(newBookings[clientRequest.RoomType], models.Bookings{
				RoomMarkStruct: models.RoomMarkStruct{
					Id:       int16(idFullName.id),
					FullName: idFullName.fullName,
					RoomMark: sql.NullString{},
					Email:    emailFromToken,
				},
				HotelReservation: models.HotelReservation{
					NumberAndGenericSubtype: models.NumberAndGenericSubtype{
						GenericSubtype: clientRequest.GenericSubtype,
						NumberOfRooms:  clientRequest.NumberOfRooms,
					},
					StartDate: clientRequest.StartDate,
					EndDate:   clientRequest.EndDate,
				},
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"response": true})

	if isWebSocketOpen {
		log.Println("New booking is sent to queue")
		msg, _ := json.Marshal(newBookings)
		message_broker.BookingProducer(msg)
	}
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

func duplicateInArray(arrayRequest *[]models.ClientRequest) error {
	var roomTypeArray [3]string
	length := len(*arrayRequest)

	for i := 0; i < length; i++ {
		roomTypeArray[i] = (*arrayRequest)[i].RoomType
	}

	for i := 0; i < length; i++ {
		var foundRoomTypeIndex int
		if i+1 < length {
			foundRoomTypeIndex = indexOf(roomTypeArray[i], roomTypeArray[i+1:])
			log.Println("PEIDA:", foundRoomTypeIndex)

			if (foundRoomTypeIndex != -1) && ((*arrayRequest)[i].GenericSubtype == (*arrayRequest)[foundRoomTypeIndex].GenericSubtype) {
				log.Println("akar:", (*arrayRequest)[i].GenericSubtype, (*arrayRequest)[foundRoomTypeIndex].GenericSubtype)
				x := ((*arrayRequest)[i].StartDate == (*arrayRequest)[foundRoomTypeIndex].StartDate) || ((*arrayRequest)[i].EndDate == (*arrayRequest)[foundRoomTypeIndex].EndDate) ||
					((*arrayRequest)[i].StartDate == (*arrayRequest)[foundRoomTypeIndex].EndDate) || ((*arrayRequest)[i].EndDate == (*arrayRequest)[foundRoomTypeIndex].StartDate)
				if x {
					return errors.New("duplicate element")
				}
			}
		}
	}

	return nil
}

func reserveRoom(clientRequest *models.ClientRequest, mysqlDb *sql.DB, email string) (HttpResponse, struct {
	id       int64
	fullName string
}, error) {
	fullname, err := findUser(email)
	var id int64
	idFullName := struct {
		id       int64
		fullName string
	}{}

	if err != nil {
		log.Printf("Mongo Database error for reserving room: %v", err)
		return HttpResponse{status: http.StatusInternalServerError}, idFullName, err
	}

	query := "INSERT INTO " + clientRequest.RoomType
	if result, err := mysqlDb.Exec(query+" (number_of_rooms, room_subtype, full_name, email, start_date, end_date) VALUES (?, ?, ?, ?, ?, ?)",
		clientRequest.NumberOfRooms, clientRequest.GenericSubtype, fullname, email, clientRequest.StartDate, clientRequest.EndDate); err != nil {
		log.Printf("Mysql Database error for reserving room: %v", err)
		return HttpResponse{status: http.StatusInternalServerError}, idFullName, err
	} else {
		id, _ = result.LastInsertId()
	}

	if err := insertRoomsToUser(email, clientRequest); err != nil {
		log.Printf("Mongo Database insertRoomsToUser error for reserving room: %v", err)
		return HttpResponse{status: http.StatusInternalServerError}, idFullName, err
	}

	idFullName = struct {
		id       int64
		fullName string
	}{id, fullname}

	return HttpResponse{status: http.StatusOK}, idFullName, nil
}

func findUser(email string) (string, error) {
	mongoDb := runDatabases.MongoDb
	ctx := *runDatabases.MongoCtxPtr
	collection := mongoDb.Database("resort").Collection("users")

	var fullnameBson bson.M
	var fullname models.FullName
	if err := collection.FindOne(ctx, bson.M{"profile.email": email}, &options.FindOneOptions{Projection: bson.M{"profile.firstName": 1, "profile.lastName": 1}}).Decode(&fullnameBson); err != nil {
		log.Printf("findUser Mongo Database error: %v", err)
		return "", err
	}
	mapstructure.Decode(fullnameBson["profile"], &fullname)

	return fmt.Sprintf("%v %v", fullname.FirstName, fullname.LastName), nil
}

func insertRoomsToUser(email string, roomArray *models.ClientRequest) error {
	mongoDb := runDatabases.MongoDb
	ctx := *runDatabases.MongoCtxPtr
	collection := mongoDb.Database("resort").Collection("users")

	if _, err := collection.UpdateOne(
		ctx,
		bson.M{"profile.email": email},
		bson.M{"$push": bson.M{
			"hotel": bson.D{
				// {Key: "_id", Value: bson.TypeObjectID},
				{Key: "roomType", Value: roomArray.RoomType},
				{Key: "roomMark", Value: nil},
				{Key: "numberOfRooms", Value: roomArray.NumberOfRooms},
				{Key: "roomSubtype", Value: roomArray.GenericSubtype},
				{Key: "startDate", Value: roomArray.StartDate},
				{Key: "endDate", Value: roomArray.EndDate},
			}}},
	); err != nil {
		return err
	}

	return nil
}
