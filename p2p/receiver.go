package p2p

import (
	"bufio"
	"encoding/json"

	"github.com/O-RD/ths_monorepo/ths"
	"github.com/libp2p/go-libp2p/core/network"
)

func connection_Stream_listener(p *ths.P2P) {
	//fmt.Println("Got a new stream!")

	p.Host.SetStreamHandler("moniker", func(s network.Stream) {
		// log.Println("sender received new stream")

		buf := bufio.NewReader(s)
		//log.Println(s)
		str, _ := buf.ReadBytes('\n')

		bytes := []byte(str)
		var message_receive ths.Moniker_message
		json.Unmarshal(bytes, &message_receive)

		l.Lock()
		p.Peers = append(p.Peers, ths.THS{Id: s.Conn().RemotePeer(),
			Moniker: message_receive.Name})
		l.Unlock()

		s.Close()

	})

}
