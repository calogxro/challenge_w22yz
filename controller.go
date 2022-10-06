package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	eventStore EventStore
}

type CreateAnswerReq struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type UpdateAnswerReq struct {
	Value string `json:"value" binding:"required"`
}

func NewController(eventStore EventStore) *Controller {
	eventStore.init()
	return &Controller{eventStore}
}

func (ctrl *Controller) createAnswer(c *gin.Context) {
	var req CreateAnswerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer := Answer{Key: req.Key, Value: req.Value}

	_, err := ctrl.eventStore.create(&answer)
	if err != nil {
		if _, keyExists := err.(*KeyExists); keyExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ok": CREATED_EVENT})
}

// return the latest answer for the given key
func (ctrl *Controller) findAnswer(c *gin.Context) {
	key := c.Param("key")

	answer, err := ctrl.eventStore.find(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if exists := answer != nil; !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": MSG_KEY_NOTFOUND})
		return
	}

	c.JSON(http.StatusOK, answer)
}

func (ctrl *Controller) updateAnswer(c *gin.Context) {
	var req UpdateAnswerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer := Answer{Key: c.Param("key"), Value: req.Value}

	_, err := ctrl.eventStore.update(&answer)
	if err != nil {
		if _, keyNotFound := err.(*KeyNotFound); keyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": UPDATED_EVENT})
}

func (ctrl *Controller) deleteAnswer(c *gin.Context) {
	key := c.Param("key")

	_, err := ctrl.eventStore.delete(key)
	if err != nil {
		if _, keyNotFound := err.(*KeyNotFound); keyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": DELETED_EVENT})
}

func (ctrl *Controller) getHistory(c *gin.Context) {
	key := c.Param("key")

	history, err := ctrl.eventStore.getHistory(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(history) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": MSG_KEY_NOTFOUND})
		return
	}

	c.JSON(http.StatusOK, history)
}
