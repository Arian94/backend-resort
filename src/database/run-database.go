package runDatabases

import (
	"context"
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MysqlDb *sql.DB
var MongoDb *mongo.Client
var PostgresDb *sql.DB
var RedisDb redis.Conn
var MongoCtxPtr *context.Context

func RunMySql() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   "arian",
		Passwd: "123",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "resort",
	}

	// Get a database handle.
	var err error
	MysqlDb, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := MysqlDb.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("MySQL Connected!")
}

func RunMongoDB() {
	// Connect to my cluster
	var err error
	MongoDb, err = mongo.NewClient(
		options.Client().ApplyURI("mongodb://mongodb0.localhost:27017"),
	)
	if err != nil {
		log.Fatal(err)
	}
	MongoCtx := context.Background() //.WithTimeout(context.Background(), 1e9*time.Second)
	MongoCtxPtr = &MongoCtx
	// cancel
	err = MongoDb.Connect(MongoCtx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Mongo connected!")
	// defer client.Disconnect(ctx)
}

func RunPostgres() {
	connStr := "user=postgres dbname=resort"
	var err error
	PostgresDb, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error postgres connection: ", err)
	} else {
		log.Println("Postgres Connected!")
	}
}

func RunRedis() {
	var err error
	RedisDb, err = redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal("Error postgres connection: ", err)
	} else {
		log.Println("Redis Connected!")
	}
}
