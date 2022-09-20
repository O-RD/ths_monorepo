package p2p

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	connmgr "github.com/libp2p/go-libp2p/p2p/net/connmgr"

	"github.com/multiformats/go-multiaddr"
)

func PeerInSlice(a peer.ID, list []peer.ID) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

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

// func create_peer(h host.Host, ctx context.Context) {

// 	//a := get_list(h, *channel_id, ctx)
// 	fmt.Println(status_struct.Chan)
// 	peerChan := initMDNS(h, status_struct.Chan)
// 	time.Sleep(time.Second * 5)

// 	//added now
// 	//var a string
// 	//fmt.Scanln(&a)

// 	var var_counter int = 1
// 	for peer := range peerChan {
// 		//log.Println(var_counter, *num_users)
// 		if peer.ID == h.ID() {
// 			continue
// 		}

// 		log.Println("Found peer:", peer, ", connecting")
// 		if err := h.Connect(ctx, peer); err != nil {
// 			log.Println("Connection failed:", err)
// 		} else {
// 			var_counter += 1
// 			//temp := peer_details{peer.ID, peer}
// 			//function_peer_chan <- temp
// 			//log.Println("Connected to ", peer.Addrs[0])
// 			//fmt.Println(peer.ID)
// 			peer_details_list = append(peer_details_list, peer_details{peer.ID, peer})
// 			//log.Println(peer_details_list[0].id)
// 			if var_counter == status_struct.Num_peers {
// 				//log.Println("get out")
// 				time.Sleep(time.Second * 6)
// 				break
// 			}
// 		}
// 	}
// }

func start_p2p() *P2P {

	//select {}

	//1. Setup Host
	var h, _ = create_host()
	log.Println(h.Addrs()[0].String() + "/p2p/" + h.ID().String())
	//2.
	ctx := context.Background()
	// create_peer(h, ctx)

	// for i, item := range peer_details_list {
	// 	log.Println(i, item.addr.Addrs[0])
	// }

	return &P2P{
		Host:    h,
		Host_ip: h.Addrs()[0].String() + "/p2p/" + h.ID().String(),
		Ctx:     ctx,
		Peers:   []string{},
	}
	//var choice int
	//fmt.Println("Enter choice")

	//choice =1
}
