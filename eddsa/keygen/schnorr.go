package keygen

import (
	SHA_256 "crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/coinbase/kryptology/pkg/core/curves"
	"gopkg.in/dedis/kyber.v2"
	"gopkg.in/dedis/kyber.v2/group/edwards25519"
	"gopkg.in/dedis/kyber.v2/util/encoding"

	// "go.dedis.ch/kyber/v3/group/edwards25519"
	"github.com/O-RD/ths_monorepo/ths"
)

var T_arr []int
var sha256 = edwards25519.NewBlakeSHA256Ed25519().Hash()
var curve = edwards25519.NewBlakeSHA256Ed25519()
var elgamal_Curve = curves.ED25519()

// type Signature struct {
// 	r kyber.Point
// 	s kyber.Scalar
// }

//secure hashing algorithm 256 used for hashing

func Hash(s string) kyber.Scalar {
	sha256.Reset()
	sha256.Write([]byte(s))

	return curve.Scalar().SetBytes(sha256.Sum(nil))
}

func Sign(m string, x kyber.Scalar) ths.Signature {
	// Get the base of the curve.
	g := curve.Point().Base()

	// Pick a random k from allowed set.
	k := curve.Scalar().Pick(curve.RandomStream())

	// r = k * G ( r = g^k)
	r := curve.Point().Mul(k, g)

	// Hash(m || r)
	e := Hash(m + r.String())

	// s = k - e * x
	s := curve.Scalar().Sub(k, curve.Scalar().Mul(e, x))

	ret := ths.Signature{R: r, S: s}
	return ret
}

// func PublicKey(m string, S Signature) kyber.Point {

// 	g := curve.Point().Base()
// 	e := Hash(m + S.r.String())
// 	y := curve.Point().Sub(S.r, curve.Point().Mul(S.s, g))
// 	y = curve.Point().Mul(curve.Scalar().Div(curve.Scalar().One(), e), y)

// 	return y
// }

func Verify(m string, S ths.Signature, y kyber.Point) bool {
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

func Preprocessing() (privateKey kyber.Scalar, publicKey kyber.Point) {
	privateKey = curve.Scalar().Pick(curve.RandomStream())
	publicKey = curve.Point().Mul(privateKey, curve.Point().Base())

	return privateKey, publicKey
}

// func main() {
// 	inputReader := bufio.NewReader(os.Stdin) //for reading multi words from console

// 	privateKey, publicKey := Preprocessing()

// 	fmt.Printf("Private key: %s\n", privateKey)
// 	fmt.Printf("Derived Public key: %s\n\n", publicKey)

// 	fmt.Println("Enter the message to sign") //getting message to sign
// 	var message string
// 	message, _ = inputReader.ReadString('\n')

// 	signature := Sign(message, privateKey)
// 	res := fmt.Sprintf("(r=%s, s=%s)", signature.r, signature.s)
// 	fmt.Printf("Signature %s\n\n", res)

// 	derived_publickey := PublicKey(message, signature)
// 	fmt.Printf("Public key : %s\n\n", publicKey)
// 	fmt.Printf("Derived Public Key? %s\n\n", derived_publickey)
// 	fmt.Printf("Verification Result : %t\n\n", Verify(message, signature, publicKey))

// }

func hash_sign(value []byte) ([]byte, error) {
	h := SHA_256.New()
	h.Write(value)
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	ret, _ := hex.DecodeString(sha1_hash)
	return ret, nil

}

func Signing_T_Unkown(U kyber.Point, x_i kyber.Scalar, Message string, peer_number string) (kyber.Scalar, kyber.Point) {

	file, _ := os.Open("Data/" + peer_number + "/Signing/R_i.txt")
	R_i, _ := encoding.ReadHexScalar(curve, file)
	U_i := curve.Point().Mul(R_i, g)

	// var T int64 = int64(Threshold)
	// j, _ := strconv.Atoi(peer_number)

	Hashing_message := Message + U.String()
	H, _ := hash_sign([]byte(Hashing_message))
	var H1 kyber.Scalar
	H1 = curve.Scalar().Zero()
	H1.SetBytes(H)
	H1 = H1.Mul(H1, x_i) //H1=H*x_i
	// H1 = H1.Mul(H1, Lambda(T, int64(j)))
	V_i := R_i.Add(R_i, H1) //Val= R_i+ H1

	return V_i, U_i
}

func Verify_sign_share(V_i kyber.Scalar, U kyber.Point, U_i kyber.Point, message string, X_i kyber.Point) bool {
	//message, U , V public key
	//V is sum of all V_i's
	//U is sum of all U_i's
	//GK is sum of all alpha[0] (group key)

	t1 := curve.Point().Mul(V_i, g)
	// h := Hash(message + U.String())
	Hashing_message := message + U.String()
	h, _ := hash_sign([]byte(Hashing_message))

	var H1 kyber.Scalar
	H1 = curve.Scalar().Zero()
	H1.SetBytes(h)

	t2 := curve.Point().Mul(H1, X_i)
	t2 = t2.Add(t2, U_i)

	if t1.Equal(t2) {
		return true
	} else {
		return false
	}

}

//When verifying only for signers(Requires Signers Array)
func Lambda(t, j int64, t_array []int) kyber.Scalar {
	var i int64
	den := curve.Scalar().One()
	var LagCoeff = curve.Scalar().One()        //
	var J kyber.Scalar = curve.Scalar().Zero() //Converting j to kyber scalar from int64
	J.SetInt64(j)

	for i = 0; i < int64(len(t_array)); i++ {
		if int64(t_array[i]) == j {
			continue
		}

		var I kyber.Scalar = curve.Scalar().Zero()
		I.SetInt64(int64(t_array[i]))
		den.Sub(I, J)               //den=(i-j)
		den.Inv(den)                //1/(i-j)
		den.Mul(den, I)             //i/(i-j)
		LagCoeff.Mul(LagCoeff, den) // product (i/(i-j)) for each i from 1 to t such that i!=j
	}
	fmt.Println(LagCoeff.String())
	return LagCoeff
}

func combine_T_Unknown(my_index int, Threshold int) (kyber.Scalar, kyber.Point) {
	// Peer_Count := len(peer_details_list)
	T := Threshold
	var Vsum kyber.Scalar = curve.Scalar().Zero()
	// err := os.MkdirAll("Received/Signing/Combine", os.ModePerm)
	// if err != nil {
	// 	panic(err)
	// }
	var Usum kyber.Point = curve.Point().Null()
	var path string
	var path2 string
	peer_number := strconv.Itoa(my_index + 1)

	fmt.Println("T_arr ARRAY:", T_arr)

	fmt.Println("T_arr OF SIGNERS:")

	for i := 0; i < len(T_arr); i++ {
		if T_arr[i] == (my_index + 1) {

			path = "Data/" + peer_number + "/Signing/V_i.txt"
			path2 = "Data/" + peer_number + "/Signing/U_i.txt"

		} else {
			path = "Received/" + fmt.Sprint(T_arr[i]) + "/Signing/V_i.txt"
			path2 = "Received/" + fmt.Sprint(T_arr[i]) + "/Signing/U_i.txt"

		}
		// path := "Broadcast/" + fmt.Sprint(i) + "/Signing/V_i.txt"
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		var j int
		var file2 *os.File
		for j = 0; j < 5; j++ {
			file2, err = os.Open(path2)
			if err != nil {
				time.Sleep(time.Second * 2)
			} else {
				break
			}

		}

		// Lambda_i := lambdas[k]
		Lambda_i := Lambda(int64(T), int64(T_arr[i]), T_arr)
		Lambda_i2 := Lambda_i
		V_i, _ := encoding.ReadHexScalar(curve, file)

		prod := curve.Scalar().Mul(Lambda_i, V_i)
		Vsum = Vsum.Add(Vsum, prod)
		// k += 1

		// // Lambda_i2 := Lambda(int64(T), int64(i))
		// U_i, _ := encoding.ReadHexPoint(curve, file)
		// prod2 := curve.Point().Mul(Lambda_i2, U_i)
		// Usum = Usum.Add(Usum, prod2)
		// path2 := "Data/" + strconv.Itoa(int(i)) + "/Signing/U_i.txt"
		// file2, _ := os.Open(path2)
		// path2 := "Broadcast/" + fmt.Sprint(i) + "/Signing/U_i.txt"

		temp, _ := encoding.ReadHexPoint(curve, file2)

		prod2 := curve.Point().Mul(Lambda_i2, temp)
		Usum = Usum.Add(Usum, prod2)
	}
	fmt.Println("Sum of all V_i:", Vsum.String())
	fmt.Println("Sum of All labda U_i:", Usum.String())
	// file, _ := os.OpenFile("Received/Signing/"+peer_number+"/V.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	// encoding.WriteHexScalar(curve, file, Vsum)
	// file, _ = os.OpenFile("Received/Signing/"+peer_number+"/U.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	// encoding.WriteHexPoint(curve, file, Usum)

	// GK := Get_Group_Key(int64(Peer_Count))
	// // x := Vsum.Clone()
	// // y := Usum.Clone()
	// fmt.Println("INSIDE GKEY:", GK.String())
	// res := Verify_sign_share(Vsum, Usum, Usum, Message, GK)
	// if res {
	// 	fmt.Println("SUCCESS VERIFICATION OF SIGNATURE")
	// } else {
	// 	fmt.Println("INSIDE FAILED TO VERIFIY")
	// }

	return Vsum, Usum

}
