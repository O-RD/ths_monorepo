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
			added_peer <- ths.THS{Id: s.Conn().RemotePeer(),
				Moniker: message_receive.Moniker,
				Round:   0,
			}
			s.Close()

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
				fmt.Println("Sent to", s.Conn().RemotePeer())
			}
		}
	})

}

func contains(peers []ths.THS, peer_id peer.ID) bool {
	for _, v := range peers {
		if v.Id == peer_id {
			return true
		}
	}

	return false
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

		receiver_ch <- ths.Payload{Sender: s.Conn().RemotePeer(), Payload: message_receive.Payload, Payload_name: message_receive.Payload_name, Type: message_receive.Type}

	})

}
