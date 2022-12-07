package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/calogxro/qaservice/config"

	//"github.com/calogxro/qaservice/config"
	"github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMySQL(cfg mysql.Config) (*sql.DB, error) {
	// Get a database handle.
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		//log.Fatal(pingErr)
		return nil, err
	}
	//fmt.Println("Connected!")

	return db, nil
}

func InitESDB(URI string) (*esdb.Client, error) {
	settings, err := esdb.ParseConnectionString(URI)
	if err != nil {
		return nil, err
	}

	db, err := esdb.NewClient(settings)
	if err != nil {
		return nil, err
	}

	return db, nil
}

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
