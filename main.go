package main

import (
	"github.com/O-RD/ths_monorepo/eddsa/keygen"
	"github.com/O-RD/ths_monorepo/p2p"
	"github.com/O-RD/ths_monorepo/ths"
)

func main() {

	start_keygen()
}

func start_keygen() {
	// party := p2p.P2p_init()
	// var version string

	var party ths.P2P

	//Init P2P values to Party
	p2p_ch := make(chan ths.P2P)
	go func(p2p_ch chan ths.P2P) {
		ths.P2p_init(p2p_ch)
	}(p2p_ch)

	party = <-p2p_ch
	party.Port, _, party.Threshold, party.Party_Size, party.Moniker = Process_flags() //Adds port number to p2p struct

	//MDNS

	var send_chan = make(chan ths.Message)
	var receiver_ch = make(chan ths.Payload)
	//add channel)
	p2p.Create_Peer(&party)

	// ths.Sort_Peers(&party) // Adds host id and sorts peers - adds this party values to peer_list

	keygen.Start(send_chan, &party, receiver_ch)
}
