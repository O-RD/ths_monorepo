package p2p

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
)

type P2P struct {

	// Represents the libp2p host
	Host             host.Host
	Host_ip          string
	Ctx              context.Context
	Peers            []peer_names
	Connectedparties int

	Threshold  int
	Party_Size int
	Moniker    string
	Port       string
	ThsType    string
}

type Peer_details struct {
	Id   peer.ID
	Addr peer.AddrInfo
}

type moniker_message struct {
	Name string
}

type peer_names struct {
	Id   peer.ID
	Name string
}

var l = sync.Mutex{}
