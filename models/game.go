package models

import (
	"database/sql"

	"go.uber.org/zap"
)

type Game struct {
	Score   int64
	Game_id int64
}

func (g *Game) UpdateScore(db *sql.DB) error {
	_, err := db.Exec("update game set score=$1 where id = $2", g.Score, g.Game_id)
	if err != nil {
		zap.L().Warn("Can not update score",
			zap.Error(err))
		return err
	}
	return nil
}
