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
	service *DeviceService
}

func (h *deviceHTTPHandler) createDevice(w http.ResponseWriter, r *http.Request) {
	var requestDevice Device

	if err := json.NewDecoder(r.Body).Decode(&requestDevice); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdDevice, err := h.service.CreateDevice(r.Context(), requestDevice)
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
	id := mux.Vars(r)["id"]

	device, err := h.service.GetByID(r.Context(), id)
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

func (h *deviceHTTPHandler) getAll(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	limit, err := h.getValueOrDefault(params.Get("limit"), 100)
	if err != nil || limit < 0 {
		log.Print(err)
		http.Error(w, "limit must be a number greater or equal to 0", http.StatusBadRequest)
		return
	}

	page, err := h.getValueOrDefault(params.Get("page"), 0)
	if err != nil || page < 0 {
		log.Print(err)
		http.Error(w, "page must be a number greater or equal to 0", http.StatusBadRequest)
		return
	}

	devices, err := h.service.GetAll(r.Context(), limit, page)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(devices); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *deviceHTTPHandler) getValueOrDefault(param string, defVal int) (int, error) {
	if param == "" {
		return defVal, nil
	}

	return strconv.Atoi(param)
}
