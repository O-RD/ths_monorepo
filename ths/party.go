package ths

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/O-RD/ths_monorepo/p2p"
	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	connmgr "github.com/libp2p/go-libp2p/p2p/net/connmgr"

	"github.com/multiformats/go-multiaddr"
)

// func PeerInSlice(a peer.ID, list []peer.ID) bool {
// 	fmt.Println(ths.THSType)
// 	for _, b := range list {
// 		if b == a {
// 			return true
// 		}
// 	}
// 	return false
// }

func create_host() (host.Host, error) {

	// Creates a new RSA key pair for this host.
	// Read RSA keys from file
	prvKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048) //, randomness)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// 0.0.0.0 will listen on any interface device.
	var port = rand.Intn(9999)

	//
	//50000-
	// log.Println("Node Port- ", 0)
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	conn_mgr, err := connmgr.NewConnManager(100, 400)
	if err != nil {
		log.Println(err, "Error in Creating conn manager")

	}
	conn := libp2p.ConnectionManager(conn_mgr)
	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	opts := libp2p.ChainOptions(libp2p.ListenAddrs(sourceMultiAddr), libp2p.Identity(prvKey), conn)

	return libp2p.New(opts)
	//return libp2p.New()
}

func P2p_init(p2p_chan chan P2P) {

	//select {}
	//Open Send channel

	//1. Setup Host
	var h, _ = create_host()
	// log.Println(h.Addrs()[0].String() + "/p2p/" + h.ID().String())
	//2.
	ctx := context.Background()
	var p P2P
	p.Host = h
	p.Host_ip = h.Addrs()[0].String() + "/p2p/" + h.ID().String()
	p.Ctx = ctx
	// Create_peer(p)
	p2p_chan <- p
	close(p2p_chan)

}

var peer_details_list []string

func Create_peer(p *P2P) {

	//a := get_list(h, *channel_id, ctx)

	//Setup listener
	p2p.Connection_Stream_listener(p)

	peerChan := initMDNS(p.Host, p.Port)
	time.Sleep(time.Second * 5)

	for external_peer := range peerChan {
		//log.Println(var_counter, *num_users)
		if external_peer.ID == p.Host.ID() {
			continue
		}

		// log.Println("Found peer:", external_peer, ", connecting")
		if err := p.Host.Connect(p.Ctx, external_peer); err != nil {
			log.Println("Connection failed:", external_peer.ID)
		} else {

			send_stream, _ := p.Host.NewStream(p.Ctx, external_peer.ID, "ths_stream")
			message := Moniker_message{
				Name: p.Moniker,
			}
			b_message, _ := json.Marshal(message)
			_, err = send_stream.Write(append(b_message, '\n'))
			p.Connectedparties += 1
			break_flag := 0
			for {
				if break_flag == 1 {
					break
				}
				for i := 0; i < len(p.Peers); i++ {
					if p.Peers[i].Id == external_peer.ID {
						log.Println("Connected to ", external_peer.ID, " with Moniker: ", p.Peers[i].Moniker)

						break_flag = 1
					}
					time.Sleep(time.Second)
				}
			}

		}
		if len(p.Peers) >= p.Party_Size-1 && p.Connectedparties >= p.Party_Size-1 {
			time.Sleep(time.Second * 2)
			break
		}
	}
}
