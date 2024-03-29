package main

import (
	"time"

	"github.com/O-RD/ths_monorepo/eddsa/keygen"
	"github.com/O-RD/ths_monorepo/p2p"
	"github.com/O-RD/ths_monorepo/ths"
)

func main() {

	start_keygen()

}

func start_keygen() {

	var party ths.P2P

	//Init P2P values to Party
	p2p_ch := make(chan ths.P2P)
	go func(p2p_ch chan ths.P2P) {
		ths.P2p_init(p2p_ch)
	}(p2p_ch)

	party = <-p2p_ch
	party.Topic, _, party.Threshold, party.Party_Size, party.Moniker = Process_flags()

	//MDNS

	var send_chan = make(chan ths.Message)
	var receiver_ch = make(chan ths.Payload)

	// ths.Find_peers_api(&party)
	p2p.Create_Peer(&party)
	time.Sleep(time.Second)
	ths.Sort_Peers(&party) // Adds host id and sorts peers - adds this party values to peer_list

	//Performs Base Keygen - and Sample Signature for interested users
	keygen.Start(send_chan, &party, receiver_ch)
}
