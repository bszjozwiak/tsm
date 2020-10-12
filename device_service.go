package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Device struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Interval int     `json:"interval"`
	Value    float32 `json:"value"`
}

type DeviceService struct {
	Database
}

func (service *DeviceService) createDevice(w http.ResponseWriter, r *http.Request) {
	var device Device

	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := service.validate(device); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	device, err := service.saveDevice(device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(device); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (service *DeviceService) validate(device Device) error {
	if len(device.Name) == 0 {
		return errors.New("device name can't be empty")
	}

	if device.Interval <= 0 {
		return errors.New("interval has to be greater than 0")
	}

	return nil
}
