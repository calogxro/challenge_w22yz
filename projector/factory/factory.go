package factory

import (
	httphandler "github.com/calogxro/qaservice/eventstore/handler/http"
	"github.com/calogxro/qaservice/eventstore/service/eventstore"
	"github.com/gin-gonic/gin"
)

func MakeService(repo eventstore.EventStore, router *gin.Engine) *gin.Engine {
	service := eventstore.New(repo)
	router = httphandler.MakeHandler(service, router)
	return router
}
