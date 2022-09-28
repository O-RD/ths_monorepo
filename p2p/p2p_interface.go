package p2p

import (
	"context"
	"sync"

	"github.com/O-RD/ths_monorepo/ths"
	"github.com/libp2p/go-libp2p/core/host"
)

type P2P struct {

	// Represents the libp2p host
	Host             host.Host
	Host_ip          string
	Ctx              context.Context
	Peers            []ths.THS
	Connectedparties int
	ThisParty        int

	send       chan ths.Message
	receive    chan ths.Message
	Threshold  int
	Party_Size int
	Moniker    string
	Port       string
	ThsType    string
}

type moniker_message struct {
	Name string
}

var l = sync.Mutex{}

type Party interface {
	Start()
}
