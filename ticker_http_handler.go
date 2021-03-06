package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type tickerHTTPHandler struct {
	ts *TickerService
}

func newTickerHTTPHandler(ds *DeviceService, publisher measurementPublisher) tickerHTTPHandler {
	ts := TickerService{ds: ds, publisher: publisher, tf: time.Tick}
	ds.AddObserver(&ts)

	return tickerHTTPHandler{ts: &ts}
}

func (h *tickerHTTPHandler) Start(w http.ResponseWriter, r *http.Request) {
	if err := h.ts.Start(r.Context()); err != nil {
		log.Print(err)
		http.Error(w, "failed to start ticker", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("ticker started"); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *tickerHTTPHandler) Stop(w http.ResponseWriter, _ *http.Request) {
	h.ts.Stop()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("ticker stopped"); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
