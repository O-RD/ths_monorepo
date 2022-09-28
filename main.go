package main

import (
	"fmt"
	"log"

	"github.com/O-RD/ths_monorepo/eddsa/keygen"
	p2p "github.com/O-RD/ths_monorepo/p2p"
	"github.com/O-RD/ths_monorepo/ths"
)

func main() {

	// party := p2p.P2p_init()
	// var version string
	var party p2p.P2P

	//Init P2P values to Party
	p2p_ch := make(chan p2p.P2P)
	go func(p2p_ch chan p2p.P2P) {
		p2p.P2p_init(p2p_ch)
	}(p2p_ch)

	party = <-p2p_ch
	party.Port, _, party.Threshold, party.Party_Size, party.Moniker = Process_flags() //Adds port number to p2p struct
	// fmt.Println(version)
	//MDNS

	party.Create_peer()
	Sort_Peers(&party) // Adds host id and sorts peers - adds this party values to peer_list

	log.Println(party)

	fmt.Println("here")
	var send_chan = make(chan ths.Message)
	keygen.Start(send_chan, party)
}
