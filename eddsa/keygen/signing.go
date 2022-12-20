package keygen

import (
	"fmt"
	"os"

	"github.com/O-RD/ths_monorepo/ths"
	"gopkg.in/dedis/kyber.v2/util/encoding"
)

func Signing(send_chan chan ths.Message, p *ths.P2P, receive_chan chan ths.Payload, message string) {
	peer_number := fmt.Sprint(p.My_Index + 1)
	Peer_Count := len(p.Sorted_Peers)
	fmt.Println("PEERCOUNT:", Peer_Count)
	fmt.Println("PEERNUMBER:", peer_number)
	var T int64 = int64(p.Threshold)
	fmt.Println("THRESHOLD:", T)

	fmt.Println("MESSAGE TO SIGN:", message)

	// T_arr = append(T_arr, my_index+1)
	// r_i:= curve.Scalar().Pick(curve.RandomStream())

	fmt.Printf("********************************************* SIGNING PHASES STARTED ******************************************\n")

	// file, _ := os.Open("Data/Signing/R_i.txt")

	file, _ := os.Open("Data/G.txt")
	x_i, _ := encoding.ReadHexScalar(curve, file)

	file, _ = os.Open("Data/Signing/U.txt")
	U, _ := encoding.ReadHexPoint(curve, file)
	fmt.Println("U from PreSign:", U.String())

	V_i, U_i := Signing_T_Unkown(U, x_i, message, peer_number)
	fmt.Println("U_i returned from sign:", U_i.String(), "\n")

	fmt.Println("SIGNATURE GENERATED:", V_i.String())

	f, _ := os.Open("Data/Final_sign.txt")
	encoding.WriteHexScalar(curve, f, V_i)

}
