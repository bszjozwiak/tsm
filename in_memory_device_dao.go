package main

import "sync"

type inMemoryDeviceDAO struct {
	mu      sync.Mutex
	devices []Device
}

func (db *inMemoryDeviceDAO) Save(device Device) (Device, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	device.Id = len(db.devices)
	db.devices = append(db.devices, device)

	return device, nil
}
