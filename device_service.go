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
	daoGetAllErr                 = "failed to get all devices"
)

type Device struct {
	Id       int     `json:"id" bson:"id"`
	Name     string  `json:"name" bson:"name"`
	Interval int     `json:"interval" bson:"interval"`
	Value    float64 `json:"value" bson:"value"`
}

type deviceDAO interface {
	Save(device Device) (Device, error)
	GetByID(id int) (*Device, error)
	GetAll(limit int, page int) ([]Device, error)
}

type DeviceCreateObserver interface {
	NotifyDeviceCreated(device Device)
}

type DeviceService struct {
	dao       deviceDAO
	observers []DeviceCreateObserver
}

func (s *DeviceService) AddObserver(observer DeviceCreateObserver) {
	s.observers = append(s.observers, observer)
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

	for _, observer := range s.observers {
		observer.NotifyDeviceCreated(savedDevice)
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

func (s *DeviceService) GetAll(limit int, page int) ([]Device, error) {
	devices, err := s.dao.GetAll(limit, page)
	if err != nil {
		log.Print(err)
		return nil, errors.New(daoGetAllErr)
	}

	return devices, nil
}
