package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"gopkg.in/dedis/kyber.v2/util/encoding"
)

var tmpl *template.Template

type Keygen_Data struct {
	IP_addresses string
	Threshold_T  int
	GroupKey     string
}

type MyInfoStruct struct {
	VaultID string
	MyIp    string
	//Make Every Property start with Capital letter
}

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
	//For accessing the assets folder in index.html file . we need to specify that ( IN HTTP ROUTER)
	// fs := http.FileServer(http.Dir("assets"))
	// http.Handle("/assets/", http.StripPrefix("/assets", fs))

	//SAME AS ABOVE FOR MUX ROUTER ( SHIFTED TO MAIN.GO)
	// fs := http.FileServer(http.Dir("assets"))
	// r.PathPrefix("/assets").Handler(http.StripPrefix("/assets", fs))
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to My world %s", r.URL.Path)

}

func Dummy_api() string {
	return "HELLO THIS IS YOUR GUILTY CONSCIENCE"
}

func Index(w http.ResponseWriter, r *http.Request) {

	tmpl.ExecuteTemplate(w, "index.html", nil)

}
func CSS2(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "home.html", nil)
}
func DisplayForm(w http.ResponseWriter, r *http.Request) {
	//P2p_func()
	this_vault = r.FormValue("vaultID")
	fmt.Println("vault:", this_vault)
	recieved_data := MyInfoStruct{
		VaultID: r.FormValue("vaultID"),
		MyIp:    p2p.Host_ip,
	}

	// recieved_data := Data2{
	// 	Item1: p2p.Host_ip,
	// 	Item2: "hello",
	// 	Item3: Dummy_api(),
	// }
	tmpl.ExecuteTemplate(w, "dealer.html", struct {
		Success bool
		Mydata  MyInfoStruct
	}{true, recieved_data})
}

func DisplayData(w http.ResponseWriter, r *http.Request) {

	tempT, _ := strconv.Atoi(r.FormValue("T"))
	IPs := r.FormValue("ip")

	Threshold = tempT
	gen_keyshares(IPs)
	peer_number := fmt.Sprint(my_index + 1)
	path := "Received/" + peer_number
	file, _ := os.Open(path + "/GroupKey.txt")
	GK, _ := encoding.ReadHexPoint(curve, file)
	// fmt.Println("GROUP KEY:", GK.String())
	recieved_data := Keygen_Data{
		IP_addresses: IPs,
		Threshold_T:  tempT,
		GroupKey:     GK.String(),
	}
	tmpl.ExecuteTemplate(w, "signing.html", struct {
		Success bool
		Mydata  Keygen_Data
	}{true, recieved_data})
	// tmpl.ExecuteTemplate(w, "display.html", struct {
	// 	Mydata Data2
	// }{recieved_data})
}

// func Init_vault(){
func Sign_Message(w http.ResponseWriter, r *http.Request) {
	Message := r.FormValue("message")
	peer_number := fmt.Sprint(my_index + 1)
	verified_value := Signing(peer_number, Message)
	fmt.Println("MYDATA DATA:", verified_value.Verify_bool)
	tmpl.ExecuteTemplate(w, "final.html", struct {
		P2p_send P2P
		Mydata   Verify_sign
	}{p2p, verified_value})

}

func DisplayNonDealer(w http.ResponseWriter, r *http.Request) {
	execute_send = 0
	this_vault = r.FormValue("vaultID2")
	fmt.Println("VALULT ND:", this_vault)

	tmpl.ExecuteTemplate(w, "nondealer.html", struct {
		P2p_send P2P
	}{p2p})
}

// }
