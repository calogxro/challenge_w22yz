package main

import (
	"encoding/json"
)

type EventStoreStub struct {
	store   []*Event
	onEvent func(*Event)
}

func NewEventStoreStub() *EventStoreStub {
	return &EventStoreStub{
		store:   []*Event{},
		onEvent: func(*Event) {},
	}
}

func (es *EventStoreStub) Subscribe(onEvent func(*Event)) error {
	es.onEvent = onEvent
	return nil
}

func (es *EventStoreStub) GetEvents() ([]*Event, error) {
	return es.store, nil
}

func (es *EventStoreStub) AddEvent(event *Event) error {
	es.store = append(es.store, event)
	es.onEvent(event)
	return nil
}

func (es *EventStoreStub) GetHistory(key string) ([]*Event, error) {
	history := []*Event{}
	events, _ := es.GetEvents()
	for _, event := range events {
		var answer Answer
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
