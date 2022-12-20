package p2p

import (
	"fmt"

	ths "github.com/O-RD/ths_monorepo/ths"
	"github.com/libp2p/go-libp2p/core/peer"
)

func contains(peers []ths.THS, peer_id peer.ID) bool {
	for _, v := range peers {
		if v.Id == peer_id {
			return true
		}
	}

	return false
}

func containsR1(peers []ths.Keygen_Store_Round1, peer_id peer.ID) bool {
	for _, v := range peers {
		if v.Id == peer_id {
			return true
		}
	}

	return false
}

func containsR2(peers []ths.Keygen_Store_Round2, peer_id peer.ID) bool {
	for _, v := range peers {
		if v.Id == peer_id {
			return true
		}
	}

	return false
}
func containsR3(peers []ths.Keygen_Store_Round3, peer_id peer.ID) bool {
	for _, v := range peers {
		if v.Id == peer_id {
			return true
		}
	}

	return false
}
func AckR1(peers []ths.Keygen_Store_Round1) bool {
	fmt.Println("First Ack-", peers)
	for _, v := range peers {

		if v.Ack == 0 {
			return false
		}
	}

	return true
}

func AckR2(peers []ths.Keygen_Store_Round2) bool {
	// fmt.Println("Second", peers)
	for _, v := range peers {
		if v.Ack == 0 {
			return false
		}
	}

	return true
}

func GetIndex(Sorted_Peers []ths.THS, peer_id peer.ID) int {
	for i, v := range Sorted_Peers {
		if v.Id == peer_id {
			return i
		}
	}

	return -1
}
