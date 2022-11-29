package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(ctrl *Controller) *gin.Engine {
	r := gin.New() //gin.Default()
	r.GET("/ping", ping)
	r.POST("/answers", ctrl.CreateAnswer)
	r.GET("/answers/:key", ctrl.FindAnswer)
	r.PATCH("/answers/:key", ctrl.UpdateAnswer)
	r.DELETE("/answers/:key", ctrl.DeleteAnswer)
	r.GET("/answers/:key/history", ctrl.GetHistory)
	return r
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
