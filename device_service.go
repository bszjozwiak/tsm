package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Interval int                `json:"interval" bson:"interval"`
	Value    float64            `json:"value" bson:"value"`
}

type deviceDAO interface {
	Save(ctx context.Context, device Device) (Device, error)
	GetByID(ctx context.Context, id string) (*Device, error)
	GetAll(ctx context.Context, limit int, page int) ([]Device, error)
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

func (s *DeviceService) CreateDevice(ctx context.Context, device Device) (Device, error) {
	if err := s.validate(device); err != nil {
		return device, err
	}

	savedDevice, err := s.dao.Save(ctx, device)
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

func (s *DeviceService) GetByID(ctx context.Context, id string) (*Device, error) {
	device, err := s.dao.GetByID(ctx, id)
	if err != nil {
		log.Print(err)
		return nil, errors.New(daoGetErr)
	}

	return device, nil
}

func (s *DeviceService) GetAll(ctx context.Context, limit int, page int) ([]Device, error) {
	devices, err := s.dao.GetAll(ctx, limit, page)
	if err != nil {
		log.Print(err)
		return nil, errors.New(daoGetAllErr)
	}

	return devices, nil
}
