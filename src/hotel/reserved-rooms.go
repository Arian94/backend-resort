package hotel

import (
	runDatabases "Resort/src/database"
	"Resort/src/message_broker"
	"Resort/src/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// var isWebSocketOpen bool = false

func ReservedRooms(c *gin.Context) {
	mysqlDb := runDatabases.MysqlDb
	tables := []string{SINGLE_ROOM, DOUBLE_ROOM, TRIPLE_ROOM, TWIN_ROOM}

	var rows *sql.Rows
	var err error
	bookedRows := map[string][]models.Bookings{
		SINGLE_ROOM: {},
		DOUBLE_ROOM: {},
		TRIPLE_ROOM: {},
		TWIN_ROOM:   {},
	}

	for _, table := range tables {
		if rows, err = mysqlDb.Query("SELECT id, full_name, room_mark, email, room_subtype, number_of_rooms, start_date, end_date FROM " + table); err != nil {
			log.Printf("Mysql Database error for getting bookings: %v", err)
			c.JSON(http.StatusInternalServerError, nil)
		}

		columnNames, _ := rows.Columns()
		columnLength := len(columnNames)
		row := models.Bookings{}

		s := reflect.ValueOf(&row).Elem()
		columns := make([]interface{}, columnLength)

		roomMarkStruct := reflect.ValueOf(s.Field(0).Addr().Interface()).Elem()
		for i := 0; i < columnLength-4; i++ {
			field := roomMarkStruct.Field(i)
			columns[i] = field.Addr().Interface() // 0, 1, 2, 3
		}

		hotelResStruct := reflect.ValueOf(s.Field(1).Addr().Interface()).Elem()
		columns[6] = hotelResStruct.Field(1).Addr().Interface()
		columns[7] = hotelResStruct.Field(2).Addr().Interface()
		numberAndGenericSubtypeStruct := reflect.ValueOf(hotelResStruct.Field(0).Addr().Interface()).Elem()
		columns[4] = numberAndGenericSubtypeStruct.Field(0).Addr().Interface()
		columns[5] = numberAndGenericSubtypeStruct.Field(1).Addr().Interface()

		for rows.Next() {
			if err := rows.Scan(columns...); err != nil {
				log.Printf("Row Scanning error: %v", err)
				c.JSON(http.StatusInternalServerError, nil)
				return
			}
			bookedRows[table] = append(bookedRows[table], row)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Error on closing rows: %v", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
	}

	c.JSON(http.StatusOK, bookedRows)
}

//webSocket returns json format
func BookingsWebSocket(c *gin.Context) {
	//Upgrade get request to webSocket protocol
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("error get connection")
		log.Fatal(err)
	}

	authHeader := c.GetHeader("Sec-Websocket-Protocol")
	token, _ := jwt.Parse(authHeader, nil)
	claims := token.Claims.(jwt.MapClaims)
	emailFromToken := fmt.Sprintf("%v", claims["email"])
	log.Println("Hotel Admin Email:", emailFromToken)

	BookingQueue, err := message_broker.MqChannel.QueueDeclare(
		message_broker.BOOKING_QUEUE_NAME, // name
		false,                             // durable
		false,                             // delete when unused
		false,                             // exclusive
		false,                             // no-wait
		nil,                               // arguments
	)
	message_broker.FailOnError(err, "Failed to declare the bookingQueue queue")

	message_broker.MqChannel.Cancel(emailFromToken, false) // for dev purposes

	msgs, err := message_broker.MqChannel.Consume(
		BookingQueue.Name, // queue
		emailFromToken,    // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	message_broker.FailOnError(err, "Failed to register a consumer")

	// forever := make(chan bool)

	go func() {
		for d := range msgs {
			err = ws.WriteJSON(string(d.Body))
			if err != nil {
				log.Println("error write json")
			}
		}
	}()

	go func() {
		code, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("error read message", err)
		}
		log.Printf("Received a message: %s", message)

		if code == -1 {
			message_broker.MqChannel.Cancel(emailFromToken, false) // client page is closed, no need to keep filling the queue since on returning back to the page, updated data will be requested by http protocol.
		}
	}()

	// log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// <-forever
}
