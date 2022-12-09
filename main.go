package main

import (
	"github.com/calogxro/qaservice/config"
	"github.com/calogxro/qaservice/controller"
	es "github.com/calogxro/qaservice/db/event_store"
	rr "github.com/calogxro/qaservice/db/read_repository"
	"github.com/calogxro/qaservice/domain"
	"github.com/calogxro/qaservice/service"
)

func main() {
	es := es.NewEventStoreDB()
	rr := rr.NewMongoRepository()

	projector := service.NewProjector(rr)
	go es.Subscribe(func(event *domain.Event) {
		projector.Project(event)
	})

	qaservice := service.NewQAService(es)
	projection := service.NewQAProjection(rr)
	ctrl := controller.NewController(qaservice, projection)
	router := controller.NewRouter(ctrl)
	router.Run(config.IP_PORT)
}
