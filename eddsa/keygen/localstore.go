package keygen

import (
	"fmt"

	"github.com/O-RD/ths_monorepo/p2p"
	"github.com/O-RD/ths_monorepo/ths"
)

func Run_listener(p *ths.P2P, receive_chan chan ths.Payload, proceed chan int) {

	p2p.Input_Stream_listener(p, receive_chan)
	for {
		// Figure how to store
		fmt.Println(<-receive_chan)
	}
}
