package db

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/calogxro/qaservice/config"
	"github.com/calogxro/qaservice/db"
	"github.com/calogxro/qaservice/domain"
)

var streamID = "answers"

type EventStoreDB struct {
	db *esdb.Client
}

func NewEventStoreDB() *EventStoreDB {
	db, _ := db.InitESDB(config.ESDB_URI)
	return &EventStoreDB{
		db: db,
	}
}

func (e *EventStoreDB) GetEvents() ([]*domain.Event, error) {
	options := esdb.ReadStreamOptions{
		From:      esdb.Start{},
		Direction: esdb.Forwards,
	}

	stream, err := e.db.ReadStream(context.Background(), streamID, options, 100)

	if err != nil {
		return nil, err
	}

	defer stream.Close()

	var events = []*domain.Event{}

	for {
		resolvedEvent, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			//panic(err) //TODO: panic: : stream 'answers' is not found
			return nil, err
		}

		eventType := resolvedEvent.Event.EventType
		eventData := resolvedEvent.Event.Data

		events = append(events, domain.NewEvent(eventType, eventData))
	}

	return events, nil
}

func (e *EventStoreDB) AddEvent(event *domain.Event) error {
	eventData := esdb.EventData{
		ContentType: esdb.ContentTypeJson,
		EventType:   event.Type,
		Data:        event.Data,
	}

	_, err := e.db.AppendToStream(context.Background(), streamID, esdb.AppendToStreamOptions{}, eventData)

	return err
}

func (e *EventStoreDB) Subscribe(onEvent func(*domain.Event)) error {
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

func (es *EventStoreDB) GetHistory(key string) ([]*domain.Event, error) {
	var history []*domain.Event
	events, _ := es.GetEvents()
	for _, event := range events {
		var answer domain.Answer
		json.Unmarshal([]byte(event.Data), &answer)

		if answer.Key == key {
			history = append(history, event)
		}
	}
	return history, nil
}

func (s *EventStoreDB) DeleteStream() error {
	opts := esdb.DeleteStreamOptions{
		//ExpectedRevision: esdb.Revision(0),
	}
	_, err := s.db.DeleteStream(context.Background(), streamID, opts)
	if err != nil {
		return err
	}
	//fmt.Printf("Drop %+v\n\n", deleteResult)
	return nil
}
