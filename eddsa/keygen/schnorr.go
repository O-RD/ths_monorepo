package keygen

import (
	"gopkg.in/dedis/kyber.v2"
	"gopkg.in/dedis/kyber.v2/group/edwards25519"

	// "go.dedis.ch/kyber/v3/group/edwards25519"
	"github.com/O-RD/ths_monorepo/ths"
)

var sha256 = edwards25519.NewBlakeSHA256Ed25519().Hash()
var curve = edwards25519.NewBlakeSHA256Ed25519()

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