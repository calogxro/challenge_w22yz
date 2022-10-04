package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Answer struct {
	ID    primitive.ObjectID `json:"-" bson:"_id"`
	Key   string             `json:"key" binding:"required"`
	Value string             `json:"value" binding:"required"`
	Event string             `json:"-"`
}
