package keygen

import (
	"fmt"

	"github.com/O-RD/ths_monorepo/p2p"
	"github.com/O-RD/ths_monorepo/ths"
)

func Start(send_chan chan ths.Message, p p2p.P2P) {

	go p2p.Send(send_chan)
	for i := 0; i < len(p.Peers); i++ {
		fmt.Println(p.Peers[i].Id)
		send_chan <- ths.Message{From: p.Host.ID(),
			To:      p.Peers[i].Id,
			Payload: "Test"}

	}
}
