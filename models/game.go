package models

import (
	"database/sql"
	"log"
)

type Game struct {
	Score   int64
	Game_id int64
}

func (g *Game) UpdateScore(db *sql.DB) error {
	_, err := db.Exec("update game set score=$1 where id = $2", g.Score, g.Game_id)
	if err != nil {
		log.Printf("cant UpdateScore: %v\n", g)
		return err
	}
	return nil
}
