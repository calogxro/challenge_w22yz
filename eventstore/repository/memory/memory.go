package memory

import (
	"encoding/json"

	"github.com/calogxro/qaservice/domain"
)

type EventStore struct {
	store   []*domain.Event
	onEvent func(*domain.Event)
}

func New() *EventStore {
	return &EventStore{
		store:   []*domain.Event{},
		onEvent: func(*domain.Event) {},
	}
}

func (es *EventStore) Subscribe(onEvent func(*domain.Event)) error {
	es.onEvent = onEvent
	return nil
}

func (es *EventStore) AddEvent(event *domain.Event) error {
	es.store = append(es.store, event)
	es.onEvent(event)
	return nil
}

// GetEvents returns all the events in the store.
func (es *EventStore) GetEvents() ([]*domain.Event, error) {
	return es.store, nil
}

// GetHistory returns all the events for a given key.
func (es *EventStore) GetHistory(key string) ([]*domain.Event, error) {
	history := []*domain.Event{}
	events, _ := es.GetEvents()
	for _, event := range events {
		var answer domain.Answer
		//err := json.Unmarshal([]byte(event.Data), &answer) // TODO
		err := json.Unmarshal(event.Data, &answer) // TODO

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
