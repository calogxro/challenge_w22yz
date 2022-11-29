package controller

import (
	"net/http"

	"github.com/calogxro/qaservice/domain"
	"github.com/calogxro/qaservice/service"
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
	service    *service.QAService
	projection *service.QAProjection
}

func NewController(s *service.QAService, p *service.QAProjection) *Controller {
	return &Controller{
		service:    s,
		projection: p,
	}
}

func (ctrl *Controller) CreateAnswer(c *gin.Context) {
	var req CreateAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer := domain.Answer(req)

	_, err := ctrl.service.CreateAnswer(answer)
	if err != nil {
		if _, keyExists := err.(*domain.KeyExists); keyExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ok": domain.ANSWER_CREATED_EVENT})
}

// return the latest answer for a given key
func (ctrl *Controller) FindAnswer(c *gin.Context) {
	key := c.Param("key")

	answer, err := ctrl.projection.GetAnswer(key)
	if err != nil {
		if _, keyNotFound := err.(*domain.KeyNotFound); keyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if exists := answer != nil; !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": domain.MSG_KEY_NOTFOUND})
		return
	}

	c.JSON(http.StatusOK, answer)
}

func (ctrl *Controller) UpdateAnswer(c *gin.Context) {
	var req UpdateAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer := domain.Answer{Key: c.Param("key"), Value: req.Value}

	_, err := ctrl.service.UpdateAnswer(answer)
	if err != nil {
		if _, keyNotFound := err.(*domain.KeyNotFound); keyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": domain.ANSWER_UPDATED_EVENT})
}

func (ctrl *Controller) DeleteAnswer(c *gin.Context) {
	key := c.Param("key")

	_, err := ctrl.service.DeleteAnswer(key)
	if err != nil {
		if _, keyNotFound := err.(*domain.KeyNotFound); keyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": domain.ANSWER_DELETED_EVENT})
}

func (ctrl *Controller) GetHistory(c *gin.Context) {
	key := c.Param("key")

	history, err := ctrl.service.GetHistory(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(history) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": domain.MSG_KEY_NOTFOUND})
		return
	}

	c.JSON(http.StatusOK, history)
}
