package db

import (
	"database/sql"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/go-sql-driver/mysql"
)

func InitMySQL() (*sql.DB, error) {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   "username", //os.Getenv("DBUSER"),
		Passwd: "password", //os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "qaservice",
	}
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

func initESDB() (*esdb.Client, error) {
	uri := "esdb://127.0.0.1:2113?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000"

	settings, err := esdb.ParseConnectionString(uri)
	if err != nil {
		return nil, err
	}

	db, err := esdb.NewClient(settings)
	if err != nil {
		return nil, err
	}

	return db, nil
}
