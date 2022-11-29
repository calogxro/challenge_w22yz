package db

import (
	"database/sql"
	"fmt"

	"github.com/calogxro/qaservice/domain"
)

type MySQLReadRepository struct {
	db *sql.DB
}

func NewMySQLReadRepository() *MySQLReadRepository {
	db, _ := InitMySQL()
	return &MySQLReadRepository{
		db: db,
	}
}

func (rr *MySQLReadRepository) GetAnswer(key string) (*domain.Answer, error) {
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

func (rr *MySQLReadRepository) CreateAnswer(answer domain.Answer) error {
	stmt := "INSERT INTO answer (key_, value) VALUES (?, ?)"
	_, err := rr.db.Exec(stmt, answer.Key, answer.Value)
	return err
}

func (rr *MySQLReadRepository) UpdateAnswer(answer domain.Answer) error {
	stmt := "UPDATE answer SET value = ? WHERE key_ = ?"
	_, err := rr.db.Exec(stmt, answer.Value, answer.Key)
	return err
}

func (rr *MySQLReadRepository) DeleteAnswer(answer domain.Answer) error {
	stmt := "DELETE FROM answer WHERE key_ = ?"
	_, err := rr.db.Exec(stmt, answer.Key)
	return err
}
