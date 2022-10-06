package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namsral/flag"
)

var addr, dbUser, dbPass, dbHost, dbPort string

const dbName = "answers"
const dbColl = "answers"

const MSG_KEY_EXISTS = "Key exists"
const MSG_KEY_NOTFOUND = "Key not found"

func init() {
	flag.StringVar(&addr, "IP_PORT", ":8080", "ip:port to expose")
	flag.StringVar(&dbUser, "DB_USER", "", "db user")
	flag.StringVar(&dbPass, "DB_PASS", "", "db password")
	flag.StringVar(&dbHost, "DB_HOST", "0.0.0.0", "db host")
	flag.StringVar(&dbPort, "DB_PORT", "", "db port")

	flag.String(flag.DefaultConfigFlagname, ".env", "path to config file")
	flag.Parse()
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
	eventStore := &MongoStore{
		dbUser: dbUser,
		dbPass: dbPass,
		dbHost: dbHost,
		dbPort: dbPort,
		dbName: dbName,
		dbColl: dbColl,
	}
	ctrl := NewController(eventStore)
	router := NewRouter(ctrl)
	router.Run(addr)
}
