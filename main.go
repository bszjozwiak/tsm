package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	deviceService := deviceHTTPHandler{service: deviceService{dao: &inMemoryDeviceDAO{}}}
	myRouter.HandleFunc("/devices", deviceService.createDevice).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}
