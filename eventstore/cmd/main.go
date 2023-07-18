package main

import (
	"github.com/calogxro/qaservice/eventstore/factory"
	"github.com/calogxro/qaservice/eventstore/repository/esdb"
	"github.com/calogxro/qaservice/eventstore/service/eventstore"
	"github.com/gin-gonic/gin"
)

func main() {
	var repo eventstore.EventStore = esdb.New()
	router := factory.MakeServer(repo, gin.New())
	router.Run(":8080")
	//router.Run(config.IP_PORT)
}
