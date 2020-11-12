package main

import (
	"errors"
	"log"
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

type deviceService struct {
	dao deviceDAO
}

func (s *deviceService) createDevice(device Device) (Device, error) {
	if err := s.validate(device); err != nil {
		return device, err
	}

	savedDevice, err := s.dao.Save(device)
	if err != nil {
		log.Print(err)
		return device, errors.New("fail to save device")
	}

	return savedDevice, nil
}

func (s *deviceService) validate(device Device) error {
	if device.Name == "" {
		return errors.New("device name can't be empty")
	}

	if device.Interval <= 0 {
		return errors.New("interval has to be greater than 0")
	}

	return nil
}
