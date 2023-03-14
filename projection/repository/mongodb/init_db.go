package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/calogxro/qaservice/config"

	//"github.com/calogxro/qaservice/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB(cfg config.MongoConfig) (*mongo.Client, error) {
	uri_template := "mongodb://%s:%s@%s:%s/"
	uri := fmt.Sprintf(uri_template, cfg.User, cfg.Pass, cfg.Host, cfg.Port)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return client, nil
}
