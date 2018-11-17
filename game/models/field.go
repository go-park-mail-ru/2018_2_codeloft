package models

import "sync"

type Cell struct {
	Val int
	sync.Mutex
}
