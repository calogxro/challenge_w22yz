package mongodb

import (
	"context"
	"log"

	"github.com/calogxro/qaservice/config"
	"github.com/calogxro/qaservice/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	coll *mongo.Collection
}

func New() *MongoRepository {
	cfg := config.MongoDB
	client, err := InitMongoDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database(cfg.DBName).Collection(cfg.DBColl)
	return &MongoRepository{
		coll: coll,
	}
}

func (s *MongoRepository) DeleteAllAnswers() {
	s.coll.Drop(context.TODO())
}

func (s *MongoRepository) GetAnswer(key string) (*domain.Answer, error) {
	filter := bson.D{{Key: "key", Value: key}}
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetLimit(1)
	cursor, err := s.coll.Find(context.TODO(), filter, opts)
	if err != nil {
		//panic(err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if found := cursor.Next(context.TODO()); found {
		var answer domain.Answer
		if err := cursor.Decode(&answer); err != nil {
			//log.Fatal(err)
			return nil, err
		}
		return &answer, nil
	}
	return nil, domain.ErrKeyNotFound
}

func (s *MongoRepository) CreateAnswer(answer domain.Answer) error {
	_, err := s.coll.InsertOne(context.TODO(), answer)
	return err
}

func (s *MongoRepository) UpdateAnswer(answer domain.Answer) error {
	filter := bson.D{{Key: "key", Value: answer.Key}}
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "value", Value: answer.Value}}},
	}
	_, err := s.coll.UpdateOne(context.TODO(), filter, update)
	return err
}

func (s *MongoRepository) DeleteAnswer(answer domain.Answer) error {
	filter := bson.D{{Key: "key", Value: answer.Key}}
	_, err := s.coll.DeleteOne(context.TODO(), filter)
	return err
}
