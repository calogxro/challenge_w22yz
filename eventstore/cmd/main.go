package main

import (
	httphandler "github.com/calogxro/qaservice/eventstore/handler/http"
	"github.com/calogxro/qaservice/eventstore/repository/esdb"
	"github.com/calogxro/qaservice/eventstore/service/eventstore"
	"github.com/gin-gonic/gin"
)

func main() {
	repo := esdb.New()
	service := eventstore.New(repo)
	router := httphandler.MakeHandler(service, gin.New())
	router.Run(":8081")
	//router.Run(config.IP_PORT)
}
