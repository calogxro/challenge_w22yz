package main

import "encoding/json"

const (
	ANSWER_CREATED_EVENT = "AnswerCreatedEvent"
	ANSWER_UPDATED_EVENT = "AnswerUpdatedEvent"
	ANSWER_DELETED_EVENT = "AnswerDeletedEvent"
)

type Event struct {
	Type string //`json:"type"`
	Data []byte //`json:"data"`
}

func NewEvent(eventType string, eventData []byte) *Event {
	return &Event{
		Type: eventType,
		Data: eventData,
	}
}

func NewAnswerCreatedEvent(answer Answer) (*Event, error) {
	data, err := json.Marshal(answer)
	return NewEvent(ANSWER_CREATED_EVENT, data), err
}

func NewAnswerUpdatedEvent(answer Answer) (*Event, error) {
	data, err := json.Marshal(answer)
	return NewEvent(ANSWER_UPDATED_EVENT, data), err
}

func NewAnswerDeletedEvent(answer Answer) (*Event, error) {
	data, err := json.Marshal(answer)
	return NewEvent(ANSWER_DELETED_EVENT, data), err
}
