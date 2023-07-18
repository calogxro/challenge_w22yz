package main

import (
	"fmt"

	"github.com/calogxro/qaservice/config"
	"github.com/calogxro/qaservice/projection/factory"
	"github.com/calogxro/qaservice/projection/service/projection"

	//"github.com/calogxro/qaservice/domain"
	//esdbgateway "github.com/calogxro/qaservice/projection/gateway/esdb"

	"github.com/calogxro/qaservice/projection/repository/mongodb"

	//"github.com/calogxro/qaservice/projection/service/projector"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println(config.MongoDB)

	var repo projection.ReadRepository = mongodb.New()
	router := factory.MakeServer(repo, gin.New())
	router.Run(":8081")
	//router.Run(config.IP_PORT)
}
