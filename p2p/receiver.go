package p2p

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"

	ths "github.com/O-RD/ths_monorepo/ths"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

func Connection_Stream_listener(p *ths.P2P, added_peer chan ths.THS) {
	//fmt.Println("Got a new stream!")

	p.Host.SetStreamHandler("ths_stream", func(s network.Stream) {
		// log.Println("sender received new stream")

		buf := bufio.NewReader(s)
		str, _ := buf.ReadBytes('\n')

		bytes := []byte(str)
		var message_receive ths.Moniker_message
		json.Unmarshal(bytes, &message_receive)
		if contains(p.Peers, s.Conn().RemotePeer()) == false {
			l.Lock()
			p.Peers = append(p.Peers, ths.THS{Id: s.Conn().RemotePeer(),
				Moniker: message_receive.Moniker,
				Round:   0,
			})
			l.Unlock()
			s.Close()
			added_peer <- ths.THS{Id: s.Conn().RemotePeer(),
				Moniker: message_receive.Moniker,
				Round:   0,
			}

			//Send Conn request if not sent prior
			s.Conn().RemoteMultiaddr()
			x, _ := peer.AddrInfoFromString(s.Conn().RemoteMultiaddr().String() + "/p2p/" + s.Conn().RemotePeer().String())
			if err := p.Host.Connect(p.Ctx, *x); err != nil {
				log.Println("Connection failed:", s.Conn().RemotePeer())
			} else {
				send_stream, _ := p.Host.NewStream(p.Ctx, s.Conn().RemotePeer(), "ths_stream")
				message := ths.Moniker_message{
					Moniker: p.Moniker,
				}
				b_message, _ := json.Marshal(message)
				_, err := send_stream.Write(append(b_message, '\n'))
				if err == nil {
					p.Connectedparties += 1
				}
				// fmt.Println("Sent to", s.Conn().RemotePeer())
			}
		} else {
			added_peer <- ths.THS{}
		}
	})

}

func Input_Stream_listener(p *ths.P2P, receiver_ch chan ths.Payload) {
	//fmt.Println("Got a new stream!")

	p.Host.SetStreamHandler("ths_stream_keygen", func(s network.Stream) {
		// log.Println("sender received new stream")

		buf := bufio.NewReader(s)
		//log.Println(s)
		str, _ := buf.ReadBytes('\n')

		bytes := []byte(str)
		var message_receive ths.Payload
		json.Unmarshal(bytes, &message_receive)
		if message_receive.Type == 1 {

			if containsR1(p.Round1, s.Conn().RemotePeer()) == false {
				receiver_ch <- ths.Payload{Sender: s.Conn().RemotePeer(), Payload: message_receive.Payload, Payload_name: message_receive.Payload_name, Type: message_receive.Type}

			}
		} else if message_receive.Type == 2 {
			if containsR2(p.Round2, s.Conn().RemotePeer()) == false {
				receiver_ch <- ths.Payload{Sender: s.Conn().RemotePeer(), Payload: message_receive.Payload, Payload_name: message_receive.Payload_name, Type: message_receive.Type}
			}
		}

	})

}

func Acknowledgement_listener(p *ths.P2P, proceed chan int) {
	p.Host.SetStreamHandler("ths_stream_ack", func(s network.Stream) {
		// log.Println("sender received new stream")

		buf := bufio.NewReader(s)
		//log.Println(s)
		str, _ := buf.ReadBytes('\n')

		bytes := []byte(str)
		var message_receive int
		json.Unmarshal(bytes, &message_receive)
		fmt.Println("Received", message_receive, "from", s.Conn().RemotePeer)
		if message_receive == 1 {

			for i := 0; i < len(p.Round1); i++ {
				if p.Round1[i].Id == s.Conn().RemotePeer() {
					p.Round1[i].Ack = 1
				}
			}
			if len(p.Round1) == len(p.Peers) && AckR1(p.Round1) {
				// p.Round1[0].Ack = 0
				fmt.Println("End")
				proceed <- 1
			}

		} else if message_receive == 2 {
			for i := 0; i < len(p.Round2); i++ {
				if p.Round2[i].Id == s.Conn().RemotePeer() {
					p.Round2[i].Ack = 2
				}
			}
			if len(p.Round2) == len(p.Peers) && AckR2(p.Round2) {
				// p.Round2[0].Ack = 0

				proceed <- 2
			}
		}
	})
}
