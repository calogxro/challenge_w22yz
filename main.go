package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namsral/flag"
)

var addr string

func init() {
	flag.StringVar(&addr, "IP_PORT", ":8080", "ip:port to expose")
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
	service := NewQAService(NewEventStoreDB())
	projection := NewQAProjection(NewMySQLReadRepository())
	ctrl := NewController(service, projection)
	router := NewRouter(ctrl)
	router.Run(addr)
}
