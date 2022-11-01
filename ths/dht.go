package ths

import (
	"fmt"
	"strings"
	"sync"

	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"

	maddr "github.com/multiformats/go-multiaddr"
)

type addrList []maddr.Multiaddr

func (al *addrList) String() string {
	strs := make([]string, len(*al))
	for i, addr := range *al {
		strs[i] = addr.String()
	}
	return strings.Join(strs, ",")
}

func (al *addrList) Set(value string) error {
	addr, err := maddr.NewMultiaddr(value)
	if err != nil {
		return err
	}
	*al = append(*al, addr)
	return nil
}

var BootstrapPeers addrList

func initDHT(p *P2P) chan peer.AddrInfo {

	kademliaDHT, err := dht.New(p.Ctx, p.Host)
	if err != nil {
		panic(err)
	}

	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.

	if err = kademliaDHT.Bootstrap(p.Ctx); err != nil {
		panic(err)
	}

	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.
	var wg sync.WaitGroup
	for _, peerAddr := range BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := p.Host.Connect(p.Ctx, *peerinfo); err != nil {
				fmt.Print()
			}
		}()
	}
	wg.Wait()

	// We use a rendezvous point "meet me here" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.

	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(p.Ctx, routingDiscovery, p.Topic)

	// Now, look for others who have announced
	// This is like your friend telling you the location to meet you.

	// PeerChan = make(chan peer.AddrInfo)
	peerChan, err := routingDiscovery.FindPeers(p.Ctx, p.Topic)
	PeerChan := make(chan peer.AddrInfo)
	for item := range peerChan {
		PeerChan <- item
	}

	if err != nil {
		panic(err)
	}

	return PeerChan
}
