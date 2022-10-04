package keygen

import (
	"fmt"

	"github.com/O-RD/ths_monorepo/p2p"
	"github.com/O-RD/ths_monorepo/ths"
)

func Run_listener(p *ths.P2P, receive_chan chan ths.Payload, proceed chan int) {

	p2p.Input_Stream_listener(p, receive_chan)
	for {
		// Figure how to store - use
		temp := <-receive_chan
		fmt.Println(temp)
		var flag = 0
		for i := 0; i < len(p.Peers); i++ {
			if p.Peers[i].Id == temp.Sender {
				p.Peers[i].Round = temp.Type
			}
			if p.Peers[i].Round != p.Round {
				flag = 1
			}
		}
		if flag == 0 {
			proceed <- 1
			fmt.Println("End of Round", p.Round)
		}

	}
}
