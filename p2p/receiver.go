package p2p

import (
	"bufio"
	"encoding/json"

	"github.com/libp2p/go-libp2p-core/network"
)

func connection_Stream_listener(p *P2P) {
	//fmt.Println("Got a new stream!")

	p.Host.SetStreamHandler("moniker", func(s network.Stream) {
		// log.Println("sender received new stream")

		buf := bufio.NewReader(s)
		//log.Println(s)
		str, _ := buf.ReadBytes('\n')

		bytes := []byte(str)
		var message_receive moniker_message
		json.Unmarshal(bytes, &message_receive)

		l.Lock()
		p.Peers = append(p.Peers, peer_names{Id: s.Conn().RemotePeer(),
			Name: message_receive.Name})
		l.Unlock()

		s.Close()

	})

}
