package hotel

import (
	runDatabases "Resort/src/database"
	"Resort/src/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/go-playground/validator.v9"
)

func UpdateRoomMark(c *gin.Context) {
	var updatedBooking *models.UpdatedRoomMark

	if err := c.BindJSON(&updatedBooking); err != nil {
		log.Printf("CheckAndReserveRooms Cannot be bound: %v", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(*updatedBooking); err != nil {
		log.Printf("Validataion error: %v", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	mysqlDb := runDatabases.MysqlDb
	requestedTableName := updatedBooking.RoomType // example: "single_room"
	requestedGenericSubtype := updatedBooking.GenericSubtype

	x0 := requestedTableName == SINGLE_ROOM || requestedTableName == DOUBLE_ROOM || requestedTableName == TRIPLE_ROOM || requestedTableName == TWIN_ROOM
	x1 := requestedGenericSubtype == SUBTYPE_STANDARD || requestedGenericSubtype == SUBTYPE_STANDARD_PLUS || requestedGenericSubtype == SUBTYPE_DELUXE
	if x0 && x1 {
		query := "UPDATE " + requestedTableName
		if _, err := mysqlDb.Query(query+" SET room_mark = ? WHERE id = ?", updatedBooking.RoomMark.String, updatedBooking.Id); err != nil {
			log.Printf("Database error: %v", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
	} else {
		log.Printf("Wrong table name or room subtype: %v, %v", requestedTableName, updatedBooking.GenericSubtype)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	mongoDb := runDatabases.MongoDb
	ctx := *runDatabases.MongoCtxPtr
	collection := mongoDb.Database("resort").Collection("users")

	if _, err := collection.UpdateOne(
		ctx,
		bson.M{
			"profile.email": updatedBooking.Email,
			"hotel": bson.M{
				"$elemMatch": bson.M{
					"roomType":      updatedBooking.RoomType,
					"roomSubtype":   updatedBooking.GenericSubtype,
					"numberOfRooms": updatedBooking.NumberOfRooms,
					"startDate":     updatedBooking.StartDate,
					"endDate":       updatedBooking.EndDate,
				}}},
		bson.M{
			"$set": bson.D{
				{Key: "hotel.$.roomMark", Value: updatedBooking.RoomMark.String},
			},
		},
	); err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, nil)

}

// UPDATE [LOW_PRIORITY] [IGNORE] table_name
// SET
//     column_name1 = expr1,
//     column_name2 = expr2,
//     ...
// [WHERE
//     condition];
