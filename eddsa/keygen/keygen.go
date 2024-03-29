package keygen

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/O-RD/ths_monorepo/p2p"
	"github.com/O-RD/ths_monorepo/ths"
)

func Start(send_chan chan ths.Message, p *ths.P2P, receive_chan chan ths.Payload) {

	fmt.Println("Passcode For Key Storage:")
	var Passcode string
	fmt.Scanln(&Passcode)
	PC := Generate_Passcode(Passcode)
	os.MkdirAll("Data/"+strconv.Itoa(p.My_Index+1), 0777)
	file, _ := os.Create("Data/" + strconv.Itoa(p.My_Index+1) + "/PC.txt")
	file.Write(PC)

	time.Sleep(time.Second * 3)

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
			KGC_sign:   ths.KGC{},
		},
	}

	//Add another channel to listener to agree to move ahead
	fmt.Println("Initiate Keygen")
	fmt.Println("Starting Round 1")
	time.Sleep(time.Second * 2)

	for i := 0; i < len(p.Sorted_Peers); i++ {

		if i == p.My_Index {
			continue
		}
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
	}
	for {
		if <-proceed_chan == 1 {

			break
		}
	}

	fmt.Println("End of Round 1")
	// time.Sleep(time.Second * 10)

	p.Round = 2
	fmt.Println("Starting Round 2")

	//Below Rounds shall Be combined
	Round2(send_chan, p, receive_chan, &Data.Keygen_All_Data)
	Round3(send_chan, p, receive_chan, &Data.Keygen_All_Data)

	Values.Keygen.Alphas = Data.Keygen_All_Data.Alphas
	fmt.Println("-->>", Values.Keygen.Alphas)
	Values.Keygen.Enc_shares = Data.Keygen_All_Data.Encrypted_Shares

	for i := 0; i < len(p.Sorted_Peers); i++ {

		if i == p.My_Index {
			continue
		}
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
	fmt.Println("End of Round 3")
	p.Round = 3

	//Below Rounds Will be Combined
	Round4(send_chan, p, receive_chan, &Data.Keygen_All_Data)
	Round5(send_chan, p, receive_chan, &Data.Keygen_All_Data)
	Values.Keygen.KGC_sign = Data.Keygen_All_Data.KGC_sign
	fmt.Println("HRELRE")
	for i := 0; i < len(p.Sorted_Peers); i++ {

		if i == p.My_Index {
			continue
		}
		send_chan <- ths.Message{From: *p,
			Type:         3,
			To:           p.Sorted_Peers[i].Id,
			Payload_name: "Third",
			Payload:      Values,
			Status:       0}
	}
	for {
		if len(p.Round3) == len(p.Peers) {
			Ack_sender <- 3
			break
		}
		time.Sleep(time.Second * 2)
	}
	for {
		if <-proceed_chan == 3 {
			break
		}
	}
	p.Round = 4
	time.Sleep(time.Second * 1)

	Round6(send_chan, p, receive_chan, &Data.Keygen_All_Data)
	Values.Keygen.Alphas_sign = Data.Keygen_All_Data.Alphas_sign
	Values.Keygen.Shares_sign = Data.Keygen_All_Data.Shares_sign
	for i := 0; i < len(p.Sorted_Peers); i++ {

		if i == p.My_Index {
			continue
		}
		send_chan <- ths.Message{From: *p,
			Type:         4,
			To:           p.Sorted_Peers[i].Id,
			Payload_name: "Forth",
			Payload:      Values,
			Status:       0}

	}
	for {
		if len(p.Round3) == len(p.Peers) {
			Ack_sender <- 4
			break
		}
		time.Sleep(time.Second * 2)
	}
	for {
		if <-proceed_chan == 4 {
			break
		}
	}
	p.Round = 5
	time.Sleep(time.Second * 1)

	Round7(send_chan, p, receive_chan, &Data.Keygen_All_Data)

	Values.Keygen.U_i = Data.Keygen_All_Data.U_i

	for i := 0; i < len(p.Sorted_Peers); i++ {

		if i == p.My_Index {
			continue
		}
		send_chan <- ths.Message{From: *p,
			Type:         5,
			To:           p.Sorted_Peers[i].Id,
			Payload_name: "Fifth",
			Payload:      Values,
			Status:       0}

	}
	for {
		if len(p.Round3) == len(p.Peers) {
			Ack_sender <- 5
			break
		}
		time.Sleep(time.Second * 2)
	}
	for {
		if <-proceed_chan == 5 {
			break
		}
	}

	Message := "RANDOM_MESSAGE"

	var choice int
	fmt.Println("Do do Sign the Message:", Message, "\n 1)Yes 2)No ")
	fmt.Scanln(&choice)
	time.Sleep(time.Second * 5)

	if choice == 1 {
		p.Round = 6

		Signing(send_chan, p, receive_chan, &Data.Keygen_All_Data, Message)
		Values.Keygen.V_i = Data.Keygen_All_Data.V_i
		for i := 0; i < len(p.Sorted_Peers); i++ {

			if i == p.My_Index {
				continue
			}
			send_chan <- ths.Message{From: *p,
				Type:         6,
				To:           p.Sorted_Peers[i].Id,
				Payload_name: "Sixth",
				Payload:      Values,
				Status:       0}

		}
		for {
			if len(p.Round3) == len(p.Peers) {
				Ack_sender <- 6
				break
			}
			time.Sleep(time.Second * 2)
		}
		for {
			if <-proceed_chan == 6 {
				break
			}
		}
	}
	time.Sleep(time.Second * 1)

	Combine(send_chan, p, receive_chan, &Data.Keygen_All_Data, Message)

}
