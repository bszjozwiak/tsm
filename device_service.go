package main

import (
	"errors"
	"log"
)

const (
	validationEmptyDeviceNameErr = "device name can't be empty"
	validationWrongIntervalErr   = "interval has to be greater than 0"
	daoSaveErr                   = "failed to save device"
	daoGetErr                    = "failed to get device"
)

type Device struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Interval int     `json:"interval"`
	Value    float32 `json:"value"`
}

type deviceDAO interface {
	Save(device Device) (Device, error)
	GetByID(id int) (*Device, error)
}

type DeviceService struct {
	dao deviceDAO
}

func (s *DeviceService) CreateDevice(device Device) (Device, error) {
	if err := s.validate(device); err != nil {
		return device, err
	}

	savedDevice, err := s.dao.Save(device)
	if err != nil {
		log.Print(err)
		return device, errors.New(daoSaveErr)
	}

	return savedDevice, nil
}

func (s *DeviceService) validate(device Device) error {
	if device.Name == "" {
		return errors.New(validationEmptyDeviceNameErr)
	}

	if device.Interval <= 0 {
		return errors.New(validationWrongIntervalErr)
	}

	return nil
}

func (s *DeviceService) GetByID(id int) (*Device, error) {
	device, err := s.dao.GetByID(id)
	if err != nil {
		log.Print(err)
		return nil, errors.New(daoGetErr)
	}

	return device, nil
}
