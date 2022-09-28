package main

import (
	p2p "github.com/O-RD/ths_monorepo/p2p"
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

}
