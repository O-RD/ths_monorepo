package keygen

import (
	"fmt"
	"io/ioutil"

	"encoding/hex"
	//	elgamal "keygen/ELGAMAL_NEW"
	"math/big"
	"os"
	"strconv"

	"github.com/coinbase/kryptology/pkg/core/curves"
	"gopkg.in/dedis/kyber.v2"

	//	"gopkg.in/dedis/kyber.v2/group/edwards25519"
	"gopkg.in/dedis/kyber.v2/util/encoding"
)

//function to Recieve KGC
func Recieve_KGC(peer_number string) {
	f1, e1 := os.Open("Broadcast/" + peer_number + "/Signature_S.txt")
	if e1 != nil {
		fmt.Println(e1)
	}
	sig, e := encoding.ReadHexScalar(curve, f1)
	if e != nil {
		fmt.Println(e)
	}

	f2, e2 := os.Open("Broadcast/" + peer_number + "/PubKey.txt")
	if e2 != nil {
		fmt.Println(e2)
	}
	pub_key, e_2 := encoding.ReadHexPoint(curve, f2)
	if e_2 != nil {
		fmt.Println(e_2)
	}

	message, e4 := ioutil.ReadFile("Broadcast/" + peer_number + "/Message.txt")
	if e4 != nil {
		fmt.Println(e4)
	}
	f1.Close()
	f2.Close()
	fmt.Printf("Recieved the KGC from peer %s \n", peer_number)
	fmt.Printf("Message: %s \n", string(message))
	fmt.Print("Public Key : ")
	fmt.Println(pub_key)
	fmt.Print("Signature  : ")
	fmt.Println(sig)
	//storing recieved values
	path := "Received/" + peer_number + "/KGC"
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	f3, e3 := os.Create(path + "/Signature.txt")
	if e3 != nil {
		fmt.Println(e3)
	}
	encoding.WriteHexScalar(curve, f3, sig)

	f4, e4 := os.Create(path + "/PubKey.txt")
	if e4 != nil {
		fmt.Println(e4)
	}
	encoding.WriteHexPoint(curve, f4, pub_key)

	f5, e5 := os.Create(path + "/Message.txt")
	if e5 != nil {
		fmt.Println(e5)
	}
	f5.WriteString(string(message))
	f5.Close()
	f4.Close()
	f3.Close()
}

//function to Recieve KGC
func Recieve_KGC_sign(peer_number string) {
	f1, e1 := os.Open("Broadcast/" + peer_number + "/Signing/Signature_S.txt")
	if e1 != nil {
		fmt.Println(e1)
	}
	sig, e := encoding.ReadHexScalar(curve, f1)
	if e != nil {
		fmt.Println(e)
	}

	f2, e2 := os.Open("Broadcast/" + peer_number + "/Signing/PubKey.txt")
	if e2 != nil {
		fmt.Println(e2)
	}
	pub_key, e_2 := encoding.ReadHexPoint(curve, f2)
	if e_2 != nil {
		fmt.Println(e_2)
	}

	message, e4 := ioutil.ReadFile("Broadcast/" + peer_number + "/Signing/Message.txt")
	if e4 != nil {
		fmt.Println(e4)
	}
	f1.Close()
	f2.Close()
	fmt.Printf("Recieved the KGC from peer %s \n", peer_number)
	fmt.Printf("Message: %s \n", string(message))
	fmt.Print("Public Key : ")
	fmt.Println(pub_key)
	fmt.Print("Signature  : ")
	fmt.Println(sig)
	//storing recieved values
	path := "Received/Signing/" + peer_number + "/KGC"
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	f3, e3 := os.Create(path + "/Signature.txt")
	if e3 != nil {
		fmt.Println(e3)
	}
	encoding.WriteHexScalar(curve, f3, sig)

	f4, e4 := os.Create(path + "/PubKey.txt")
	if e4 != nil {
		fmt.Println(e4)
	}
	encoding.WriteHexPoint(curve, f4, pub_key)

	f5, e5 := os.Create(path + "/Message.txt")
	if e5 != nil {
		fmt.Println(e5)
	}
	f5.WriteString(string(message))
	f5.Close()
	f4.Close()
	f3.Close()
}

//Broadcasting Function for KGC
func Broadcast_KGC(peer_number string) {
	f1, e1 := os.Open("Commitment/KGC/Signature_S.txt")
	if e1 != nil {
		fmt.Println(e1)
	}
	sig, e := encoding.ReadHexScalar(curve, f1)
	if e != nil {
		fmt.Println(e)
	}

	f2, e2 := os.Open("Commitment/KGC/PubKey.txt")
	if e2 != nil {
		fmt.Println(e2)
	}
	pub_key, e_2 := encoding.ReadHexPoint(curve, f2)
	if e_2 != nil {
		fmt.Println(e_2)
	}

	message, e4 := ioutil.ReadFile("Commitment/KGC/Message.txt")
	if e4 != nil {
		fmt.Println(e4)
	}
	f1.Close()
	f2.Close()
	//creatinf directory  broadcast/kgc_i to store kgc of ith peer
	path := "Broadcast/KGC" + peer_number
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	//Broadcasting read message ,sig and pub key
	f3, e3 := os.Create(path + "/Signature_S.txt")
	if e3 != nil {
		fmt.Println(e3)
	}
	encoding.WriteHexScalar(curve, f3, sig)
	f4, e4 := os.Create(path + "/Pubkey.txt")
	if e4 != nil {
		fmt.Println(e4)
	}
	encoding.WriteHexPoint(curve, f4, pub_key)
	f5, e5 := os.Create(path + "/Message.txt")
	if e5 != nil {
		fmt.Println(e5)
	}
	f5.WriteString(string(message))
	f5.Close()
	f4.Close()
	f3.Close()
}

//Broadcasting Function for EPK(Elgamal Public Key)
func Broadcast_EPK(peer_number string, EPK curves.Point) {
	//Creating path for broadcast/EPki to store EPK
	path := "Broadcast/EPK"
	err := os.MkdirAll(path, os.ModePerm) //Saves the EPK's t Broadcast folder
	if err != nil {
		fmt.Println(err)
	}
	path = path + "/epk" + fmt.Sprint(peer_number) + ".txt" // Saves as 'Broadcast/EPK/EPKi.txt'
	file, err2 := os.Create(path)
	if err2 != nil {
		fmt.Println(err2)
	}

	_, err = fmt.Fprintf(file, "%x\n", EPK.ToAffineCompressed()) //Saving the EPK in hex format
	if err != nil {
		fmt.Println(err)
	}
	file.Close()
}

func Recieve_EPK(peer_number string) {
	//Creating path for broadcast/EPki to store EPK
	path := "Broadcast/EPK/epk" + fmt.Sprint(peer_number) + ".txt"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	data, _ = hex.DecodeString(string(data))
	elg_curve := Setup()

	EPK_generated, _ := elg_curve.Point.FromAffineCompressed(data)
	path = "Received"
	f1, e1 := os.Create(path + "/EPK" + peer_number + ".txt")
	if e1 != nil {
		fmt.Println(e1)
	}
	_, err = fmt.Fprintf(f1, "%x\n", EPK_generated.ToAffineCompressed())
	if err != nil {
		fmt.Println(err)
	}

}

//Broadcasting Function for KGD

func Broadcast_KGD(peer_number string) {
	//reading KGD from commitment
	f1, e1 := os.Open("Commitment/KGD.txt")
	if e1 != nil {
		fmt.Println(e1)
	}
	KGD_i, e_1 := encoding.ReadHexPoint(curve, f1)
	if e_1 != nil {
		panic(e_1)
	}
	//Broadcasting KGD_i
	path := "Broadcast/KGD" + peer_number + ".txt"
	f2, e2 := os.Create(path)
	if e2 != nil {
		panic(e2)
	}
	encoding.WriteHexPoint(curve, f2, KGD_i)
	f2.Close()
	f1.Close()
}

func Recieve_KGD(peer_number string) {
	path := "Broadcast/" + peer_number + "/KGD.txt"
	f1, e1 := os.Open(path)
	if e1 != nil {
		panic(e1)
	}
	KGD_i, err := encoding.ReadHexPoint(curve, f1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Recieved value of KGD from peer %s \n", peer_number)
	fmt.Println(KGD_i)
	f1.Close()
	//storing KGD
	path = "Received/" + peer_number
	f2, e2 := os.Create(path + "/KGC/KGD" + ".txt")
	if e2 != nil {
		fmt.Println(e2)
	}
	encoding.WriteHexPoint(curve, f2, KGD_i)
}
func Recieve_KGD_sign(peer_number string) {
	path := "Broadcast/" + peer_number + "/Signing/KGD.txt"
	f1, e1 := os.Open(path)
	if e1 != nil {
		panic(e1)
	}
	KGD_i, err := encoding.ReadHexPoint(curve, f1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Recieved value of KGD from peer %s \n", peer_number)
	fmt.Println(KGD_i)
	f1.Close()
	//storing KGD
	path = "Received/Signing/" + peer_number
	f2, e2 := os.Create(path + "/KGC/KGD" + ".txt")
	if e2 != nil {
		fmt.Println(e2)
	}
	encoding.WriteHexPoint(curve, f2, KGD_i)
}

func toHexInt(n *big.Int) string {
	return fmt.Sprintf("%X", n) // or %x or upper case
}

func Brodcast_Alphas(peer_number string, alphas []*big.Int, k int64) {
	path := "Broadcast/alphas" + peer_number + ".txt"
	f1, e1 := os.Create(path)
	if e1 != nil {
		panic(e1)
	}
	var i int64
	for i = 0; i < k; i++ {
		f1.WriteString(fmt.Sprintf("%d", alphas[i]))
	}
	f1.Close()
}

// func Recieve_Alphas(peer_number string, k int64) {
// 	var i int64

// 	for i = 0; i < k; i++ {

// 		path := "Broadcast/" + peer_number + "/Alphas/alpha" + fmt.Sprint(i) + ".txt"
// 		data, err := os.ReadFile(path)
// 		if err != nil {
// 			panic(err)
// 		}
// 		b := string(data)
// 		path1 := "Received/" + peer_number + "/ALPHAS"
// 		err1 := os.MkdirAll(path1, os.ModePerm)
// 		if err != nil {
// 			panic(err1)
// 		}
// 		f2, err2 := os.Create(path1 + "/alpha" + fmt.Sprint(i) + ".txt")
// 		if err2 != nil {
// 			panic(err2)
// 		}
// 		f2.WriteString(b)
// 		f2.Close()
// 	}
// }

func Recieve_Alphas(Peer_Count int64, peer_number string, k int64, my_index int) {
	var i int64
	var j int64
	for j = 1; j <= Peer_Count; j++ {
		if j == int64(my_index+1) {
			continue
		}
		path1 := "Received/" + peer_number + "/Alphas/" + fmt.Sprint(j)
		err1 := os.MkdirAll(path1, os.ModePerm)
		if err1 != nil {
			panic(err1)
		}
		for i = 0; i < k; i++ {

			path := "Broadcast/" + fmt.Sprint(j) + "/Alphas/alpha" + fmt.Sprint(i) + ".txt"
			data, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}
			b := string(data)

			f2, err2 := os.Create(path1 + "/alpha" + fmt.Sprint(i) + ".txt")
			if err2 != nil {
				panic(err2)
			}
			f2.WriteString(b)
			f2.Close()
		}
	}
}

func Recieve_Alphas_sign(Peer_Count int64, peer_number string, k int64, my_index int) {
	var i int64
	var j int64
	for j = 1; j <= Peer_Count; j++ {
		if j == int64(my_index+1) {
			continue
		}
		path1 := "Received/Signing/" + peer_number + "/Alphas/" + fmt.Sprint(j)
		err1 := os.MkdirAll(path1, os.ModePerm)
		if err1 != nil {
			panic(err1)
		}

		for i = 0; i < k; i++ {

			path := "Broadcast/" + fmt.Sprint(j) + "/Signing/Alphas/alpha" + fmt.Sprint(i) + ".txt"
			data, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}
			b := string(data)

			// f2, err2 := os.OpenFile(path1+"/alpha"+fmt.Sprint(i)+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
			f2, err2 := os.Create(path1 + "/alpha" + fmt.Sprint(i) + ".txt")
			if err2 != nil {
				panic(err2)
			}
			f2.WriteString(b)
			f2.Close()
		}
	}
}

// func Recieve_Alphas_sign(peer_number string, k int64) {
// 	var i int64
// 	var j int64
// 	for j=0;j<=Peer_Count;j++{

// 	for i = 0; i < k; i++ {
// 		path := "Broadcast/" + peer_number + "/Signing/Alphas/alpha" + fmt.Sprint(i) + ".txt"
// 		data, err := os.ReadFile(path)
// 		if err != nil {
// 			panic(err)
// 		}
// 		b := string(data)
// 		path1 := "Received/Signing/" + peer_number + "/ALPHAS"
// 		err1 := os.MkdirAll(path1, os.ModePerm)
// 		if err != nil {
// 			panic(err1)
// 		}
// 		f2, err2 := os.Create(path1 + "/alpha" + fmt.Sprint(i) + ".txt")
// 		if err2 != nil {
// 			panic(err2)
// 		}
// 		f2.WriteString(b)
// 		f2.Close()
// 	}
// }

//To get the summation of all U_i's
func Get_U(peer_Count, peer_number int64) kyber.Point {
	U := curve.Point().Null()
	path := "Broadcast/U/U_"
	var i int64
	for i = 1; i <= peer_Count; i++ {
		f2, err := os.Open(path + strconv.Itoa(int(i)) + ".txt")
		if err != nil {
			panic(err)
		}
		// f_2, _ := encoding.ReadHexScalar(curve, f2)
		// val, _ := encoding.ScalarToStringHex(curve, f_2)
		// temp := new(big.Int)
		// temp.SetString(val, 16)
		// U = U.Add(U, temp)
		U_i, _ := encoding.ReadHexPoint(curve, f2)
		//fmt.Print("***U_i Read: ", U_i, "\n")
		if U == curve.Point().Null() {
			U = U_i
		} else {
			U = U.Add(U, U_i)

		}

		//fmt.Print("U before returning:", U, "\n")

	}
	return U
}

func Get_G(peer_Count, peer_number int64) *big.Int {
	G_i := new(big.Int)
	path := "Private/shares" + fmt.Sprint(peer_number)

	var i int64
	for i = 1; i <= int64(peer_Count); i++ {
		path2 := path + "/share" + fmt.Sprint(i) + ".txt"
		data, _ := ioutil.ReadFile(path2)

		temp := new(big.Int)
		temp, _ = temp.SetString(string(data), 10)
		//fmt.Println("temp", temp)
		//fmt.Print("TEMP->", temp)

		G_i = G_i.Add(G_i, temp)
	}
	//fmt.Print("G_I wehn returning", G_i, "\n")
	return G_i
}

//To return the summation of all the G_i's
func Get_R(peer_Count, peer_number int64) *big.Int {
	var G_i *big.Int
	path := "Broadcast/VerificationSet/SetBy"
	var i int64
	//G_i, _ = G_i.SetString("0", 10)
	for i = 1; i <= int64(peer_Count); i++ {
		path2 := path + fmt.Sprint(i) + "/F_" + fmt.Sprint(i) + "(" + fmt.Sprint(peer_number) + ")" + ".txt"
		//fmt.Print(path2, "\n\n")
		data, _ := ioutil.ReadFile(path2)

		temp := new(big.Int)
		temp, _ = temp.SetString(string(data), 10)

		//fmt.Print("TEMP->", temp)
		if G_i == nil {
			G_i = temp
			//fmt.Print("G_I->", G_i)
			continue
		}
		G_i = G_i.Add(G_i, temp)

	}
	// G_i is summation of shares recieved from other peers in the network
	return G_i
}

func Store_Scaler(x kyber.Scalar, path, filename string) error {
	os.MkdirAll(path, os.ModePerm)
	file, err := os.OpenFile(path+"/"+filename+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	encoding.WriteHexScalar(curve, file, x)
	if err != nil {
		return err
	}
	return nil
}

func Store_Point(x kyber.Point, path, filename string) error {
	os.MkdirAll(path, os.ModePerm)
	file, err := os.OpenFile(path+"/"+filename+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	encoding.WriteHexPoint(curve, file, x)
	if err != nil {
		return err
	}
	return nil
}

func Get_Stored_Point(path string) kyber.Point {
	file, _ := os.OpenFile(path, os.O_RDONLY, 0777)

	x, err := encoding.ReadHexPoint(curve, file)

	if err != nil {
		fmt.Println("error reading kyber Point")
		return nil
	}
	return x
}

func Get_Stored_Scaler(path string) kyber.Scalar {
	file, _ := os.OpenFile(path, os.O_RDONLY, 0777)

	x, err := encoding.ReadHexScalar(curve, file)

	if err != nil {
		fmt.Println("error reading kyber Scaler")
		return nil
	}
	return x
}

func Get_ESK(path string) (curves.Scalar, error) {
	elg_curve := Setup()
	_ESK_file, err := ioutil.ReadFile(path)
	_ESK_file, _ = hex.DecodeString(string(_ESK_file))
	_ESK, err2 := elg_curve.Scalar.SetBytes(_ESK_file)
	if err != nil || err2 != nil {
		if err != nil {
			return nil, err
		} else {
			return nil, err2
		}
	}

	return _ESK, nil
}

func Get_EPK(path string) (curves.Point, error) {
	elg_curve := Setup()
	_EPK_file, err := ioutil.ReadFile(path)
	_EPK_file, _ = hex.DecodeString(string(_EPK_file))
	_EPK, err2 := elg_curve.Point.FromAffineCompressed(_EPK_file)
	if err != nil || err2 != nil {
		if err != nil {
			return nil, err
		} else {
			return nil, err2
		}
	}

	return _EPK, nil

}

func Get_Group_Key(Peer_Count int64, my_index int) kyber.Point {
	var i int64
	var GK kyber.Point
	GK = curve.Point().Null()
	for i = 1; i <= Peer_Count; i++ {
		if i == int64(my_index+1) {
			path2 := "vss/" + fmt.Sprint(i) + "/alpha0.txt"

			file, _ := os.Open(path2)
			temp, _ := encoding.ReadHexPoint(curve, file)
			GK = GK.Add(GK, temp)
			continue
		}
		path := "Broadcast/" + fmt.Sprint(i) + "/Alphas/alpha0.txt"
		file, _ := os.Open(path)
		temp, _ := encoding.ReadHexPoint(curve, file)
		GK = GK.Add(GK, temp)
	}
	return GK
}
