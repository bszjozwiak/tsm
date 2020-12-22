package main

import (
	"github.com/gorilla/mux"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"log"
	"net/http"
	"os"
)

func main() {
	measurements := make(chan Measurement, 10)
	client := influxdb2.NewClient(os.Getenv("TSM_INFLUX_URL"), os.Getenv("TSM_INFLUX_TOKEN"))
	mw := MeasurementsWriter{measurements: measurements, writeAPI: client.WriteAPIBlocking("tsm", "mydb")}
	mw.AsyncStart()

	myRouter := mux.NewRouter().StrictSlash(true)
	deviceService := DeviceService{dao: &inMemoryDeviceDAO{}}

	deviceHandler := deviceHTTPHandler{service: &deviceService}
	myRouter.HandleFunc("/devices", deviceHandler.createDevice).Methods(http.MethodPost)
	myRouter.HandleFunc("/devices/{id}", deviceHandler.getByID).Methods(http.MethodGet)
	myRouter.HandleFunc("/devices", deviceHandler.getAll).Methods(http.MethodGet)

	tickerHandler := newTickerHTTPHandler(&deviceService, measurements)
	myRouter.HandleFunc("/start", tickerHandler.Start).Methods(http.MethodPost)
	myRouter.HandleFunc("/stop", tickerHandler.Stop).Methods(http.MethodPost)

	log.Printf("tsm started and listening on %v", getAddr())
	log.Fatal(http.ListenAndServe(getAddr(), myRouter))
}

func getAddr() string {
	port := os.Getenv("TSM_PORT")

	if port != "" {
		return ":" + port
	}

	return ":8000"
}
