package game

import (
	gamemodels "github.com/go-park-mail-ru/2018_2_codeloft/game/game/models"
	"sync"
)

type DiffCell struct {
	Pos gamemodels.Position `json:"pos"`
	Val string              `json:"color"`
}

type Diff struct {
	sync.Mutex `json:"-"`
	DiffArray  []DiffCell `json:"diff_array"`
}

func (d *Diff) Add(cell DiffCell) {
	d.Lock()
	d.DiffArray = append(d.DiffArray, cell)
	d.Unlock()
}
