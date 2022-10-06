package keygen

import (
	"fmt"
	"os"
	"strconv"

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
			x := fmt.Sprintf("%x", temp.Payload.Keygen_Data.EPK.ToAffineCompressed())
			fmt.Println("EPK:", x)
			Sender_index := p2p.GetIndex(p.Sorted_Peers, temp.Sender)
			fmt.Println("Sender:", Sender_index)

			err1 := os.MkdirAll("Broadcast/"+strconv.Itoa(Sender_index), os.ModePerm)
			if err1 != nil {
				fmt.Println("Error")
			}
			_f, err := os.OpenFile("Broadcast/"+strconv.Itoa(Sender_index)+"/EPK.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
			if err != nil {
				fmt.Println("ERROR:")
			}

			_f.WriteString(x)

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
