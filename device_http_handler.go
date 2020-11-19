package main

import (
	"encoding/json"
	"net/http"
)

type deviceHTTPHandler struct {
	service DeviceService
}

func (h *deviceHTTPHandler) createDevice(w http.ResponseWriter, r *http.Request) {
	var requestDevice Device

	if err := json.NewDecoder(r.Body).Decode(&requestDevice); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdDevice, err := h.service.CreateDevice(requestDevice)
	if err != nil {
		switch err.Error() {
		case validationEmptyDeviceNameErr, validationWrongIntervalErr:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case savingErr:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			//3BJO_TODO: Programming error. panic("not handled status") or log.Fatal("not handled status")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdDevice); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
