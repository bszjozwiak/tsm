package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type deviceHTTPHandler struct {
	service DeviceService
}

func (h *deviceHTTPHandler) createDevice(w http.ResponseWriter, r *http.Request) {
	var requestDevice Device

	if err := json.NewDecoder(r.Body).Decode(&requestDevice); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdDevice, err := h.service.CreateDevice(requestDevice)
	if err != nil {
		log.Print(err)
		switch err.Error() {
		case validationEmptyDeviceNameErr, validationWrongIntervalErr:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case daoSaveErr:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			//3BJO_TODO: Programming error. panic("not handled status") or log.Fatal("not handled status")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdDevice); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *deviceHTTPHandler) getByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Print(err)
		http.Error(w, "device id must be a number", http.StatusBadRequest)
		return
	}

	device, err := h.service.GetByID(id)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if device == nil {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(fmt.Sprintf("device with id %v not found", id)); err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(device); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
