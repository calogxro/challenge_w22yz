package main

import (
	"flag"
	"log"

	"github.com/calogxro/qaservice/config"
	"github.com/calogxro/qaservice/domain"
	esFactory "github.com/calogxro/qaservice/eventstore/factory"
	"github.com/calogxro/qaservice/eventstore/repository/esdb"
	esmemdb "github.com/calogxro/qaservice/eventstore/repository/memory"
	prFactory "github.com/calogxro/qaservice/projection/factory"
	rrmemdb "github.com/calogxro/qaservice/projection/repository/memory"
	projectionRepo "github.com/calogxro/qaservice/projection/repository/mongodb"
	gateway "github.com/calogxro/qaservice/projector/gateway/esdbgw"
	projectorRepo "github.com/calogxro/qaservice/projector/repository/mongodb"
	"github.com/calogxro/qaservice/projector/service/projector"
	"github.com/gin-gonic/gin"
)

func main() {
	memdbFlag := flag.Bool("memdb", false, "use memory database")
	flag.Parse()

	if *memdbFlag {
		runInMemory()
	} else {
		runWithDatabase()
	}
}

func runInMemory() {
	router := gin.New()

	// Event store
	es := esmemdb.New()
	router = esFactory.MakeServer(es, router)

	// Projection
	rr := rrmemdb.New()
	router = prFactory.MakeServer(rr, router)

	// Projector
	p := projector.New(rr)
	es.Subscribe(func(event *domain.Event) {
		log.Printf("%s", event)
		p.Project(event)
	})

	router.Run(config.IP_PORT)
}

func runWithDatabase() {
	router := gin.New()

	// Event store
	router = esFactory.MakeServer(esdb.New(), router)

	// Projection
	router = prFactory.MakeServer(projectionRepo.New(), router)

	// Projector
	projectorRepo := projectorRepo.New()
	p := projector.New(projectorRepo)
	esgw := gateway.New()
	go esgw.Subscribe(func(event *domain.Event) {
		p.Project(event)
	})

	router.Run(config.IP_PORT)
}
