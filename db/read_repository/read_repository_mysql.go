package readrepository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/calogxro/qaservice/config"
	"github.com/calogxro/qaservice/db"
	"github.com/calogxro/qaservice/domain"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository() *MySQLRepository {
	db, err := db.InitMySQL(config.MySQL)
	if err != nil {
		log.Fatal(err)
	}
	return &MySQLRepository{
		db: db,
	}
}

func (rr *MySQLRepository) GetAnswer(key string) (*domain.Answer, error) {
	var answer domain.Answer

	row := rr.db.QueryRow("SELECT key_, value FROM answer WHERE key_ = ?", key)
	if err := row.Scan(&answer.Key, &answer.Value); err != nil {
		if err == sql.ErrNoRows {
			return nil, &domain.KeyNotFound{}
		}
		return nil, fmt.Errorf("GetAnswer %s: %v", key, err)
	}
	return &answer, nil
}

func (rr *MySQLRepository) CreateAnswer(answer domain.Answer) error {
	stmt := "INSERT INTO answer (key_, value) VALUES (?, ?)"
	_, err := rr.db.Exec(stmt, answer.Key, answer.Value)
	return err
}

func (rr *MySQLRepository) UpdateAnswer(answer domain.Answer) error {
	stmt := "UPDATE answer SET value = ? WHERE key_ = ?"
	_, err := rr.db.Exec(stmt, answer.Value, answer.Key)
	return err
}

func (rr *MySQLRepository) DeleteAnswer(answer domain.Answer) error {
	stmt := "DELETE FROM answer WHERE key_ = ?"
	_, err := rr.db.Exec(stmt, answer.Key)
	return err
}

func (rr *MySQLRepository) DeleteAllAnswers() error {
	stmt := "DELETE FROM answer"
	_, err := rr.db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
