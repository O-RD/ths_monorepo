package ths

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"

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
	var port = rand.Intn(1000)
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

	peerChan := initMDNS(p.Host, p.Topic)
	// time.Sleep(time.Second * 5)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for external_peer := range peerChan {

			if external_peer.ID == p.Host.ID() {
				continue
			}

			log.Println("Found peer:", external_peer, ", connecting")
			if err := p.Host.Connect(p.Ctx, external_peer); err != nil {
				log.Println("Connection failed:", external_peer.ID)
			} else {

				send_stream, _ := p.Host.NewStream(p.Ctx, external_peer.ID, "ths_stream")
				message := Moniker_message{
					Moniker: p.Moniker,
				}
				b_message, _ := json.Marshal(message)
				_, err = send_stream.Write(append(b_message, '\n'))
				if err == nil {
					p.Connectedparties += 1
				}
				fmt.Println("Sent to", external_peer.ID)

			}
		}
	}()
	wg.Wait()
}

func Sort_Peers(party *P2P) {

	list_of_external := party.Peers
	my_id := THS{
		Id:      party.Host.ID(),
		Moniker: party.Moniker,
		Round:   0,
	}
	list_of_external = append(list_of_external, my_id)

	//Bubble Sort peers
	for i := 0; i < len(list_of_external)-1; i++ {
		for j := 0; j < len(list_of_external)-i-1; j++ {
			if list_of_external[j].Id.String() > list_of_external[j+1].Id.String() {
				list_of_external[j], list_of_external[j+1] = list_of_external[j+1], list_of_external[j]
			}
		}
	}
	for i := 0; i < len(list_of_external); i++ {
		if list_of_external[i].Id == party.Host.ID() {
			party.My_Index = i
		}
	}
	party.Sorted_Peers = list_of_external

}
