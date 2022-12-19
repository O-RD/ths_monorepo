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
	go p2p.Send(send_chan)

	Round_Values := Fill_default_Keygen()

	Data := ths.Data{
		Keygen_All_Data: Round_Values,
		Sign_All_Data:   "nil",
	}

	Round1(send_chan, p, receive_chan, &Data.Keygen_All_Data)

	Values := ths.Round_Data{
		Keygen: ths.Keygen_Data{
			EPK:        Data.Keygen_All_Data.EPK, //curves.Point
			SPK:        Data.Keygen_All_Data.SPK, //kyber.Point
			KGC:        Data.Keygen_All_Data.KGC,
			Alphas:     []string{},
			Enc_shares: []ths.Encrypted_Share{},
		},
	}

	// fmt.Println("AFTER ROUND1:", Data.Keygen_Data)

	//Add another channel to listener to agree to move ahead
	// fmt.Printf("EPK:%x\n", Data.Keygen_Data.EPK.ToAffineCompressed())

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
			Payload:      Values,
			Status:       0}

	}
	for {
		if len(p.Round1) == len(p.Peers) {

			Ack_sender <- 1
			break
		}
		//time.Sleep(time.Millisecond)
	}
	for {
		if <-proceed_chan == 1 {

			break
		}
	}

	fmt.Println("End of Round 1")
	// time.Sleep(time.Second * 10)
	//compute after round and proceed - replaces wait_until()
	p.Round = 2
	fmt.Println("Starting Round 2")
	Round2(send_chan, p, receive_chan, &Data.Keygen_All_Data)
	Round3(send_chan, p, receive_chan, &Data.Keygen_All_Data)

	Values.Keygen.Alphas = Data.Keygen_All_Data.Alphas
	fmt.Println("-->>", Values.Keygen.Alphas)
	Values.Keygen.Enc_shares = Data.Keygen_All_Data.Encrypted_Shares

	for i := 0; i < len(p.Sorted_Peers); i++ {

		if i == p.My_Index {
			continue
		}
		//if p.Peers[i].Id != p.Host.ID() -> Continue
		send_chan <- ths.Message{From: *p,
			Type:         2,
			To:           p.Sorted_Peers[i].Id,
			Payload_name: "Second",
			Payload:      Values,
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

	// time.Sleep(time.Second * 10)
}
