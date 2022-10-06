package keygen

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	crand "crypto/rand"
	SHA "crypto/sha256"

	//"crypto/sha256"

	//"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	//"keygen/ELGAMAL_NEW/AES"
	"log"
	"os"
	"strings"

	"github.com/coinbase/kryptology/pkg/core/curves"
)

//Here we define Public Parameters ( Curve Generator , Receiver's Public Key)

//var KeyPub curves.Point

//To checkc the size of arrays
func sanity_check(target []byte, value int) bool {
	if target == nil { //Check if The given value is equal to nil
		return false
	}
	if len(target)*8 != value { // Check if Bit size of byte array passed is equal to 'value' of size specified
		return false
	}
	return true
}

//Hash the given byte array using SHA256
func hash(value []byte) ([]byte, error) {
	h := SHA.New()
	h.Write(value)
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	ret, _ := hex.DecodeString(sha1_hash)
	if !sanity_check(ret, 256) {
		return []byte(""), errors.New("error in hashing")
	}

	return ret, nil
}

func Setup() *curves.Curve {
	curve := curves.ED25519() // Choosen curve : ED25519
	path := "./Generator.json"
	G := curve.Point.Generator()

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalf("[-] Error Creating Generator File: %s", err)
	}

	_, err = fmt.Fprintln(file, G)
	if err != nil {
		log.Fatalf("[-] Error Writing")
	}

	fmt.Println("[+] Elgamal Setup Completed")
	return curve
}

func Elgamal_KeyGen() (curves.Scalar, curves.Point) { //Generates <Key_pri,KeyPub> Pair
	curve := Setup()
	G := curve.Point.Generator()
	private := curve.Scalar.Random(crand.Reader)
	public := G.Mul(private)
	path := "./PublicParameters.json"

	file, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalf("[-] Error Creating PublicParameters File: %s", err)
	}

	_, err = fmt.Fprintf(file, "Public Key(Hex):\n%x\n", public.ToAffineCompressed())
	if err != nil {
		log.Fatalf("[-] Error Writing")
	}

	path = "./PrivateKey.json"

	file, err = os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalf("[-] Error Creating PrivateKey File: %s", err)
	}

	_, err = fmt.Fprintf(file, "Private Key(Hex):\n%x\n", private.Bytes())
	if err != nil {
		log.Fatalf("[-] Error Writing")
	}

	fmt.Println("[+] KeyGen Completed")
	return private, public
}

func Encryption(msg string, key_pub_r curves.Point) (curves.Point, []byte, error) {

	curve := Setup()

	if len(msg) == 0 {
		return curve.Point.Random(crand.Reader), []byte(""), errors.New("invalid Message")
	}
	//Choose random Number 'r'
	r := curve.Scalar.Random(crand.Reader)

	//C1=P.r
	C1 := curve.Point.Generator().Mul(r)

	//Computing AESKEY = r.Reciever_public_key
	temp := key_pub_r.Mul(r)
	aesKey := temp.ToAffineCompressed()

	// Generating C2: Encrypted Message
	C2, _ := Encrypt(aesKey, msg)

	// Saving the Public Parameters <C1,C2> in a file
	path := "./PublicParameters.json"
	file, _ := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0777)

	_, err := fmt.Fprintf(file, "C1(Hex):\n%x\n", C1.ToAffineCompressed())
	if err != nil {
		log.Fatalf("Error Creating PublicParameters file")
	}
	_, err = fmt.Fprintf(file, "C2(Hex):\n%x\n", C2)
	if err != nil {
		log.Fatalf("Error Writing\n")
	}

	if !sanity_check(r.Point().ToAffineCompressed(), 256) && !sanity_check(C1.ToAffineCompressed(), 256) && !sanity_check(aesKey, 256) {
		return curve.Point.Random(crand.Reader), []byte(""), errors.New("error in encryption")
	}
	//println("INSIDE ENCRYPTION:", C2)

	return C1, C2, nil
}

func AuthEncryption(msg string, key_priv curves.Scalar, key_pub_s curves.Point, key_pub_r curves.Point) (curves.Point, []byte, []byte, error) {
	curve := Setup()

	if len(msg) == 0 {
		return curve.Point.Random(crand.Reader), []byte(""), []byte(""), errors.New("invalid message")
	}
	//Inputs to function-> Curve used , Message to encrypt , Private key
	//Choosing Random Value 'r'
	r := curve.Scalar.Random(crand.Reader)

	//Calculating C1 = P . r
	C1 := curve.Point.Generator().Mul(r)

	//Generating AES_key = r . Reciever_Public_key
	temp := key_pub_r.Mul(r)
	aes_Key, hash_err := hash(temp.ToAffineCompressed()) //Changing curves.point to []byte and then hashing to [256]byte array
	if hash_err != nil {
		return curve.Point.Random(crand.Reader), []byte(""), []byte(""), hash_err //Returns Random Point,empty bytes,error
	}

	//Calculating C2: Encrypted Message
	C2, _ := Encrypt(aes_Key, msg)

	//Calculating Symmetric_key = Sender_Private_Key . Receiver_Public_key
	symm_key := key_pub_r.Mul(key_priv) //Same Symmetric Key used by both parties

	//Combing message, Aes_key, symmetric_key, key_public_sender, key_public_reciever into a single []byte array
	Hashing_message := append([]byte(msg), aes_Key...)
	Hashing_message = append(Hashing_message, symm_key.ToAffineCompressed()...)
	Hashing_message = append(Hashing_message, key_pub_s.ToAffineCompressed()...)
	Hashing_message = append(Hashing_message, key_pub_r.ToAffineCompressed()...)

	C3, hash_err2 := hash(Hashing_message)

	if hash_err2 != nil {
		return curve.Point.Random(crand.Reader), []byte(""), []byte(""), hash_err //Returns Random Point,empty bytes,error
	}

	// Cheking The values Generated:
	if !sanity_check(r.Point().ToAffineCompressed(), 256) && !sanity_check(C1.ToAffineCompressed(), 256) && len(C2) == 0 && !sanity_check(C3, 256) {
		return curve.Point.Random(crand.Reader), []byte(""), []byte(""), hash_err //Returns Random Point,empty bytes,error

	}
	return C1, []byte(C2), C3, nil
}

func Decryption(C1 curves.Point, C2 []byte, key_pri curves.Scalar) ([]byte, error) {

	//fmt.Println("INSIDE DECRYP:", C2)
	if !sanity_check(C1.ToAffineCompressed(), 256) && !sanity_check(key_pri.Point().ToAffineCompressed(), 256) {
		fmt.Printf("SANITY Error")
		return []byte(""), errors.New("invalid input")
	}
	temp_key := C1.Mul(key_pri) //Recovering Symm. key
	aesKey := temp_key.ToAffineCompressed()

	dec, _ := Decrypt(aesKey, string(C2))

	return []byte(dec), nil
}

func AuthDecryption(C1 curves.Point, C2 []byte, C3 []byte, key_pub_s curves.Point, key_pub_r curves.Point, key_pri curves.Scalar) (string, error) {

	if !sanity_check(C1.ToAffineCompressed(), 256) && !sanity_check(key_pri.Point().ToAffineCompressed(), 256) && !sanity_check(key_pub_r.ToAffineCompressed(), 256) && !sanity_check(key_pub_s.ToAffineCompressed(), 256) {
		return "", errors.New("invalid input")
	}

	//Computing R = C1 . key_pri
	R := C1.Mul(key_pri)

	//Generating aes_Key = hash(r)
	aes_Key, hash_err := hash(R.ToAffineCompressed())
	if hash_err != nil {
		return "", hash_err //Returns Random Point,empty bytes,error
	}

	dec, _ := Decrypt(aes_Key, string(C2))

	symm_key := key_pub_s.Mul(key_pri) //Recovering Symm. key

	//Hashing_message:=append()
	Hashing_message := append([]byte(dec), aes_Key...)
	Hashing_message = append(Hashing_message, symm_key.ToAffineCompressed()...)
	Hashing_message = append(Hashing_message, key_pub_s.ToAffineCompressed()...)
	Hashing_message = append(Hashing_message, key_pub_r.ToAffineCompressed()...)

	H, hash_err2 := hash(Hashing_message)
	if hash_err2 != nil {
		return "", hash_err2 //Returns Random Point,empty bytes,error
	}

	if !sanity_check(symm_key.ToAffineCompressed(), 256) {
		return "", errors.New("invalid symmetric key generated") //Returns Random Point,empty bytes,error

	}

	//fmt.Println(temp_C3)
	//fmt.Println(C3)
	//fmt.Println(dec)

	//Checking Authentication , if temp H==C3-> valid , else -> Invalid

	res := bytes.Compare(H, C3)
	if res == 0 {
		return dec, nil
	}

	//Retuning empty msg and invlaid error
	return "", errors.New("Decryption Error")

}

//to add padding to the passed string
func addBase64Padding(value string) string {
	m := len(value) % 4
	if m != 0 {
		value += strings.Repeat("=", 4-m)
	}

	return value
}

//to pad the bytes in src
func Pad(src []byte) []byte {
	//to calculate padding size
	padding := aes.BlockSize - len(src)%aes.BlockSize
	//padding is added to the padtext using repeat function
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

//to unpad the bytes in src
func Unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error")
	}

	return src[:(length - unpadding)], nil
}

//to encrypt the text with the key passed
func Encrypt(key []byte, text string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte(""), err
	}

	msg := Pad([]byte(text))
	//to create a byte array of size blocksize plus length of msg
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	//to read the file with specified blocksize
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte(""), err
	}

	//returns a Stream which encrypts with cipher feedback mode, using the given Block.//returns a Stream which encrypts with cipher feedback mode, using the given Block.
	cfb := cipher.NewCFBEncrypter(block, iv)
	//encrypts texts with xor function
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(msg))
	//padding is removed from the final msg and converted to string
	finalMsg := removeBase64Padding(base64.URLEncoding.EncodeToString(ciphertext))
	return []byte(finalMsg), nil
}

//to decrypt the text with the key passed
func Decrypt(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//padding is done to the ciphertext passed
	decodedMsg, err := base64.URLEncoding.DecodeString(addBase64Padding(text))
	if err != nil {
		return "", err
	}

	if (len(decodedMsg) % aes.BlockSize) != 0 {
		return "", errors.New("blocksize must be multiple of decoded message length")
	}

	iv := decodedMsg[:aes.BlockSize]
	msg := decodedMsg[aes.BlockSize:]

	//decrypts the message
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	//message is unpadded
	unpadMsg, err := Unpad(msg)
	if err != nil {
		return "", err
	}

	//message is returned as a string
	return string(unpadMsg), nil
}

//to remove padding from the passed string
func removeBase64Padding(value string) string {
	return strings.Replace(value, "=", "", -1)
}

// //to pad the bytes in src
// func Pad(src []byte) []byte {
// 	//to calculate padding size
// 	padding := aes.BlockSize - len(src)%aes.BlockSize
// 	//padding is added to the padtext using repeat function
// 	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
// 	return append(src, padtext...)
// }

// //to unpad the bytes in src
// func Unpad(src []byte) ([]byte, error) {
// 	length := len(src)
// 	unpadding := int(src[length-1])

// 	if unpadding > length {
// 		return nil, errors.New("unpad error")
// 	}

// 	return src[:(length - unpadding)], nil
// }

// //to encrypt the text with the key passed
// func AESencrypt(key []byte, text string) (string, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return "", err
// 	}

// 	msg := Pad([]byte(text))
// 	//to create a byte array of size blocksize plus length of msg
// 	ciphertext := make([]byte, aes.BlockSize+len(msg))
// 	iv := ciphertext[:aes.BlockSize]
// 	//to read the file with specified blocksize
// 	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
// 		return "", err
// 	}

// 	//returns a Stream which encrypts with cipher feedback mode, using the given Block.//returns a Stream which encrypts with cipher feedback mode, using the given Block.
// 	cfb := cipher.NewCFBEncrypter(block, iv)
// 	//encrypts texts with xor function
// 	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(msg))
// 	//padding is removed from the final msg and converted to string
// 	finalMsg := removeBase64Padding(base64.URLEncoding.EncodeToString(ciphertext))
// 	return finalMsg, nil
// }

// //to decrypt the text with the key passed
// func AESdecrypt(key []byte, text string) (string, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return "", err
// 	}

// 	//padding is done to the ciphertext passed
// 	decodedMsg, err := base64.URLEncoding.DecodeString(addBase64Padding(text))
// 	if err != nil {
// 		return "", err
// 	}

// 	if (len(decodedMsg) % aes.BlockSize) != 0 {
// 		return "", errors.New("blocksize must be multiple of decoded message length")
// 	}

// 	iv := decodedMsg[:aes.BlockSize]
// 	msg := decodedMsg[aes.BlockSize:]

// 	//decrypts the message
// 	cfb := cipher.NewCFBDecrypter(block, iv)
// 	cfb.XORKeyStream(msg, msg)

// 	//message is unpadded
// 	unpadMsg, err := Unpad(msg)
// 	if err != nil {
// 		return "", err
// 	}

// 	//message is returned as a string
// 	return string(unpadMsg), nil
// }
