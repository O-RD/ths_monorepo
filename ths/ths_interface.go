package ths

import (
	"context"

	host "github.com/libp2p/go-libp2p/core/host"
	peer "github.com/libp2p/go-libp2p/core/peer"
)

type THS struct {
	Id      peer.ID
	Moniker string
}

type Message struct {
	From    peer.ID
	To      peer.ID
	Payload string
}

type P2P struct {

	// Represents the libp2p host
	Host             host.Host
	Host_ip          string
	Ctx              context.Context
	Peers            []THS
	Connectedparties int
	ThisParty        int

	send       chan Message
	receive    chan Message
	Threshold  int
	Party_Size int
	Moniker    string
	Port       string
	ThsType    string
}

type Moniker_message struct {
	Name string
}
