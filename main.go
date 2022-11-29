package main

import (
	"github.com/calogxro/qaservice/controller"
	"github.com/calogxro/qaservice/db"
	"github.com/calogxro/qaservice/service"
	"github.com/namsral/flag"
)

var addr string

func init() {
	flag.StringVar(&addr, "IP_PORT", ":8080", "ip:port to expose")
	flag.Parse()
}

func main() {
	qaservice := service.NewQAService(db.NewEventStoreDB())
	projection := service.NewQAProjection(db.NewMySQLReadRepository())
	ctrl := controller.NewController(qaservice, projection)
	router := controller.NewRouter(ctrl)
	router.Run(addr)
}
