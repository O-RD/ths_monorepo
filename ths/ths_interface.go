package ths

import (
	"context"

	host "github.com/libp2p/go-libp2p/core/host"
	peer "github.com/libp2p/go-libp2p/core/peer"
)

type THS struct {
	Id      peer.ID
	Moniker string
	Round   int
}

type Message struct {
	From         P2P
	To           peer.ID
	Type         int
	Payload_name string
	Payload      string
	End          int
}
type Payload struct {
	Sender       peer.ID
	Type         int
	Payload_name string //"C1,C2,C3"
	Payload      string //"drhdrhdrh,hdhdth,shsdthsdth"
}
type P2P struct {

	// Represents the libp2p host
	Host             host.Host
	Host_ip          string
	Ctx              context.Context
	Peers            []THS
	Connectedparties int
	ThisParty        int

	//Used for indexing peers
	Sorted_Peers []THS
	My_Index     int

	//Used for comms
	send       chan Message
	receive    chan Message
	Threshold  int
	Party_Size int
	Moniker    string
	Port       string
	ThsType    string

	Round int
}

type Moniker_message struct {
	Moniker string
}

type Keygen_Store struct {
	V1 string
	V2 string
	V3 string
}

// var Participants []THS
