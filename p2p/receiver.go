package p2p

import (
	"bufio"
	"encoding/json"

	ths "github.com/O-RD/ths_monorepo/ths"
	"github.com/libp2p/go-libp2p/core/network"
)

func Connection_Stream_listener(p *ths.P2P) {
	//fmt.Println("Got a new stream!")

	p.Host.SetStreamHandler("ths_stream", func(s network.Stream) {
		// log.Println("sender received new stream")

		buf := bufio.NewReader(s)
		//log.Println(s)
		str, _ := buf.ReadBytes('\n')

		bytes := []byte(str)
		var message_receive ths.Moniker_message
		json.Unmarshal(bytes, &message_receive)
		l.Lock()
		p.Peers = append(p.Peers, ths.THS{Id: s.Conn().RemotePeer(),
			Moniker: message_receive.Moniker})
		l.Unlock()

		s.Close()
	})

}

func Input_Stream_listener(p *ths.P2P, receiver_ch chan ths.Payload) {
	//fmt.Println("Got a new stream!")

	p.Host.SetStreamHandler("ths_stream", func(s network.Stream) {
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
