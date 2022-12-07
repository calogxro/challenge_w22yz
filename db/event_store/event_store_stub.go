package db

import (
	"encoding/json"

	"github.com/calogxro/qaservice/domain"
)

type EventStoreStub struct {
	store   []*domain.Event
	onEvent func(*domain.Event)
}

func NewEventStoreStub() *EventStoreStub {
	return &EventStoreStub{
		store:   []*domain.Event{},
		onEvent: func(*domain.Event) {},
	}
}

func (es *EventStoreStub) Subscribe(onEvent func(*domain.Event)) error {
	es.onEvent = onEvent
	return nil
}

func (es *EventStoreStub) GetEvents() ([]*domain.Event, error) {
	return es.store, nil
}

func (es *EventStoreStub) AddEvent(event *domain.Event) error {
	es.store = append(es.store, event)
	es.onEvent(event)
	return nil
}

func (es *EventStoreStub) GetHistory(key string) ([]*domain.Event, error) {
	history := []*domain.Event{}
	events, _ := es.GetEvents()
	for _, event := range events {
		var answer domain.Answer
		err := json.Unmarshal([]byte(event.Data), &answer) // TODO

		if err != nil {
			return nil, err
		}

		if answer.Key == key {
			history = append(history, event)
		}
	}
	//fmt.Printf("history: %v\n", history)
	return history, nil
}
