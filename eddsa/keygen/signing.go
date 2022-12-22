package keygen

import (
	"fmt"
	"os"

	"github.com/O-RD/ths_monorepo/ths"
	"gopkg.in/dedis/kyber.v2/util/encoding"
)

func Signing(send_chan chan ths.Message, p *ths.P2P, receive_chan chan ths.Payload, Sign_Values *ths.Keygen_Store, message string) {
	peer_number := fmt.Sprint(p.My_Index + 1)
	Peer_Count := len(p.Sorted_Peers)
	fmt.Println("PEERCOUNT:", Peer_Count)
	fmt.Println("PEERNUMBER:", peer_number)
	var T int64 = int64(p.Threshold)
	fmt.Println("THRESHOLD:", T)

	fmt.Println("MESSAGE TO SIGN:", message)

	T_arr = append(T_arr, p.My_Index+1)
	// r_i:= curve.Scalar().Pick(curve.RandomStream())

	fmt.Printf("********************************************* SIGNING PHASES STARTED ******************************************\n")

	// file, _ := os.Open("Data/Signing/R_i.txt")

	file, _ := os.Open("Data/" + peer_number + "/G.txt")
	x_i, _ := encoding.ReadHexScalar(curve, file)

	file, _ = os.Open("Data/" + peer_number + "/Signing/U.txt")
	U, _ := encoding.ReadHexPoint(curve, file)
	fmt.Println("U from PreSign:", U.String())

	V_i, U_i := Signing_T_Unkown(U, x_i, message, peer_number)
	fmt.Println("U_i returned from sign:", U_i.String(), "\n")

	fmt.Println("SIGNATURE GENERATED:", V_i.String())

	f, _ := os.Open("Data/Final_sign.txt")
	encoding.WriteHexScalar(curve, f, V_i)

	X_i := curve.Point().Mul(x_i, g)
	if Verify_sign_share(V_i, U, U_i, message, X_i) {
		fmt.Println("INDIVIDUAL SHARES ARE VERIFIED")
	} else {
		fmt.Println("NOT VERIFIED INDIVIDUAL SHARES")
	}

	Sign_Values.V_i, _ = encoding.ScalarToStringHex(curve, V_i)
	file, _ = os.Create("Data/" + peer_number + "/Signing/V_i.txt")
	encoding.WriteHexScalar(curve, file, V_i)

}

func Combine(send_chan chan ths.Message, p *ths.P2P, receive_chan chan ths.Payload, Sign_Values *ths.Keygen_Store, Message string) {
	peer_number := fmt.Sprint(p.My_Index + 1)
	Peer_Count := len(p.Sorted_Peers)
	fmt.Println("PEERCOUNT:", Peer_Count)
	fmt.Println("PEERNUMBER:", peer_number)
	var T int = p.Threshold
	fmt.Println("THRESHOLD:", T)

	fmt.Println("************ COMBINATION PHASE ****************")
	fmt.Println(Peer_Count)
	GKey := Get_Group_Key(int64(Peer_Count), p.My_Index)
	Vsum, Usum := combine_T_Unknown(p.My_Index, T)
	fmt.Println("************ VERIFYING ****************")

	fmt.Println("GKEY:", GKey.String())
	fmt.Println("VSUM:", Vsum.String())
	fmt.Println("Usum:", Usum.String())

	res := Verify_sign_share(Vsum, Usum, Usum, Message, GKey)
	if res {
		fmt.Println("SUCCESS VERIFICATION OF SIGNATURE")
	} else {
		fmt.Println("FAILED TO VERIFIY")
	}
	fmt.Println("FINAL SIGNATURE:", Vsum.String())

}
