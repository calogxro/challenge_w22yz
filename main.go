package main

import (
	"github.com/calogxro/qaservice/config"
	"github.com/calogxro/qaservice/controller"
	es "github.com/calogxro/qaservice/db/event_store"
	rr "github.com/calogxro/qaservice/db/read_repository"
	"github.com/calogxro/qaservice/service"
)

func main() {
	qaservice := service.NewQAService(es.NewEventStoreDB())
	projection := service.NewQAProjection(rr.NewMySQLRepository())
	ctrl := controller.NewController(qaservice, projection)
	router := controller.NewRouter(ctrl)
	router.Run(config.IP_PORT)
}
