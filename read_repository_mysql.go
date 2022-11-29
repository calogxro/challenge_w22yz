package main

import (
	"database/sql"
	"fmt"
)

type MySQLReadRepository struct {
	db *sql.DB
}

func NewMySQLReadRepository() *MySQLReadRepository {
	db, _ := initMySQL()
	return &MySQLReadRepository{
		db: db,
	}
}

func (rr *MySQLReadRepository) GetAnswer(key string) (*Answer, error) {
	var answer Answer

	row := rr.db.QueryRow("SELECT key_, value FROM answer WHERE key_ = ?", key)
	if err := row.Scan(&answer.Key, &answer.Value); err != nil {
		if err == sql.ErrNoRows {
			return nil, &KeyNotFound{}
		}
		return nil, fmt.Errorf("GetAnswer %s: %v", key, err)
	}
	return &answer, nil
}

func (rr *MySQLReadRepository) CreateAnswer(answer Answer) error {
	stmt := "INSERT INTO answer (key_, value) VALUES (?, ?)"
	_, err := rr.db.Exec(stmt, answer.Key, answer.Value)
	return err
}

func (rr *MySQLReadRepository) UpdateAnswer(answer Answer) error {
	stmt := "UPDATE answer SET value = ? WHERE key_ = ?"
	_, err := rr.db.Exec(stmt, answer.Value, answer.Key)
	return err
}

func (rr *MySQLReadRepository) DeleteAnswer(answer Answer) error {
	stmt := "DELETE FROM answer WHERE key_ = ?"
	_, err := rr.db.Exec(stmt, answer.Key)
	return err
}
