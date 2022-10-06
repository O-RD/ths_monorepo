package ths

import (
	"context"

	host "github.com/libp2p/go-libp2p/core/host"
	peer "github.com/libp2p/go-libp2p/core/peer"
	"gopkg.in/dedis/kyber.v2"

	"github.com/coinbase/kryptology/pkg/core/curves"
	"gopkg.in/dedis/kyber.v2/group/edwards25519"
)

type THS struct {
	Id      peer.ID
	Moniker string
	Round   int
}

type Signature struct {
	R kyber.Point
	S kyber.Scalar
}

type Message struct {
	From         P2P
	To           peer.ID
	Type         int
	Payload_name string
	Payload      Data
	Status       int
}

type Data struct {
	Keygen_Data Keygen_Store
	Sign_Data   string
}

type Payload struct {
	Sender peer.ID
	// Sender       string
	Type         int
	Payload_name string //"C1,C2,C3"
	Payload      Data   //"drhdrhdrh,hdhdth,shsdthsdth"
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

type Keygen_Store struct {
	Curve            *edwards25519.SuiteEd25519
	EPK              curves.Point
	ESK              curves.Scalar
	SSK              kyber.Scalar
	SPK              kyber.Point
	KGC              KGC
	Alphas           []kyber.Point
	Encrypted_Shares []Encrypted_Share
	V2               string
	V3               string
}

type Encrypted_Share struct {
	C1 curves.Point
	C2 []byte
	C3 []byte
}

type KGC struct {
	Sign       kyber.Scalar
	Public_key kyber.Point
	Message    string
	KGD        kyber.Point
}

type Keygen_Store_Round1 struct {
	Id  peer.ID
	V1  Data
	Ack int
}
type Keygen_Store_Round2 struct {
	Id  peer.ID
	V1  Data
	Ack int
}

// var Participants []THS
