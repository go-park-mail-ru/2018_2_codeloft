package models

import (
	"sync"
)

//easyjson:json
type Cell struct {
	Val int        `json:"id"`
	Mu  sync.Mutex `json:"-"`
}
