package main

import (
	"flag"
	"net/http"

	"github.com/calogxro/qaservice/config"
	"github.com/calogxro/qaservice/domain"
	es_httphandler "github.com/calogxro/qaservice/eventstore/handler/http"
	"github.com/calogxro/qaservice/eventstore/repository/esdb"
	esmemdb "github.com/calogxro/qaservice/eventstore/repository/memory"
	rrmemdb "github.com/calogxro/qaservice/projection/repository/memory"

	"github.com/calogxro/qaservice/eventstore/service/eventstore"
	esdbgateway "github.com/calogxro/qaservice/projection/gateway/esdb"
	pr_httphandler "github.com/calogxro/qaservice/projection/handler/http"
	"github.com/calogxro/qaservice/projection/repository/mongodb"
	"github.com/calogxro/qaservice/projection/service/projection"
	"github.com/calogxro/qaservice/projection/service/projector"
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
	es := esmemdb.New()
	rr := rrmemdb.New()
	p := projector.New(rr)

	es.Subscribe(func(event *domain.Event) {
		p.Project(event)
	})

	es_service := eventstore.New(es)
	pr_service := projection.New(rr)

	router := gin.New()
	router.GET("/ping", ping)
	router = es_httphandler.MakeHandler(es_service, router)
	router = pr_httphandler.MakeHandler(pr_service, router)

	//router.Run(":8080")
	router.Run(config.IP_PORT)
}

func runWithDatabase() {
	repo := esdb.New()
	readRepository := mongodb.New()
	p := projector.New(readRepository)

	esgw := esdbgateway.New()
	go esgw.Subscribe(func(event *domain.Event) {
		p.Project(event)
	})

	es_service := eventstore.New(repo)
	pr_service := projection.New(readRepository)

	router := gin.New()
	router.GET("/ping", ping)
	router = es_httphandler.MakeHandler(es_service, router)
	router = pr_httphandler.MakeHandler(pr_service, router)

	//router.Run(":8080")
	router.Run(config.IP_PORT)
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
