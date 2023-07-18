package factory

import (
	httphandler "github.com/calogxro/qaservice/projection/handler/http"
	"github.com/calogxro/qaservice/projection/service/projection"
	"github.com/gin-gonic/gin"
)

func MakeServer(repo projection.ReadRepository, router *gin.Engine) *gin.Engine {
	service := projection.New(repo)
	router = httphandler.MakeHandler(service, router)
	return router
}
