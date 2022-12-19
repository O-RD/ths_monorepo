package keygen

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"gopkg.in/dedis/kyber.v2"
	"gopkg.in/dedis/kyber.v2/util/encoding"
)

// var secret kyber.Scalar      //secret key
var secret_sign kyber.Scalar //secret key
var g = curve.Point().Base() //Generator

// Generating polynomial and storing coefficients
func Generate_Polynomial_coefficients(T int64, poly []kyber.Scalar, peer_number string, secret kyber.Scalar, path string) {
	//storing coefficients in a text file
	// path := "vss/" + peer_number
	err := os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}
	path += "/coeff"
	//poly := make([]kyber.Scalar, T) //Declaring the array to store coeff
	poly[0] = secret //constant term of the polynomial will be our secret
	// fmt.Println("Secret:", secret, "-->", poly[0])
	var i int64
	file, _ := os.Create(path + "0.txt")
	encoding.WriteHexScalar(curve, file, poly[0])
	file.Close()
	for i = 1; i < T; i++ {
		file, _ := os.Create(path + strconv.Itoa(int(i)) + ".txt")
		poly[i] = curve.Scalar().Pick(curve.RandomStream()) //generating random value for coefficients
		//writing coeffiecient to the text file
		encoding.WriteHexScalar(curve, file, poly[i])
		file.Close()
	}
	fmt.Printf("Coefficients Generation for Peer %s Completed \n", peer_number)
}

// Generating polynomial and storing coefficients
func Generate_Polynomial_coefficients_sign(T int64, poly []kyber.Scalar, peer_number string) {
	//storing coefficients in a text file
	path := "vss/Signing/" + peer_number
	err := os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}
	path += "/coeff"
	//poly := make([]kyber.Scalar, T) //Declaring the array to store coeff
	poly[0] = secret_sign //constant term of the polynomial will be our secret
	fmt.Println("Secret:", secret_sign, "-->", poly[0])
	var i int64
	file, _ := os.Create(path + "0.txt")
	encoding.WriteHexScalar(curve, file, poly[0])
	file.Close()
	for i = 1; i < T; i++ {
		file, _ := os.Create(path + strconv.Itoa(int(i)) + ".txt")
		poly[i] = curve.Scalar().Pick(curve.RandomStream()) //generating random value for coefficients
		//writing coeffiecient to the text file
		encoding.WriteHexScalar(curve, file, poly[i])
		file.Close()
	}
	fmt.Printf("Sign Coefficients Generation for Peer %s Completed \n", peer_number)
}

// Calculates and return F(i)
func f_of_i(i int64, T int64, coeff []kyber.Scalar) kyber.Scalar {
	var j int64
	var val kyber.Scalar
	val = curve.Scalar().Zero()
	for j = 0; j < T; j++ {
		// calculating x^j
		X := math.Pow(float64(i), float64(j))
		x := curve.Scalar().SetInt64(int64(X))
		var v2 kyber.Scalar
		v2 = curve.Scalar().Zero()
		coeff_j := coeff[j]
		v2.Mul(x, coeff_j) // multiplying coefficient of j with x^j *coeff[j]
		val.Add(val, v2)   //adding the values from (j-1)*x^(j-1)
	}
	return val
}

func Generate_share(N int64, T int64, coeff []kyber.Scalar, shares []kyber.Scalar, peer_number string, path string) {
	var i int64
	// path := "vss/" + peer_number
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	path += "/Indivisual_Share"
	for i = 1; i <= N; i++ {
		file, _ := os.Create(path + strconv.Itoa(int(i)) + ".txt")
		shares[i-1] = f_of_i(i, T, coeff)
		//fmt.Println(share)                          //to calculate f(i)
		encoding.WriteHexScalar(curve, file, shares[i-1]) //writing share to the file
		file.Close()
	}
	fmt.Printf("Share Generation for Peer %s Completed \n", peer_number)
}

func Generate_share_sign(N int64, T int64, coeff []kyber.Scalar, shares []kyber.Scalar, peer_number string) {
	var i int64
	path := "vss/Signing/" + peer_number
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	path += "/Indivisual_Share"
	for i = 0; i <= N; i++ {
		file, _ := os.Create(path + strconv.Itoa(int(i)) + ".txt")
		shares[i] = f_of_i(i, T, coeff)
		//fmt.Println(share)                          //to calculate f(i)
		encoding.WriteHexScalar(curve, file, shares[i]) //writing share to the file
		file.Close()
	}
	fmt.Printf("Sign Share Generation for Peer %s Completed \n", peer_number)
}

//calculates coeff_i*G
func Generate_Alpha_i(coeff_i kyber.Scalar) kyber.Point {
	val := g
	val.Mul(coeff_i, val)
	return val
}

//Generates Alphas for each coefficient i.e alpha_i= coeff_i *G
func Generate_Alphas(T int64, alphas []kyber.Point, coeff []kyber.Scalar, peer_number string, path string) {
	// path := "vss/" + peer_number
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	path += "/alpha"
	var i int64
	for i = 0; i < T; i++ {
		//fmt.Println(coeff[i])
		file, e1 := os.Create(path + strconv.Itoa(int(i)) + ".txt")
		if e1 != nil {
			panic(e1)
		}
		var val kyber.Point
		val = curve.Point().Null()
		//coeff_i := Readcoeff(i)
		coeff_i := coeff[i]
		val.Mul(coeff_i, g)
		alphas[i] = val
		encoding.WriteHexPoint(curve, file, val)
		file.Close()
	}
	fmt.Printf("Alpha Generation for Peer %s Completed \n", peer_number)
}

//Generates Alphas for each coefficient i.e alpha_i= coeff_i *G
func Generate_Alphas_sign(T int64, alphas []kyber.Point, coeff []kyber.Scalar, peer_number string) {
	path := "vss/Signing/" + peer_number
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	path += "/alpha"
	var i int64
	for i = 0; i < T; i++ {
		//fmt.Println(coeff[i])
		file, e1 := os.Create(path + strconv.Itoa(int(i)) + ".txt")
		if e1 != nil {
			panic(e1)
		}
		var val kyber.Point
		val = curve.Point().Null()
		//coeff_i := Readcoeff(i)
		coeff_i := coeff[i]
		val.Mul(coeff_i, g)
		alphas[i] = val
		encoding.WriteHexPoint(curve, file, val)
		file.Close()
	}
	fmt.Printf("Sign Alpha Generation for Peer %s Completed \n", peer_number)
}

// To verify if the share i is correct or not
// func Verify_i2(i int64, f_i kyber.Scalar, T int64, alphas []kyber.Point) bool {
// 	//fmt.Println("we are in vss")

// 	var v1 kyber.Point
// 	v1 = curve.Point().Null()
// 	// share_i := Readshare(i)
// 	// share_i := f_i
// 	fmt.Print("SHARE INSIDE VERIFY ", f_i.String())
// 	// g2 := curve.Point().Base()
// 	v1.Mul(f_i, g)
// 	var j int64
// 	var sum kyber.Point
// 	sum = curve.Point().Null()
// 	for j = 0; j < T; j++ {
// 		X := math.Pow(float64(i), float64(j))
// 		x := curve.Scalar().SetInt64(int64(X))
// 		var v4 kyber.Point = curve.Point().Null()
// 		// v4 = Readalpha(j)
// 		v4 = alphas[j]
// 		var val kyber.Point
// 		val = curve.Point().Null()
// 		val.Mul(x, v4)
// 		sum.Add(val, sum)
// 	}
// 	fmt.Println(v1)
// 	fmt.Println(sum)
// 	return v1.Equal(sum)
// }
func Verify_i(i int64, f_i kyber.Scalar, T int64, alphas []kyber.Point) bool {
	//fmt.Println("we are in vss")
	var v1 kyber.Point
	v1 = curve.Point().Null()
	// share_i := Readshare(i)
	share_i := f_i
	v1 = v1.Mul(share_i, g)
	var j int64
	var sum kyber.Point
	sum = curve.Point().Null()
	for j = 0; j < T; j++ {
		X := math.Pow(float64(i), float64(j))
		x := curve.Scalar().SetInt64(int64(X))
		var v4 kyber.Point = curve.Point().Null()
		// v4 = Readalpha(j)
		v4 = alphas[j]
		fmt.Println("ALPHAS\n", alphas[j])
		var val kyber.Point
		val = curve.Point().Null()
		val.Mul(x, v4)
		sum.Add(val, sum)
	}

	fmt.Println(v1)
	fmt.Println(sum)
	return v1.Equal(sum)
}

//Sets tha value of secret of the polynomial
// func Set_secret(val kyber.Scalar) {
// 	secret = val
// }
func Set_secret_sign(val kyber.Scalar) {
	secret_sign = val
}
func Verify_Share(peer_number string, N int64, T int64, signing bool, my_index int) kyber.Scalar {
	var i int64
	path := ""
	path2 := ""
	if signing {
		path = "Received/Signing/" + peer_number + "/Shares/share"
	} else {
		path = "Data/Keygen_shares/"
	}
	if signing {
		path2 = "vss/Signing/" + peer_number + "/Indivisual_Share" + peer_number + ".txt"
	} else {
		path2 = "vss/" + peer_number + "/Indivisual_Share" + peer_number + ".txt"
	}

	share := []kyber.Scalar{} // to store share
	for i = 0; i < int64(N); i++ {
		share = append(share, curve.Scalar().Zero())
	}

	for i = 1; i <= N; i++ {
		if i == int64(my_index+1) { //To Add own share to G_i

			file, err := os.Open(path2)
			if err != nil {
				panic(err)
			}
			val, e1 := encoding.ReadHexScalar(curve, file)
			if e1 != nil {
				panic(e1)
			}
			share[i-1] = val
			continue

		}

		file, err := os.Open(path + "share" + strconv.Itoa(int(i)) + ".txt")
		if err != nil {
			panic(err)
		}
		val, e1 := encoding.ReadHexScalar(curve, file)
		if e1 != nil {
			panic(e1)
		}
		share[i-1] = val
	}

	verify_each_share(peer_number, N, share, T, signing, my_index)

	var G_i kyber.Scalar
	G_i = curve.Scalar().Zero()
	for i = 1; i <= int64(N); i++ {

		G_i.Add(share[i-1], G_i)

	}
	return G_i
}

// func Verify_Share_sign(peer_number string, N int64, T int64) kyber.Scalar {
// 	var i int64
// 	path := "Received/Signing/" + peer_number + "/Shares/share"
// 	share := []kyber.Scalar{} // to store share
// 	for i = 0; i <= int64(N); i++ {
// 		share = append(share, curve.Scalar().Zero())
// 	}

// 	for i = 0; i <= N; i++ {
// 		if strconv.Itoa(int(i)) == peer_number { //To Add own share to G_i
// 			file, err := os.Open("vss/Signing/" + peer_number + "/Indivisual_Share" + peer_number + ".txt")
// 			if err != nil {
// 				panic(err)
// 			}
// 			val, e1 := encoding.ReadHexScalar(curve, file)
// 			if e1 != nil {
// 				panic(e1)
// 			}
// 			share[i] = val

// 		} else {
// 			file, err := os.Open(path + strconv.Itoa(int(i)) + ".txt")
// 			if err != nil {
// 				panic(err)
// 			}
// 			val, e1 := encoding.ReadHexScalar(curve, file)
// 			if e1 != nil {
// 				panic(e1)
// 			}
// 			share[i] = val
// 		}
// 	}
// 	// for i = 0; i <= N; i++ {
// 	// 	fmt.Println(i, "--->", share[i].String())
// 	// }
// 	verify_each_share_sign(peer_number, N, share, T,my_index)

// 	var R_i kyber.Scalar
// 	R_i = curve.Scalar().Zero()
// 	for i = 0; i <= int64(N); i++ {
// 		R_i.Add(share[i], R_i)
// 	}
// 	return R_i
// }

//verify_each_share
func verify_each_share(peer_number string, peer_count int64, share []kyber.Scalar, T int64, signing bool, my_index int) {
	var i int64
	fmt.Printf("Verifying shares Recieved By %s \n", peer_number)
	for i = 1; i <= peer_count; i++ {
		if i == int64(my_index+1) {
			continue
		}
		var j int64
		path := ""
		if signing {
			path = "Broadcast/" + strconv.Itoa(int(i)) + "/Signing/Alphas/"
		} else {
			path = "Received/" + strconv.Itoa(int(i)) + "/Keygen_alphas/"
		}
		alphas := []kyber.Point{} // to store alphas

		for j = 0; j < T; j++ {
			alphas = append(alphas, curve.Point().Null())
		}
		for j = 0; j < T; j++ {
			f1, e1 := os.Open(path + "alpha" + strconv.Itoa(int(j)) + ".txt")
			if e1 != nil {
				panic(e1)
			}
			alpha, _ := encoding.ReadHexPoint(curve, f1)
			alphas[j] = alpha
			f1.Close()
		}
		//fmt.Println(alphas)
		I, err := strconv.Atoi(peer_number)
		if err != nil {
			panic(err)
		}
		fmt.Println("INSDIE VerifyEachShare")
		fmt.Println("share->", share[i-1].String())
		var k int
		for k = 0; k < int(T); k++ {
			fmt.Println(alphas[k].String())
		}
		fmt.Println("")
		if !Verify_i(int64(I), share[i-1], T, alphas) {
			fmt.Printf("Peer %d shared wrong values mission aborting \n", i)
		} else {
			fmt.Printf("Shared Verified for Peer %d \n", i)
		}

	}
}

//verify_each_share
func verify_each_share_sign(peer_number string, peer_count int64, share []kyber.Scalar, T int64, my_index int) {
	var i int64
	fmt.Printf("Verifying shares Recieved By %s \n", peer_number)
	for i = 0; i <= peer_count; i++ {
		if i == int64(my_index) {
			continue
		}

		var j int64
		path := "Broadcast/" + strconv.Itoa(int(i)) + "/Signing/Alphas/"
		alphas := []kyber.Point{} // to store alphas

		for j = 0; j < T; j++ {
			alphas = append(alphas, curve.Point().Null())
		}
		for j = 0; j < T; j++ {
			f1, e1 := os.Open(path + "alpha" + strconv.Itoa(int(j)) + ".txt")
			if e1 != nil {
				panic(e1)
			}
			alpha, _ := encoding.ReadHexPoint(curve, f1)
			alphas[j] = alpha
			f1.Close()
		}
		//fmt.Println(alphas)
		I, err := strconv.Atoi(peer_number)
		if err != nil {
			panic(err)
		}
		fmt.Println("INSDIE VerifyEachShare Sign")
		fmt.Println("share->", share[i].String())
		var k int
		for k = 0; k < int(T); k++ {
			fmt.Println(alphas[k].String())
		}
		fmt.Println("")
		if !Verify_i(int64(I), share[i], T, alphas) {
			fmt.Printf("Peer %d shared wrong values mission aborting \n", i)
		} else {
			fmt.Printf("Shared Verified for Peer %d \n", i)
		}
	}
}

func Recieve_Share_sign(peer_number string, Peer_Count int64, my_index int) {
	var i int64

	for i = 1; i <= Peer_Count; i++ {
		if i == int64(my_index+1) {
			continue
		}
		// path := "Broadcast/" + fmt.Sprint(i) + "/Signing/Shares/shareTo" + peer_number + ".txt"
		// file, _ := os.Open(path)
		// share, _ := encoding.ReadHexScalar(curve, file)
		// path2 := "Received/Signing/" + peer_number + "/Shares/share" + fmt.Sprint(i) + ".txt"

		// os.MkdirAll("Received/Signing/"+peer_number+"/Shares/", os.ModePerm)

		// file2, _ := os.Create(path2)

		// encoding.WriteHexScalar(curve, file2, share)
		path := "Broadcast/" + fmt.Sprint(i) + "/Signing/Shares/shareTo" + peer_number + ".txt"
		file, _ := os.ReadFile(path)
		share := string(file)

		fmt.Println("RECIEVED:", share)

		path2 := "Received/Signing/" + peer_number + "/Shares/share" + fmt.Sprint(i) + ".txt"
		os.MkdirAll("Received/Signing/"+peer_number+"/Shares/", os.ModePerm)
		file2, _ := os.OpenFile(path2, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		_, _ = fmt.Fprint(file2, share)

	}
}

// func Verify_final_sign(V kyber.Scalar, U kyber.Point, message string, GK kyber.Point) bool {
// 	//message, U , V public key
// 	//V is sum of all V_i's
// 	//U is sum of all U_i's
// 	//GK is sum of all alpha[0] (group key)

// 	t1 := curve.Point().Mul(V, g)
// 	// h := Hash(message + U.String())
// 	Hashing_message := message + U.String()
// 	h, _ := hash_sign([]byte(Hashing_message))

// 	var H1 kyber.Scalar
// 	H1 = curve.Scalar().Zero()
// 	H1.SetBytes(h)

// 	t2 := curve.Point().Mul(H1, GK)
// 	t2 = t2.Add(t2, U)

// 	fmt.Println(t1)
// 	fmt.Println(t2)

// 	if t1.Equal(t2) {
// 		return true
// 	} else {
// 		return false
// 	}

// }

// func Verify_sign_share(V_i kyber.Scalar, U kyber.Point, U_i kyber.Point, message string, X_i kyber.Point) bool {
// 	//message, U , V public key
// 	//V is sum of all V_i's
// 	//U is sum of all U_i's
// 	//GK is sum of all alpha[0] (group key)

// 	t1 := curve.Point().Mul(V_i, g)
// 	// h := Hash(message + U.String())
// 	Hashing_message := message + U.String()
// 	h, _ := hash_sign([]byte(Hashing_message))

// 	var H1 kyber.Scalar
// 	H1 = curve.Scalar().Zero()
// 	H1.SetBytes(h)

// 	t2 := curve.Point().Mul(H1, X_i)
// 	t2 = t2.Add(t2, U_i)

// 	fmt.Println(t1)
// 	fmt.Println(t2)

// 	if t1.Equal(t2) {
// 		return true
// 	} else {
// 		return false
// 	}

// }
