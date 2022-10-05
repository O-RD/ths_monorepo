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
	Status       int
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
	Topic      string
	ThsType    string

	Round int

	Round1 []Keygen_Store_Round1

	Round2 []Keygen_Store_Round2
}

type Moniker_message struct {
	Moniker string
}

type Keygen_Store_Round1 struct {
	Id  peer.ID
	V1  string
	Ack int
}
type Keygen_Store_Round2 struct {
	Id  peer.ID
	V1  string
	Ack int
}

// var Participants []THS
