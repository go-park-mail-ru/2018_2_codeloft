package models

import (
	"database/sql"
)

type Session struct {
	Value   string
	User_id int64
}

func (s *Session) AddCookie(db *sql.DB) error {
	_, err := db.Exec("insert into sessions(value, id) values ($1, $2) on CONFLICT do nothing", s.Value, s.User_id)
	if err != nil {

		log.Printf("cant AddCookie: %v\n", s)

		return err
	}
	return nil
}

func (s *Session) DeleteCookie(db *sql.DB) error {
	_, err := db.Exec("delete from sessions where id = $1", s.User_id)
	if err != nil {

		log.Printf("cant AddCookie: %v\n", s)

		return err
	}
	return nil
}

func (s *Session) CheckCookie(db *sql.DB, value string) bool {
	row := db.QueryRow("select * from sessions where value = $1", value)
	err := row.Scan(&s.Value, &s.User_id)
	if err != nil {
		return false
	}
	return true
}
