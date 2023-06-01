package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

var router = mux.NewRouter().StrictSlash(true)

func main() {
	fmt.Println("Listening on IP 172.0.0.1:8080")

	ServiceHandle()

	handle := cors.AllowAll().Handler(router)
	err := http.ListenAndServe(":8080", handle)
	if err != nil {
		fmt.Println("error", err)
	}
}

func ServiceHandle() {
	router.HandleFunc("/services/insert", InsertService).Methods("POST")
}
