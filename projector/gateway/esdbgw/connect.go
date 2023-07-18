package gateway

import (
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	//"github.com/calogxro/qaservice/config"
)

func connect(URI string) (*esdb.Client, error) {
	settings, err := esdb.ParseConnectionString(URI)
	if err != nil {
		return nil, err
	}

	db, err := esdb.NewClient(settings)
	if err != nil {
		return nil, err
	}

	return db, nil
}
