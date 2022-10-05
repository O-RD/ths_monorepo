package p2p

import (
	"time"

	ths "github.com/O-RD/ths_monorepo/ths"
)

func Create_Peer(p *ths.P2P) {

	addedpeer := make(chan ths.THS)
	Connection_Stream_listener(p, addedpeer)
	time.Sleep(time.Second * 3)
	go ths.Create_peer(p)
	for {
		<-addedpeer
		time.Sleep(time.Second * 3)
		// fmt.Println(p.Connectedparties, p.Party_Size-1, len(p.Peers))
		if len(p.Peers) >= p.Party_Size-1 && p.Connectedparties >= p.Party_Size-1 {

			// time.Sleep(time.Second * 5)
			break
		}
	}
}
