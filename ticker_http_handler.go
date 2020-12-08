package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type tickerHTTPHandler struct {
	ts *TickerService
	mw *MeasurementsWriter
}

func newTickerHTTPHandler(ds *DeviceService) tickerHTTPHandler {
	measurements := make(chan Measurement, 10)

	ts := TickerService{ds: ds, measurements: measurements}
	ds.AddObserver(&ts)

	mw := MeasurementsWriter{measurements: measurements}

	return tickerHTTPHandler{ts: &ts, mw: &mw}
}

func (h *tickerHTTPHandler) Start(w http.ResponseWriter, _ *http.Request) {
	if err := h.ts.Start(); err != nil {
		log.Print(err)
		http.Error(w, "failed to start ticker", http.StatusInternalServerError)
	}

	h.mw.Start()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("ticker started"); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *tickerHTTPHandler) Stop(w http.ResponseWriter, _ *http.Request) {
	h.ts.Stop()
	h.mw.Stop()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("ticker stopped"); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
