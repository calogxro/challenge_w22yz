package main

import (
	"context"
	"fmt"
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

type Event struct {
	Event string  `json:"event"`
	Data  *Answer `json:"data"`
}

type EventStore interface {
	init()
	create(*Answer) (interface{}, error)
	find(string) (*Answer, error)
	update(*Answer) (interface{}, error)
	delete(string) (interface{}, error)
	getHistory(string) ([]*Event, error)
	drop()
}

type MongoStore struct {
	dbUser string
	dbPass string
	dbHost string
	dbPort string

	dbName string
	dbColl string

	coll *mongo.Collection
}

func (s *MongoStore) init() {
	uri_template := "mongodb://%s:%s@%s:%s/"
	uri := fmt.Sprintf(uri_template, s.dbUser, s.dbPass, s.dbHost, s.dbPort)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	s.coll = client.Database(s.dbName).Collection(s.dbColl)
}

func (s *MongoStore) drop() {
	s.coll.Drop(context.TODO())
}

// https://stackoverflow.com/a/57777011
func (s *MongoStore) save(answer *Answer) (interface{}, error) {
	answer.ID = primitive.NewObjectID()
	result, err := s.coll.InsertOne(context.TODO(), answer)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (s *MongoStore) find(key string) (*Answer, error) {
	filter := bson.D{{Key: "key", Value: key}}
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetLimit(1)
	cursor, err := s.coll.Find(context.TODO(), filter, opts)
	if err != nil {
		//panic(err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if found := cursor.Next(context.TODO()); found {
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

func (s *MongoStore) exists(key string) (bool, error) {
	answer, err := s.find(key)
	if err != nil {
		return false, err
	}
	return answer != nil, nil
}

func (s *MongoStore) create(answer *Answer) (interface{}, error) {
	exists, err := s.exists(answer.Key)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, &KeyExists{}
	}

	answer.Event = CREATED_EVENT
	id, _ := s.save(answer)
	return id, nil
}

func (s *MongoStore) update(answer *Answer) (interface{}, error) {
	exists, err := s.exists(answer.Key)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, &KeyNotFound{}
	}

	answer.Event = UPDATED_EVENT
	id, _ := s.save(answer)
	return id, nil
}

func (s *MongoStore) delete(key string) (interface{}, error) {
	answer, err := s.find(key)
	if err != nil {
		return nil, err
	}

	if exists := answer != nil; !exists {
		return nil, &KeyNotFound{}
	}

	answer.Event = DELETED_EVENT
	id, _ := s.save(answer)
	return id, nil
}

func (s *MongoStore) getHistory(key string) ([]*Event, error) {
	filter := bson.D{{Key: "key", Value: key}}
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: 1}})
	cursor, err := s.coll.Find(context.TODO(), filter, opts)
	if err != nil {
		//panic(err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []*Answer
	if err = cursor.All(context.TODO(), &results); err != nil {
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
