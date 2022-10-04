package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Client
var ctx = context.TODO()

const DB_URI = "mongodb://root:example@%s:27017/"
const DB_HOST = "localhost"
const DB_NAME = "answers"
const DB_COLL = "answers"
const DB_COLL_TEST = "test_" + DB_COLL

const PORT = ":8080"

const MSG_KEY_EXISTS = "Key exists"
const MSG_KEY_NOTFOUND = "Key not found"

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func init() {
	dbHost := getEnv("DB_HOST", DB_HOST)
	db = connect(fmt.Sprintf(DB_URI, dbHost))
}

func NewRouter(ctrl *Controller) *gin.Engine {
	r := gin.New() //gin.Default()
	r.GET("/ping", ping)
	r.POST("/answers", ctrl.createAnswer)
	r.GET("/answers/:key", ctrl.findAnswer)
	r.PATCH("/answers/:key", ctrl.updateAnswer)
	r.DELETE("/answers/:key", ctrl.deleteAnswer)
	r.GET("/answers/:key/history", ctrl.getHistory)
	return r
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func main() {
	coll := db.Database(DB_NAME).Collection(DB_COLL)
	eventStore := &EventStore{coll}
	ctrl := Controller{eventStore}
	router := NewRouter(&ctrl)
	router.Run(PORT)
}
