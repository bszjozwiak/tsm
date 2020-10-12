package main

type Database interface {
	saveDevice(device Device) (Device, error)
}

type InMemoryDatabase struct {
	devices []Device
}

func (db *InMemoryDatabase) saveDevice(device Device) (Device, error) {
	device.Id = len(db.devices)
	db.devices = append(db.devices, device)

	return device, nil
}
