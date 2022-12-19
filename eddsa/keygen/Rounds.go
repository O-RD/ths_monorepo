package keygen

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/O-RD/ths_monorepo/ths"
	"gopkg.in/dedis/kyber.v2"

	// "gopkg.in/dedis/kyber.v2/util/encoding"
	"gopkg.in/dedis/kyber.v2/util/encoding"
	Encode "gopkg.in/dedis/kyber.v2/util/encoding"
)

func Round1(send_chan chan ths.Message, p *ths.P2P, receive_chan chan ths.Payload, Round1_Values *ths.Keygen_Store) {

	// var protocolID protocol.ID = "/keygen/0.0.1"

	//Generate broadcast wait time
	time.Sleep(time.Second * 5)
	// fmt.Println(GeneratePrime(1024))

	peer_number := fmt.Sprint(p.My_Index + 1)
	Peer_Count := len(p.Sorted_Peers)
	fmt.Println("PEERCOUNT:", Peer_Count)
	fmt.Println("MYINDEX:", peer_number)
	var T int64 = int64(p.Threshold)
	fmt.Println("THRESHOLD:", T)

	// fmt.Scan(&choice)
	time.Sleep(time.Second)
	// fmt.Println(choice)
	// switch choice {
	// case 1:
	time.Sleep(time.Second * 2)
	// Generation of Elgamal Keys
	ESK, EPK := Elgamal_KeyGen()

	Round1_Values.EPK = fmt.Sprintf("%x", EPK.ToAffineCompressed())
	Round1_Values.ESK = string(ESK.Bytes())

	fmt.Println(" \n ")
	fmt.Println("Elgamal Public Key:")
	fmt.Println(&EPK)
	fmt.Println("Elgamal Secret Key:")
	fmt.Println(&ESK)
	fmt.Printf("\n")

	//Generating Schnorr Public and  Secret Key
	SSK, SPK := Preprocessing()
	fmt.Println("Schnorr Public Key:")
	fmt.Println(SPK)
	fmt.Println("Schnorr Secret Key:")
	fmt.Println(SSK)
	fmt.Printf("\n")

	Round1_Values.SSK, _ = Encode.ScalarToStringHex(curve, SSK)
	Round1_Values.SPK, _ = Encode.PointToStringHex(curve, SPK)

	//storing the schnorr secret key to Prvate Folder

	//commiting SSK
	Commitment(SSK, "hello world", peer_number, Round1_Values)

}

func Round2(send_chan chan ths.Message, p *ths.P2P, receive_chan chan ths.Payload, Round2_Values *ths.Keygen_Store) {

	peer_number := fmt.Sprint(p.My_Index + 1)
	Peer_Count := len(p.Sorted_Peers)
	fmt.Println("PEERCOUNT:", Peer_Count)
	fmt.Println("MYINDEX:", p.My_Index)
	fmt.Println("PEER_NUMBER:", peer_number)
	var T int64 = int64(p.Threshold)
	fmt.Println("THRESHOLD:", T)

	poly := []kyber.Scalar{}  // to store coefficients
	share := []kyber.Scalar{} // to store share
	alphas := []kyber.Point{} // to store alphas

	var i int64

	for i = 0; i < T; i++ {
		poly = append(poly, curve.Scalar().Zero())
	}

	for i = 0; i < T; i++ {
		alphas = append(alphas, curve.Point().Null())
	}

	for i = 1; i <= int64(Peer_Count); i++ {
		share = append(share, curve.Scalar().Zero())
	}

	// to generate coefficients of the polynomial
	SSK, _ := Encode.StringHexToScalar(curve, Round2_Values.SSK)
	Generate_Polynomial_coefficients(T, poly, peer_number, SSK, "vss/"+peer_number)

	// Generating the shares and storing in share array
	Generate_share(int64(Peer_Count), T, poly, share, peer_number, "vss/"+peer_number)
	//Generating Alphas
	Generate_Alphas(T, alphas, poly, peer_number, "vss/"+peer_number)
	// Round2_Values.Alphas = alphas
	for i = 0; i < T; i++ {
		Round2_Values.Poly[i], _ = Encode.ScalarToStringHex(curve, poly[i])
		Round2_Values.Alphas[i], _ = Encode.PointToStringHex(curve, alphas[i])
	}
	for i = 0; i < int64(Peer_Count); i++ {
		Round2_Values.Shares[i], _ = Encode.ScalarToStringHex(curve, share[i])
	}

}

func Round3(send_chan chan ths.Message, p *ths.P2P, receive_chan chan ths.Payload, Round3_Values *ths.Keygen_Store) {
	peer_number := fmt.Sprint(p.My_Index + 1)
	Peer_Count := len(p.Sorted_Peers)
	fmt.Println("PEERCOUNT:", Peer_Count)
	fmt.Println("PEERNUMBER:", peer_number)
	var T int64 = int64(p.Threshold)
	fmt.Println("THRESHOLD:", T)

	//Receiving alphas from other peers
	// Recieve_Alphas(int64(Peer_Count), peer_number, T, p.My_Index)

	fmt.Println("ENCRYPTING SHARES & Broadcasting")

	//Defining Elgamal Curve
	elg_curve := Setup()

	//Reading Sender's elgamal Public key

	// path13 := "Data/" + peer_number
	// Sender_EPK_file, _ := ioutil.ReadFile(path13 + "/EPK.txt")
	// Sender_EPK_file, _ = hex.DecodeString(string(Sender_EPK_file))
	// Sender_EPK, _ := elg_curve.Point.FromAffineCompressed(Sender_EPK_file)

	temp, _ := hex.DecodeString(Round3_Values.EPK)
	Sender_EPK, _ := elg_curve.Point.FromAffineCompressed(temp)

	//Reading Sender's elgamal Secret Key
	// Sender_ESK_file, _ := ioutil.ReadFile(path13 + "/ESK.txt")
	temp, _ = hex.DecodeString(Round3_Values.ESK)
	Sender_ESK, _ := elg_curve.Scalar.SetBytes(temp)

	//Path to vss generated parameters
	// path3 := "vss/" + peer_number

	var i int64

	for i = 1; i <= int64(Peer_Count); i++ {
		if i == int64(p.My_Index+1) {
			continue
		}

		// _f, err := os.Open(path3 + "/Indivisual_Share" + strconv.Itoa(int(i)) + ".txt")
		// if err != nil {
		// 	panic(err)
		// }

		//share for ith user generated by current peer
		share, _ := Encode.StringHexToScalar(curve, Round3_Values.Shares[i-1])

		//Reading Elgamal Public key of ith user
		data, err := os.ReadFile("Data/" + "INSERT HERE" + "/keys/EPK.txt")
		if err != nil {
			fmt.Println(err)
		}
		data, _ = hex.DecodeString(string(data))
		elg_curve := Setup()
		EPK_receiver, _ := elg_curve.Point.FromAffineCompressed(data)

		// fmt.Println("Sender_ESK:", Sender_ESK)
		// fmt.Println("Sender_EPK:", Sender_EPK)
		// fmt.Println("EPK_receiver:", EPK_receiver)
		//encoding.WriteHexScalar(curve, file, shares[i])
		// toEncrypt, _ := encoding.ScalarToStringHex(curve, share)

		//Share to encrypt(in string format)
		toEncrypt := share.String()
		fmt.Println("TO ENCRYPT:", toEncrypt)

		//Ecryption using( current peer's Secret key,current peer's public key, ith users(receivers) public key)
		C1, C2, C3, _ := AuthEncryption(toEncrypt, Sender_ESK, Sender_EPK, EPK_receiver)
		encrypted := ths.Encrypted_Share{
			C1: string(C1.ToAffineCompressed()),
			C2: C2,
			C3: C3,
		}

		Round3_Values.Encrypted_Shares = append(Round3_Values.Encrypted_Shares, encrypted)

	}
	// wait_until(7)

}
func Round4(send_chan chan ths.Message, p *ths.P2P, receive_chan chan ths.Payload, Round3_Values *ths.Keygen_Store) {

	peer_number := fmt.Sprint(p.My_Index + 1)
	Peer_Count := len(p.Sorted_Peers)
	fmt.Println("PEERCOUNT:", Peer_Count)
	fmt.Println("PEERNUMBER:", peer_number)
	var T int64 = int64(p.Threshold)
	fmt.Println("THRESHOLD:", T)

	path4 := "Data/" + peer_number
	// Reciever_EPK_file, _ := ioutil.ReadFile(path + "/EPK.txt")
	// Reciever_EPK_file, _ = hex.DecodeString(string(Reciever_EPK_file))
	// Reciever_EPK, _ := elg_curve.Point.FromAffineCompressed(Reciever_EPK_file)

	Reciever_EPK, _ := Get_EPK(path4 + "/EPK.txt")

	//Reading elgamal Secret key of current peer
	// Reciever_ESK_file, _ := ioutil.ReadFile(path + "/ESK.txt")
	// Reciever_ESK_file, _ = hex.DecodeString(string(Reciever_ESK_file))
	// Reciever_ESK, _ := elg_curve.Scalar.SetBytes(Reciever_ESK_file)

	Reciever_ESK, _ := Get_ESK(path4 + "/ESK.txt")

	var i int64
	for i = 1; i <= int64(Peer_Count); i++ {
		if i == int64(p.My_Index+1) {
			continue
		}

		//Reading Elgamal Public key of Sender(i)
		path2 := "Received/" + fmt.Sprint(i) + "/Keys/EPK.txt" //Get the epk of sender
		// data, err := ioutil.ReadFile(path2)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// data, _ = hex.DecodeString(string(data))
		// Sender_EPK, _ := elg_curve.Point.FromAffineCompressed(data)
		Sender_EPK, _ := Get_EPK(path2)

		//Reading the Encrypted data sent by ith user to current user from broadcast folder(channel)
		C1_Data, _ := ioutil.ReadFile("Received/" + fmt.Sprint(i) + "/Keygen_shares/C1.txt")
		C2_Data, _ := ioutil.ReadFile("Received/" + fmt.Sprint(i) + "/Keygen_shares/C1.txt")
		C3_Data, _ := ioutil.ReadFile("Received/" + fmt.Sprint(i) + "/Keygen_shares/C1.txt")

		// C1_Data, _ := ioutil.ReadFile("C's/" + fmt.Sprint(i) + "/" + peer_number + "/C1.txt")
		// C2_Data, _ := ioutil.ReadFile("C's/" + fmt.Sprint(i) + "/" + peer_number + "/C2.txt")
		// C3_Data, _ := ioutil.ReadFile("C's/" + fmt.Sprint(i) + "/" + peer_number + "/C3.txt")

		//Changing the data read into C1,C2,C3 format
		data, _ := hex.DecodeString(string(C1_Data))
		C1, _ := elgamal_Curve.Point.FromAffineCompressed(data)
		C2 := C2_Data
		C3, _ := hex.DecodeString(string(C3_Data))
		// C3 := C3_Data

		//Decryption of shares(C1,C2,C3)
		time.Sleep(time.Second * 1)
		share, err := AuthDecryption(C1, C2, C3, Sender_EPK, Reciever_EPK, Reciever_ESK)
		fmt.Println("DECRPYPteD:", share)

		if err != nil {
			fmt.Println("Error Decrypting")
		}

		//Saving the decrypted message into the received folder of current user
		path2 = "Data/Keygen_shares/"
		os.MkdirAll(path2, os.ModePerm)
		file, _ := os.OpenFile(path2+"share"+strconv.Itoa(int(i))+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		_, _ = fmt.Fprint(file, share)

		fmt.Println(peer_number, "Verifying Shares")

		G := Verify_Share(peer_number, int64(Peer_Count), int64(T), false, p.My_Index)
		fmt.Println("Private Key Share:", G.String(), "\n")
		path5 := "Data"
		os.MkdirAll(path5, os.ModePerm)
		file5, _ := os.Create(path5 + "/G.txt")
		encoding.WriteHexScalar(curve, file5, G)

		// if verify_GK(int64(Peer_Count), T) {
		// 	fmt.Println("VERIFIED G")
		// } else {
		// 	fmt.Println("NOT VERIFIED G")
		// }
		//BROADCAST GROUP PUBLIC KEY

		//G-> input to sign t unknwn
		GK := Get_Group_Key(int64(Peer_Count), p.My_Index)
		file5, _ = os.Create(path5 + "/GroupKey.txt")
		encoding.WriteHexPoint(curve, file5, GK)
		fmt.Println("GROUP KEY:", GK.String())
	}

}
