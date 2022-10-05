package keygen

import (
	"fmt"
	"time"

	"github.com/O-RD/ths_monorepo/p2p"
	"github.com/O-RD/ths_monorepo/ths"
)

func Start(send_chan chan ths.Message, p *ths.P2P, receive_chan chan ths.Payload) {

	//listener runs the libp2p listener and store received values
	proceed_chan := make(chan int)
	Ack_sender := make(chan int)
	p.Round = 1
	go Run_listener(p, receive_chan, proceed_chan, Ack_sender)
	//Add another channel to listener to agree to move ahead

	go p2p.Send(send_chan)
	fmt.Println("Initiate Keygen")
	fmt.Println("Starting Round 1")
	time.Sleep(time.Second * 3)

	for i := 0; i < len(p.Sorted_Peers); i++ {

		if i == p.My_Index {
			continue
		}
		//if p.Peers[i].Id != p.Host.ID() -> Continue
		send_chan <- ths.Message{From: *p,
			Type:         1,
			To:           p.Sorted_Peers[i].Id,
			Payload_name: "First",
			Payload:      "Test",
			Status:       0}

	}
	for {
		if len(p.Round1) == len(p.Peers) {
			Ack_sender <- 1
			break
		}
		time.Sleep(time.Millisecond)
	}
	for {
		if <-proceed_chan == 1 {
			break
		}
	}
	fmt.Println("End of Round 1")
	//compute after round and proceed - replaces wait_until()
	p.Round = 2
	fmt.Println("Starting Round 2")
	for i := 0; i < len(p.Sorted_Peers); i++ {

		if i == p.My_Index {
			continue
		}
		//if p.Peers[i].Id != p.Host.ID() -> Continue
		send_chan <- ths.Message{From: *p,
			Type:         2,
			To:           p.Sorted_Peers[i].Id,
			Payload_name: "Second",
			Payload:      "Test2",
			Status:       0}

	}
	for {
		if len(p.Round2) == len(p.Peers) {
			Ack_sender <- 2
			break
		}
		time.Sleep(time.Second * 2)
	}
	for {
		if <-proceed_chan == 2 {
			break
		}
	}
	fmt.Println("End of Round 2")
}
