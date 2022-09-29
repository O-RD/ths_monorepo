package p2p

import ths "github.com/O-RD/ths_monorepo/ths"

func Create_Peer(p *ths.P2P) {
	connection_Stream_listener(p)
	ths.Create_peer(p)
}
