package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateAnswerRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type UpdateAnswerRequest struct {
	Value string `json:"value" binding:"required"`
}

type Controller struct {
	service    *QAService
	projection *QAProjection
}

func NewController(s *QAService, p *QAProjection) *Controller {
	return &Controller{
		service:    s,
		projection: p,
	}
}

func (ctrl *Controller) createAnswer(c *gin.Context) {
	var req CreateAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer := Answer(req)

	_, err := ctrl.service.CreateAnswer(answer)
	if err != nil {
		if _, keyExists := err.(*KeyExists); keyExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ok": ANSWER_CREATED_EVENT})
}

// return the latest answer for a given key
func (ctrl *Controller) findAnswer(c *gin.Context) {
	key := c.Param("key")

	answer, err := ctrl.projection.GetAnswer(key)
	if err != nil {
		if _, keyNotFound := err.(*KeyNotFound); keyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
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
	var req UpdateAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer := Answer{Key: c.Param("key"), Value: req.Value}

	_, err := ctrl.service.UpdateAnswer(answer)
	if err != nil {
		if _, keyNotFound := err.(*KeyNotFound); keyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": ANSWER_UPDATED_EVENT})
}

func (ctrl *Controller) deleteAnswer(c *gin.Context) {
	key := c.Param("key")

	_, err := ctrl.service.DeleteAnswer(key)
	if err != nil {
		if _, keyNotFound := err.(*KeyNotFound); keyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": ANSWER_DELETED_EVENT})
}

func (ctrl *Controller) getHistory(c *gin.Context) {
	key := c.Param("key")

	history, err := ctrl.service.GetHistory(key)
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
