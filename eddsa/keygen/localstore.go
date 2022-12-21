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
		message_receive := <-receive_chan
		if message_receive.Type == 1 {
			p.Round1 = append(p.Round1, ths.Keygen_Store_Round1{Id: message_receive.Sender,
				V1:  message_receive.Payload,
				Ack: 0,
			})
			sender_index := p2p.Get_index(p.Sorted_Peers, message_receive.Sender)
			if sender_index <= 0 {
				fmt.Println("User Not Found")
				return
			}
			os.MkdirAll("Received/"+strconv.Itoa(sender_index)+"/Keys", 0777)

			path := "Received/" + strconv.Itoa(sender_index) + "/Keys"
			f, _ := os.Create(path + "/EPK.txt")
			f.WriteString(message_receive.Payload.Keygen.EPK)
			f, _ = os.Create(path + "/SPK.txt")
			f.WriteString(message_receive.Payload.Keygen.SPK)

			os.MkdirAll("Received/"+strconv.Itoa(sender_index)+"/Keygen_commit", 0777)
			path = "Received/" + strconv.Itoa(sender_index) + "/Keygen_commit"
			f, _ = os.Create(path + "/Signature_S.txt")
			f.WriteString(message_receive.Payload.Keygen.KGC.Signature_S)
			f, _ = os.Create(path + "/Pubkey.txt")
			f.WriteString(message_receive.Payload.Keygen.KGC.Public_key)
			f, _ = os.Create(path + "/Message.txt")
			f.WriteString(message_receive.Payload.Keygen.KGC.Message)
			f, _ = os.Create(path + "/KGD.txt")
			f.WriteString(message_receive.Payload.Keygen.KGC.KGD)

		} else if message_receive.Type == 2 {
			p.Round2 = append(p.Round2, ths.Keygen_Store_Round2{Id: message_receive.Sender,
				V1:  message_receive.Payload,
				Ack: 0,
			})
			sender_index := p2p.Get_index(p.Sorted_Peers, message_receive.Sender)

			if sender_index == p.My_Index+1 {
				continue
			} else {
				os.MkdirAll("Received/"+strconv.Itoa(sender_index)+"/Keygen_alphas", 0777)
				var i int
				path := "Received/" + strconv.Itoa(sender_index) + "/Keygen_alphas"
				for i = 0; i < p.Threshold; i++ {
					f, _ := os.Create(path + "/alpha" + strconv.Itoa(i) + ".txt")
					f.WriteString(message_receive.Payload.Keygen.Alphas[i])
				}
				os.MkdirAll("Received/"+strconv.Itoa(sender_index)+"/Keygen_shares", 0777)

				f, _ := os.Create("Received/" + strconv.Itoa(sender_index) + "/Keygen_shares/C1.txt")
				f.WriteString(message_receive.Payload.Keygen.Enc_shares[p.My_Index].C1)
				f, _ = os.Create("Received/" + strconv.Itoa(sender_index) + "/Keygen_shares/C2.txt")
				f.WriteString(message_receive.Payload.Keygen.Enc_shares[p.My_Index].C2)
				f, _ = os.Create("Received/" + strconv.Itoa(sender_index) + "/Keygen_shares/C3.txt")
				f.WriteString(message_receive.Payload.Keygen.Enc_shares[p.My_Index].C3)
			}
			// if len(p.Round2) == len(p.Peers) {
			// 	Ack_sender <- 2
			// }

		} else if message_receive.Type == 3 {
			p.Round3 = append(p.Round3, ths.Keygen_Store_Round3{Id: message_receive.Sender,
				V1:  message_receive.Payload,
				Ack: 0,
			})
			sender_index := p2p.Get_index(p.Sorted_Peers, message_receive.Sender)

			os.MkdirAll("Received/"+strconv.Itoa(sender_index)+"/Presigning_commit", 0777)
			path := "Received/" + strconv.Itoa(sender_index) + "/Presigning_commit"
			f, _ := os.Create(path + "/Signature_S.txt")
			f.WriteString(message_receive.Payload.Keygen.KGC_sign.Signature_S)
			f, _ = os.Create(path + "/Pubkey.txt")
			f.WriteString(message_receive.Payload.Keygen.KGC_sign.Public_key)
			f, _ = os.Create(path + "/Message.txt")
			f.WriteString(message_receive.Payload.Keygen.KGC_sign.Message)
			f, _ = os.Create(path + "/KGD.txt")
			f.WriteString(message_receive.Payload.Keygen.KGC_sign.KGD)
		} else if message_receive.Type == 4 {
			p.Round4 = append(p.Round4, ths.Keygen_Store_Round4{Id: message_receive.Sender,
				V1:  message_receive.Payload,
				Ack: 0,
			})
			sender_index := p2p.Get_index(p.Sorted_Peers, message_receive.Sender)
			os.MkdirAll("Received/"+strconv.Itoa(sender_index)+"/Presigning_alphas", 0777)

			path := "Received/" + strconv.Itoa(sender_index) + "/Presigning_alphas"
			var i int
			for i = 0; i < p.Threshold; i++ {
				f, _ := os.Create(path + "/alpha" + strconv.Itoa(i) + ".txt")
				f.WriteString(message_receive.Payload.Keygen.Alphas_sign[i])
			}
			os.MkdirAll("Data/Presigning_shares", 0777)
			f, _ := os.Create("Data/Presigning_shares/share" + strconv.Itoa(sender_index) + ".txt")
			f.WriteString(message_receive.Payload.Keygen.Shares_sign[p.My_Index])

		} else if message_receive.Type == 5 {
			p.Round5 = append(p.Round5, ths.Keygen_Store_Round5{Id: message_receive.Sender,
				V1:  message_receive.Payload,
				Ack: 0,
			})
			sender_index := p2p.Get_index(p.Sorted_Peers, message_receive.Sender)
			os.MkdirAll("Received/Signing/U_i", 0777)
			f, _ := os.Create("Received/Signing/U_i/U_" + strconv.Itoa(sender_index) + ".txt")
			f.WriteString(message_receive.Payload.Keygen.U_i)

		}
		fmt.Println(message_receive)

	}
}
