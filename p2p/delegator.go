package p2p

import (
	"fmt"
	"time"

	ths "github.com/O-RD/ths_monorepo/ths"
)

func Create_Peer(p *ths.P2P) {

	addedpeer := make(chan ths.THS)
	Connection_Stream_listener(p, addedpeer)
	go ths.Create_peer(p)
	for {
		fmt.Println("Connected to", <-addedpeer)
		p.Connectedparties += 1
		// fmt.Println(p.Connectedparties, p.Party_Size-1, len(p.Peers))
		if len(p.Peers) >= p.Party_Size-1 && p.Connectedparties >= p.Party_Size-1 {
			time.Sleep(time.Second * 2)
			break
		}
	}
}
