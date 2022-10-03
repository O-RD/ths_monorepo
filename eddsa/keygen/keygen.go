package keygen

import (
	"time"

	"github.com/O-RD/ths_monorepo/p2p"
	"github.com/O-RD/ths_monorepo/ths"
)

func Start(send_chan chan ths.Message, p *ths.P2P, receive_chan chan ths.Payload) {

	//listener runs the libp2p listener and store received values
	proceed_chan := make(chan int)
	Run_listener(p, receive_chan, proceed_chan)
	//Add another channel to listener to agree to move ahead

	go p2p.Send(send_chan)
	time.Sleep(time.Second * 2)
	for i := 0; i < len(p.Peers); i++ {

		//if p.Peers[i].Id != p.Host.ID() -> Continue
		send_chan <- ths.Message{From: *p,
			Type:         2,
			To:           p.Peers[i].Id,
			Payload_name: "First",
			Payload:      "Test",
			End:          0}

	}
	<-proceed_chan
	//compute after round and proceed - replaces wait_until()

}
