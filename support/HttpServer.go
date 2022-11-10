package main

// import (
// 	// "fmt"
// 	"net/http"

// 	"github.com/gorilla/mux"
// )

// var r *mux.Router

// func main2() {
// 	//http routing
// 	// http.HandleFunc("/", handlerFunc)
// 	// http.HandleFunc("/css1", CSS1)
// 	// http.HandleFunc("/css2", CSS2)
// 	// http.ListenAndServe(":8080", nil)
// 	//WE have advanced router MUX
// 	r = mux.NewRouter()

// 	fs := http.FileServer(http.Dir("assets"))
// 	r.PathPrefix("/assets").Handler(http.StripPrefix("/assets", fs))
// 	r.HandleFunc("/", handlerFunc)
// 	// r.HandleFunc("/css1", CSS1)
// 	r.HandleFunc("/css2", CSS2)
// 	r.HandleFunc("/form", DisplayForm)
// 	r.HandleFunc("/display", DisplayData)
// 	http.ListenAndServe(":8080", r)

// }
