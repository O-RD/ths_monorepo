package p2p

import (
	"sync"
)

var l = sync.Mutex{}

type Party interface {
	Start()
}
