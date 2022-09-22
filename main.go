package main

import (
	"fmt"

	p2p "github.com/O-RD/ths_monorepo/p2p"
)

func main() {

	party := p2p.P2p_init()
	var version string
	party.Port, version, party.Threshold = process_flags() //Adds port number to p2p struct
	fmt.Println(party, version)
	//MDNS
	party.Create_peer()
}
