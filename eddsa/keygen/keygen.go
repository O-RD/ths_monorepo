package keygen

import (
	"fmt"
	"time"

	"github.com/O-RD/ths_monorepo/p2p"
	"github.com/O-RD/ths_monorepo/ths"
)

func Start(send_chan chan ths.Message, p *ths.P2P, receive_chan chan ths.Payload) {

	p2p.Input_Stream_listener(p, receive_chan)

	go p2p.Send(send_chan)
	time.Sleep(time.Second * 2)
	for i := 0; i < len(p.Peers); i++ {

		send_chan <- ths.Message{From: *p,
			Type:    2,
			To:      p.Peers[i].Id,
			Payload: "Test",
			End:     0}

	}

	for {
		fmt.Println(<-receive_chan)
	}

}
