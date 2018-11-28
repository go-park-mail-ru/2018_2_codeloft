package models

import (
	"sync"
)

type Cell struct {
	Val int        `json:"id"`
	Mu  sync.Mutex `json:"-"`
}
