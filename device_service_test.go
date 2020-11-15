package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateDeviceWithWrongInterval(t *testing.T) {
	underTest := DeviceService{dao: &inMemoryDeviceDAO{}}

	device := Device{Name: "name", Interval: 0, Value: 1}

	status, _, err := underTest.CreateDevice(device)

	assert.Equal(t, validationError, status)
	assert.Error(t, err, "Device with Interval less than 1 isn't valid.")
}

func TestCreateDeviceWithoutName(t *testing.T) {
	underTest := DeviceService{dao: &inMemoryDeviceDAO{}}

	device := Device{Interval: 1, Value: 1}

	status, _, err := underTest.CreateDevice(device)

	assert.Equal(t, validationError, status)
	assert.Error(t, err, "Device without name isn't valid")
}

func TestCreateDeviceWithCorrectData(t *testing.T) {
	underTest := DeviceService{dao: &inMemoryDeviceDAO{}}

	device := Device{Name: "name", Interval: 1, Value: 1}

	status, _, err := underTest.CreateDevice(device)

	assert.Equal(t, deviceCreated, status)
	assert.NoError(t, err, "Device with correct data should be created")
}

func TestCreateValidDeviceButErrorWhenSavingByDAO(t *testing.T) {
	underTest := DeviceService{dao: &failingDeviceDAO{}}

	device := Device{Name: "name", Interval: 1, Value: 1}

	status, _, err := underTest.CreateDevice(device)

	assert.Equal(t, savingError, status)
	assert.Error(t, err, "The error should be return when DAO fails")
}
