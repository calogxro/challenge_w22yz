package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/calogxro/qaservice/domain"
	gateway "github.com/calogxro/qaservice/projector/gateway/esdbgw"
	"github.com/calogxro/qaservice/projector/repository/mongodb"
	"github.com/calogxro/qaservice/projector/service/projector"
)

func main() {
	readRepository := mongodb.New()
	p := projector.New(readRepository)

	esgw := gateway.New()
	go esgw.Subscribe(func(event *domain.Event) {
		log.Printf("%s", event)
		p.Project(event)
	})

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press ctrl+c to continue...")
	<-done // Will block here until user hits ctrl+c
}
