package ths

import (
	"context"

	host "github.com/libp2p/go-libp2p/core/host"
	peer "github.com/libp2p/go-libp2p/core/peer"
	"gopkg.in/dedis/kyber.v2"
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
	Payload      Round_Data
	Status       int
}

type Data struct {
	Keygen_All_Data Keygen_Store
	Sign_All_Data   string
}

//Contains Data to be broadcasted
type Round_Data struct {
	Keygen Keygen_Data
}
type Payload struct {
	Sender peer.ID
	// Sender       string
	Type         int
	Payload_name string     //"C1,C2,C3"
	Payload      Round_Data //"drhdrhdrh,hdhdth,shsdthsdth"
}

type P2P struct {

	// Represents the libp2p host
	Host             host.Host
	Host_ip          string
	Ctx              context.Context
	Peers            []THS
	Connectedparties int
	ThisParty        int

	//Required for Peer finding API
	Host_ips []string
	Host_id  string //peer.ID to []byte

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

	Round3 []Keygen_Store_Round3

	Round4 []Keygen_Store_Round4

	Round5 []Keygen_Store_Round5
}

type Moniker_message struct {
	Moniker string
}

//The Data to be Broadcasted
type Keygen_Data struct {
	EPK         string //curves.Point
	SPK         string //kyber.Point
	KGC         KGC
	Alphas      []string
	Enc_shares  []Encrypted_Share
	KGC_sign    KGC
	Alphas_sign []string
	Shares_sign []string
	U_i         string
}

//All the key values generated will be stored in this structure
type Keygen_Store struct {
	EPK    string //curves.Point
	ESK    string //curves.Scalar
	SSK    string //kyber.Scalar
	SPK    string //kyber.Point
	KGC    KGC
	Alphas []string
	Poly   []string

	Shares           []string
	Encrypted_Shares []Encrypted_Share
	V2               string
	V3               string
	R_i              string //Secret of Presign Poly
	KGC_sign         KGC
	Poly_sign        []string
	Alphas_sign      []string
	Shares_sign      []string

	R   string //Presing R_i
	U_i string
	U   string
}

type Encrypted_Share struct {
	C1 string
	C2 []byte
	C3 []byte
}

type KGC struct {
	Signature_S string //kyber.Scalar
	Public_key  string //kyber.Point
	Message     string
	KGD         string //kyber.Point
}

type Keygen_Store_Round1 struct {
	Id  peer.ID
	V1  Round_Data
	Ack int
}
type Keygen_Store_Round2 struct {
	Id  peer.ID
	V1  Round_Data
	Ack int
}
type Keygen_Store_Round3 struct {
	Id  peer.ID
	V1  Round_Data
	Ack int
}

type Keygen_Store_Round4 struct {
	Id  peer.ID
	V1  Round_Data
	Ack int
}

type Keygen_Store_Round5 struct {
	Id  peer.ID
	V1  Round_Data
	Ack int
}

// var Participants []THS
