package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
)

type P2P struct {

	// Represents the libp2p host
	Host    host.Host
	Host_ip string
	Ctx     context.Context
	Peers   []string

	Threshold int
	Port      string
	ThsType   string
}

type Peer_details struct {
	Id   peer.ID
	Addr peer.AddrInfo
}
