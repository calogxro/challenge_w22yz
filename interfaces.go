package main

type IEventStore interface {
	GetEvents() ([]*Event, error)
	AddEvent(event *Event) error
	GetHistory(key string) ([]*Event, error)
	Subscribe(onEvent func(*Event)) error
}

type IReadRepository interface {
	GetAnswer(key string) (*Answer, error)
	CreateAnswer(answer Answer) error
	UpdateAnswer(answer Answer) error
	DeleteAnswer(answer Answer) error
}
