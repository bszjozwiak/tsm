package main

import (
	"encoding/json"
	"net/http"
)

type deviceHTTPHandler struct {
	service deviceService
}

func (h *deviceHTTPHandler) createDevice(w http.ResponseWriter, r *http.Request) {
	var requestDevice Device

	if err := json.NewDecoder(r.Body).Decode(&requestDevice); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdDevice, err := h.service.createDevice(requestDevice)
	if err != nil {
		code := http.StatusBadRequest
		if err.Error() == "fail to save device" {
			code = http.StatusInternalServerError
		}

		http.Error(w, err.Error(), code)
		return
	}

	if err := json.NewEncoder(w).Encode(createdDevice); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
