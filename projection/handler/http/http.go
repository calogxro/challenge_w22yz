package http

import (
	"errors"
	"net/http"

	"github.com/calogxro/qaservice/domain"
	"github.com/calogxro/qaservice/projection/service/projection"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	//service    *service.QAService
	projection *projection.Projection
}

func MakeHandler(svc *projection.Projection, r *gin.Engine) *gin.Engine {
	handler := &Handler{
		projection: svc,
	}

	r.GET("/answers/:key", handler.FindAnswer)

	r.GET("/projection/ping", ping)

	return r
}

// func New(p *projection.Projection) *Handler {
// 	return &Handler{
// 		//service:    s,
// 		projection: p,
// 	}
// }

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// FindAnswer returns the latest answer for a given key.
func (h *Handler) FindAnswer(c *gin.Context) {
	key := c.Param("key")

	answer, err := h.projection.GetAnswer(key)
	if err != nil {
		if errors.Is(err, domain.ErrKeyNotFound) {
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
