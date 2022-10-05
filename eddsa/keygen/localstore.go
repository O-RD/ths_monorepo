package keygen

import (
	"fmt"

	"github.com/O-RD/ths_monorepo/p2p"
	"github.com/O-RD/ths_monorepo/ths"
)

func Run_listener(p *ths.P2P, receive_chan chan ths.Payload, proceed chan int, Ack_sender chan int) {

	p2p.Input_Stream_listener(p, receive_chan)
	p2p.Acknowledgement_listener(p, proceed)

	go p2p.Send_Ack(p, Ack_sender)
	for {
		// Figure how to store - use
		temp := <-receive_chan
		if temp.Type == 1 {
			p.Round1 = append(p.Round1, ths.Keygen_Store_Round1{Id: temp.Sender,
				V1:  temp.Payload,
				Ack: 0,
			})

			// if len(p.Round1) == len(p.Peers) {
			// 	Ack_sender <- 1
			// }

		} else if temp.Type == 2 {
			p.Round2 = append(p.Round2, ths.Keygen_Store_Round2{Id: temp.Sender,
				V1:  temp.Payload,
				Ack: 0,
			})

			// if len(p.Round2) == len(p.Peers) {
			// 	Ack_sender <- 2
			// }

		}
		fmt.Println(temp)

	}
}
