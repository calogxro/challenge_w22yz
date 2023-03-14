package main

import (
	"github.com/calogxro/qaservice/domain"
	esdbgateway "github.com/calogxro/qaservice/projection/gateway/esdb"
	httphandler "github.com/calogxro/qaservice/projection/handler/http"
	"github.com/calogxro/qaservice/projection/repository/mongodb"
	"github.com/calogxro/qaservice/projection/service/projection"
	"github.com/calogxro/qaservice/projection/service/projector"
	"github.com/gin-gonic/gin"
)

func main() {
	readRepository := mongodb.New()
	p := projector.New(readRepository)

	esgw := esdbgateway.New()
	go esgw.Subscribe(func(event *domain.Event) {
		p.Project(event)
	})

	service := projection.New(readRepository)
	router := httphandler.MakeHandler(service, gin.New())
	router.Run(":8082")
	//router.Run(config.IP_PORT)
}
