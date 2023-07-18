package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/calogxro/qaservice/domain"
	"github.com/calogxro/qaservice/eventstore/service/eventstore"
	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	Key   string `json:"key,omitempty" binding:"required"` // ,omitempty is useful during tests ???????
	Value string `json:"value,omitempty" binding:"required"`
}

type UpdateRequest struct {
	Value string `json:"value" binding:"required"`
}

type Response struct {
	Key  string      `json:"key"`
	Data interface{} `json:"data,omitempty"`
	Err  string      `json:"err,omitempty"` // errors don't JSON-marshal, so we use a stringponse
}

type event struct {
	Type string
	Data domain.Answer
}

type Handler struct {
	service *eventstore.Service
	//projection *projection.QAProjection
}

func MakeHandler(svc *eventstore.Service, r *gin.Engine) *gin.Engine {
	handler := &Handler{
		service: svc,
	}

	r.POST("/answers", handler.CreateAnswer)
	r.PATCH("/answers/:key", handler.UpdateAnswer)
	r.DELETE("/answers/:key", handler.DeleteAnswer)
	r.GET("/answers/:key/history", handler.GetHistory)

	r.GET("/eventstore/ping", ping)

	return r
}

// func New(s *eventstore.Service) *Handler {
// 	return &Handler{
// 		service: s,
// 		//projection: p,
// 	}
// }

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *Handler) CreateAnswer(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer := domain.Answer(req)

	_, err := h.service.CreateAnswer(answer)
	if err != nil {
		if errors.Is(err, domain.ErrKeyExists) {
			//c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			c.JSON(http.StatusConflict, Response{
				Key: answer.Key,
				Err: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//c.JSON(http.StatusCreated, gin.H{"ok": domain.ANSWER_CREATED_EVENT})
	c.JSON(http.StatusCreated, Response{
		Key: answer.Key,
	})
}

func (h *Handler) UpdateAnswer(c *gin.Context) {
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer := domain.Answer{Key: c.Param("key"), Value: req.Value}

	_, err := h.service.UpdateAnswer(answer)
	if err != nil {
		if errors.Is(err, domain.ErrKeyNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//c.JSON(http.StatusOK, gin.H{"ok": domain.ANSWER_UPDATED_EVENT})
	c.JSON(http.StatusOK, Response{
		Key: answer.Key,
	})
}

func (h *Handler) DeleteAnswer(c *gin.Context) {
	key := c.Param("key")

	_, err := h.service.DeleteAnswer(key)
	if err != nil {
		if errors.Is(err, domain.ErrKeyNotFound) {
			//c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			c.JSON(http.StatusNotFound, Response{
				Key: key,
				Err: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//c.JSON(http.StatusOK, gin.H{"ok": domain.ANSWER_DELETED_EVENT})
	c.JSON(http.StatusOK, Response{
		Key: key,
	})
}

func (h *Handler) GetHistory(c *gin.Context) {
	key := c.Param("key")

	events, err := h.service.GetHistory(key)
	if err != nil {
		if errors.Is(err, domain.ErrKeyNotFound) {
			//c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			c.JSON(http.StatusNotFound, Response{
				Key: key,
				Err: err.Error(),
			})
			return
		}
		// db connection error
		c.JSON(http.StatusInternalServerError, nil) //, gin.H{"error": err.Error()})
		return
	}

	var history []event

	for _, ev := range events {
		var answer domain.Answer
		//err := json.Unmarshal([]byte(event.Data), &answer)
		err := json.Unmarshal(ev.Data, &answer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
			return
		}
		history = append(history, event{
			Type: ev.Type,
			Data: answer,
		})
	}

	c.JSON(http.StatusOK, Response{
		Key:  key,
		Data: history,
	})
}
