package esdb

import (
	"context"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/calogxro/qaservice/config"
	"github.com/calogxro/qaservice/domain"
)

//var streamID = "answers"

type EventStoreDB struct {
	db *esdb.Client
}

func New() *EventStoreDB {
	db, _ := connect(config.ESDB_URI)
	return &EventStoreDB{
		db: db,
	}
}

func (e *EventStoreDB) Subscribe(onEvent func(*domain.Event)) error {
	// TODO SubscribeToStream(streamID)
	stream, err := e.db.SubscribeToAll(context.Background(), esdb.SubscribeToAllOptions{})

	if err != nil {
		return err
	}

	defer stream.Close()

	for {
		event := stream.Recv()

		if event.EventAppeared != nil {
			eventType := event.EventAppeared.Event.EventType
			eventData := event.EventAppeared.Event.Data

			onEvent(domain.NewEvent(eventType, eventData))
		}

		if event.SubscriptionDropped != nil {
			break
		}
	}
	return nil
}
