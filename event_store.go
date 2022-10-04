package main

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CREATED_EVENT = "created"
	UPDATED_EVENT = "updated"
	DELETED_EVENT = "deleted"
)

type IEventStore interface {
	create(*Answer) (interface{}, error)
	find(string) (*Answer, error)
	update(*Answer) (interface{}, error)
	delete(string) (interface{}, error)
	getHistory(string) ([]*Event, error)
	drop()
}

type EventStore struct {
	coll *mongo.Collection
}

func connect(db_uri string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(db_uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func (es *EventStore) drop() {
	es.coll.Drop(ctx)
}

// https://stackoverflow.com/a/57777011
func (es *EventStore) save(answer *Answer) (interface{}, error) {
	answer.ID = primitive.NewObjectID()
	result, err := es.coll.InsertOne(ctx, answer)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (es *EventStore) find(key string) (*Answer, error) {
	filter := bson.D{{Key: "key", Value: key}}
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetLimit(1)
	cursor, err := es.coll.Find(ctx, filter, opts)
	if err != nil {
		//panic(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	if found := cursor.Next(ctx); found {
		var answer Answer
		if err := cursor.Decode(&answer); err != nil {
			//log.Fatal(err)
			return nil, err
		}
		if answer.Event != DELETED_EVENT {
			return &answer, nil
		}
	}

	// not found or deleted
	return nil, nil
}

func (es *EventStore) exists(key string) (bool, error) {
	answer, err := es.find(key)
	if err != nil {
		return false, err
	}
	return answer != nil, nil
}

func (es *EventStore) create(answer *Answer) (interface{}, error) {
	exists, err := es.exists(answer.Key)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, &KeyExists{}
	}

	answer.Event = CREATED_EVENT
	id, _ := es.save(answer)
	return id, nil
}

func (es *EventStore) update(answer *Answer) (interface{}, error) {
	exists, err := es.exists(answer.Key)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, &KeyNotFound{}
	}

	answer.Event = UPDATED_EVENT
	id, _ := es.save(answer)
	return id, nil
}

func (es *EventStore) delete(key string) (interface{}, error) {
	answer, err := es.find(key)
	if err != nil {
		return nil, err
	}

	if exists := answer != nil; !exists {
		return nil, &KeyNotFound{}
	}

	answer.Event = DELETED_EVENT
	id, _ := es.save(answer)
	return id, nil
}

func (es *EventStore) getHistory(key string) ([]*Event, error) {
	filter := bson.D{{Key: "key", Value: key}}
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: 1}})
	cursor, err := es.coll.Find(ctx, filter, opts)
	if err != nil {
		//panic(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*Answer
	if err = cursor.All(ctx, &results); err != nil {
		//log.Fatal(err)
		return nil, err
	}

	var events []*Event
	for _, answer := range results {
		events = append(events, &Event{
			Event: answer.Event,
			Data:  answer,
		})
	}

	return events, nil
}
