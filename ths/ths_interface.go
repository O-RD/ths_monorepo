package ths

import peer "github.com/libp2p/go-libp2p/core/peer"

type THS struct {
	Id      peer.ID
	Moniker string
}

type Message struct {
	From    peer.ID
	To      peer.ID
	Payload string
}
