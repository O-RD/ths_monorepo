package ths

import (
	"context"

	"github.com/libp2p/go-libp2p-core/host"
)

type P2P struct {

	// Represents the libp2p host
	Host    host.Host
	Host_ip string
	Ctx     context.Context
	Peers   []string
}
