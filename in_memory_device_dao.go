package main

import (
	"context"
	"errors"
	"sync"
)

type inMemoryDeviceDAO struct {
	mu      sync.Mutex
	devices []Device
}

func (db *inMemoryDeviceDAO) Save(_ context.Context, device Device) (Device, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	device.Id = len(db.devices)
	db.devices = append(db.devices, device)

	return device, nil
}

func (db *inMemoryDeviceDAO) GetByID(_ context.Context, id int) (*Device, error) {
	if len(db.devices) > id {
		device := db.devices[id]
		return &device, nil
	}

	return nil, nil
}

func (db *inMemoryDeviceDAO) GetAll(_ context.Context, limit int, page int) ([]Device, error) {
	if limit < 0 {
		return nil, errors.New("limit can't be negative")
	}

	if limit == 0 {
		return append([]Device(nil), db.devices...), nil
	}

	start := limit * page
	if start >= len(db.devices) {
		return []Device{}, nil
	}

	end := start + limit
	if end > len(db.devices) {
		end = len(db.devices)
	}

	return append([]Device(nil), db.devices[start:end]...), nil
}
