package main

import (
	"errors"
	"log"
)

const (
	deviceCreated = iota
	validationError
	savingError
)

type Device struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Interval int     `json:"interval"`
	Value    float32 `json:"value"`
}

type deviceDAO interface {
	Save(device Device) (Device, error)
}

type DeviceService struct {
	dao deviceDAO
}

func (s *DeviceService) CreateDevice(device Device) (int, Device, error) {
	if err := s.validate(device); err != nil {
		return validationError, device, err
	}

	savedDevice, err := s.dao.Save(device)
	if err != nil {
		log.Print(err)
		return savingError, device, errors.New("fail to save device")
	}

	return deviceCreated, savedDevice, nil
}

func (s *DeviceService) validate(device Device) error {
	if device.Name == "" {
		return errors.New("device name can't be empty")
	}

	if device.Interval <= 0 {
		return errors.New("interval has to be greater than 0")
	}

	return nil
}
