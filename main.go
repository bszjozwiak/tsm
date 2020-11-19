package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	deviceService := deviceHTTPHandler{service: DeviceService{dao: &inMemoryDeviceDAO{}}}
	myRouter.HandleFunc("/devices", deviceService.createDevice).Methods("POST")
	log.Fatal(http.ListenAndServe(getAddr(), myRouter))
}

func getAddr() string {
	port := os.Getenv("TSM_PORT")

	if port != "" {
		return ":" + port
	}

	return ":8000"
}
