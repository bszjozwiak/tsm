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

	status, createdDevice, err := h.service.CreateDevice(requestDevice)
	switch status {
	case deviceCreated:
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(createdDevice); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case validationError:
		http.Error(w, getErrorMessage(err), http.StatusBadRequest)
	case savingError:
		http.Error(w, getErrorMessage(err), http.StatusInternalServerError)
	default:
		//3BJO_TODO: Programming error. panic("not handled status") or log.Fatal("not handled status")
	}
}

func getErrorMessage(err error) string {
	if err != nil {
		return err.Error()
	}

	return ""
}
