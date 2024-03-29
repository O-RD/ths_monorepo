//Using Schnorr to commit and decommit
//Make changes to the commit and decommit function according to your need i.e. take file names as arguments for the function

package keygen

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/O-RD/ths_monorepo/ths"
	"gopkg.in/dedis/kyber.v2"

	//"gopkg.in/dedis/kyber.v2/group/edwards25519"
	"gopkg.in/dedis/kyber.v2/util/encoding"
)

type Data struct {
	s       string
	pub     string
	message string
}

// m: Message
// S: Signature
func Comit_PublicKey(m string, S ths.Signature) kyber.Point {
	// Create a generator.
	g := curve.Point().Base()

	// e = Hash(m || r)
	e := Hash(m + S.R.String())

	// y = (r - s * G) * (1 / e)
	y := curve.Point().Sub(S.R, curve.Point().Mul(S.S, g))
	y = curve.Point().Mul(curve.Scalar().Div(curve.Scalar().One(), e), y)

	return y
}

// m: Message
// s: Signature
// y: Public key
func Comit_Verify(m string, S ths.Signature, y kyber.Point) bool {
	// Create a generator.
	g := curve.Point().Base()

	// e = Hash(m || r)
	e := Hash(m + S.R.String())

	// Attempt to reconstruct 's * G' with a provided signature; s * G = r - e * y
	sGv := curve.Point().Sub(S.R, curve.Point().Mul(e, y))

	// Construct the actual 's * G'
	sG := curve.Point().Mul(S.S, g)

	// Equality check; ensure signature and public key outputs to s * G.
	return sG.Equal(sGv)
}

func Commitment(x kyber.Scalar, m string, peer_number string, value_struct *ths.Keygen_Store) {
	path1 := "Temp/Commitment/" + peer_number + "/KGC"
	err := os.MkdirAll(path1, os.ModePerm)
	if err != nil {
		panic(err)
	}
	publicKey := curve.Point().Mul(x, curve.Point().Base())
	sig := Sign(m, x)

	f1, e1 := os.OpenFile(path1+"/Signature_S.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if e1 != nil {
		fmt.Println(e1)
	}
	encoding.WriteHexScalar(curve, f1, sig.S)

	f2, e2 := os.OpenFile(path1+"/PubKey.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if e2 != nil {
		fmt.Println(e2)
	}
	encoding.WriteHexPoint(curve, f2, publicKey)

	f3, e3 := os.OpenFile(path1+"/Message.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if e3 != nil {
		fmt.Println(e3)
	}
	f3.WriteString(m)
	f3.Close()
	f4, e4 := os.OpenFile("Temp/Commitment/"+peer_number+"/KGD.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if e4 != nil {
		fmt.Println(e4)
	}
	encoding.WriteHexPoint(curve, f4, sig.R)

	value_struct.KGC.Signature_S, _ = encoding.ScalarToStringHex(curve, sig.S)
	value_struct.KGC.Message = m
	value_struct.KGC.Public_key, _ = encoding.PointToStringHex(curve, publicKey)
	value_struct.KGC.KGD, _ = encoding.PointToStringHex(curve, sig.R)

	fmt.Printf("Commitment Done for Peer %s \n", peer_number)
}
func Commitment_sign(x kyber.Scalar, m string, peer_number string, value_struct *ths.Keygen_Store) {
	path1 := "Temp/Commitment/Signing/" + peer_number + "/KGC"
	err := os.MkdirAll(path1, os.ModePerm)
	if err != nil {
		panic(err)
	}
	publicKey := curve.Point().Mul(x, curve.Point().Base())
	sig := Sign(m, x)

	f1, e1 := os.OpenFile(path1+"/Signature_S.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if e1 != nil {
		fmt.Println(e1)
	}
	encoding.WriteHexScalar(curve, f1, sig.S)

	f2, e2 := os.OpenFile(path1+"/PubKey.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if e2 != nil {
		fmt.Println(e2)
	}
	encoding.WriteHexPoint(curve, f2, publicKey)

	f3, e3 := os.OpenFile(path1+"/Message.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if e3 != nil {
		fmt.Println(e3)
	}
	f3.WriteString(m)
	f3.Close()
	f4, e4 := os.OpenFile("Temp/Commitment/Signing/"+peer_number+"/KGD.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if e4 != nil {
		fmt.Println(e4)
	}
	encoding.WriteHexPoint(curve, f4, sig.R)
	value_struct.KGC_sign.Signature_S, _ = encoding.ScalarToStringHex(curve, sig.S)
	value_struct.KGC_sign.Message = m
	value_struct.KGC_sign.Public_key, _ = encoding.PointToStringHex(curve, publicKey)
	value_struct.KGC_sign.KGD, _ = encoding.PointToStringHex(curve, sig.R)
	fmt.Printf("Sign Commitment Done for Peer %s \n", peer_number)
}

func Decommitment_j(peer_number string) string {

	path := "Received/" + peer_number + "/Keygen_commit/"
	f1, e1 := os.Open(path + "Signature_S.txt")
	if e1 != nil {
		fmt.Println(e1)
	}
	sig_d, e := encoding.ReadHexScalar(curve, f1)
	if e != nil {
		fmt.Println(e)
	}

	f2, e2 := os.Open(path + "Pubkey.txt")
	if e2 != nil {
		fmt.Println(e2)
	}
	pub_key, e_2 := encoding.ReadHexPoint(curve, f2)
	if e_2 != nil {
		fmt.Println(e_2)
	}
	f3, e3 := os.Open(path + "KGD.txt")
	if e3 != nil {
		fmt.Println(e3)
	}
	KGD_j, e_3 := encoding.ReadHexPoint(curve, f3)
	if e_3 != nil {
		fmt.Println(e_3)
	}

	message, e4 := ioutil.ReadFile(path + "Message.txt")
	if e4 != nil {
		fmt.Println(e4)
	}

	newS := ths.Signature{}
	newS.S = sig_d
	newS.R = KGD_j

	t := Comit_Verify(string(message), newS, pub_key)
	if t {

		return pub_key.String()
	} else {
		return "Invalid"
	}
}

func Decommitment_j_sign(peer_number string) string {
	path := "Received/" + peer_number + "/Presigning_commit"
	f1, e1 := os.Open(path + "/Signature_S.txt")
	if e1 != nil {
		fmt.Println(e1)
	}
	sig_d, e := encoding.ReadHexScalar(curve, f1)
	if e != nil {
		fmt.Println(e)
	}

	f2, e2 := os.Open(path + "/Pubkey.txt")
	if e2 != nil {
		fmt.Println(e2)
	}
	pub_key, e_2 := encoding.ReadHexPoint(curve, f2)
	if e_2 != nil {
		fmt.Println(e_2)
	}
	f3, e3 := os.Open(path + "/KGD.txt")
	if e3 != nil {
		fmt.Println(e3)
	}
	KGD_j, e_3 := encoding.ReadHexPoint(curve, f3)
	if e_3 != nil {
		fmt.Println(e_3)
	}

	message, e4 := ioutil.ReadFile(path + "/Message.txt")
	if e4 != nil {
		fmt.Println(e4)
	}

	newS := ths.Signature{}
	newS.S = sig_d
	newS.R = KGD_j

	t := Comit_Verify(string(message), newS, pub_key)
	if t {

		return pub_key.String()
	} else {
		return "Invalid"
	}
}
