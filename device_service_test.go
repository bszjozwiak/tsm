package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateDeviceWithWrongInterval(t *testing.T) {
	underTest := DeviceService{dao: &inMemoryDeviceDAO{}}

	device := Device{Name: "name", Interval: 0, Value: 1}

	_, err := underTest.CreateDevice(context.Background(), device)

	assert.Error(t, err, "Device with Interval less than 1 isn't valid.")
}

func TestCreateDeviceWithoutName(t *testing.T) {
	underTest := DeviceService{dao: &inMemoryDeviceDAO{}}

	device := Device{Interval: 1, Value: 1}

	_, err := underTest.CreateDevice(context.Background(), device)

	assert.Error(t, err, "Device without name isn't valid")
}

func TestCreateDeviceWithCorrectData(t *testing.T) {
	underTest := DeviceService{dao: &inMemoryDeviceDAO{}}

	device := Device{Name: "name", Interval: 1, Value: 1}

	_, err := underTest.CreateDevice(context.Background(), device)

	assert.NoError(t, err, "Device with correct data should be created")
}

func TestCreateValidDeviceButErrorWhenSavingByDAO(t *testing.T) {
	underTest := DeviceService{dao: &failingDeviceDAO{}}

	device := Device{Name: "name", Interval: 1, Value: 1}

	_, err := underTest.CreateDevice(context.Background(), device)

	assert.Error(t, err, "The error should be return when DAO fails")
	assert.EqualError(t, err, daoSaveErr)
}

func TestGetByIDErrorInDAO(t *testing.T) {
	underTest := DeviceService{dao: &failingDeviceDAO{}}

	_, err := underTest.GetByID(context.Background(), 1)

	assert.Error(t, err, "The error should be return when DAO fails")
	assert.EqualError(t, err, daoGetErr)
}

func TestGetAllErrorInDAO(t *testing.T) {
	underTest := DeviceService{dao: &failingDeviceDAO{}}

	_, err := underTest.GetAll(context.Background(), 0, 0)

	assert.Error(t, err, "The error should be return when DAO fails")
	assert.EqualError(t, err, daoGetAllErr)
}

func TestDeviceService_CreateDevice_NotifyAllObservers(t *testing.T) {
	firstObserver := deviceServiceObserver{}
	secondObserver := deviceServiceObserver{}

	underTest := DeviceService{dao: &inMemoryDeviceDAO{}}
	underTest.AddObserver(&firstObserver)
	underTest.AddObserver(&secondObserver)

	device := Device{Name: "name", Interval: 1, Value: 1}

	_, _ = underTest.CreateDevice(context.Background(), device)

	assert.True(t, firstObserver.notified)
	assert.True(t, secondObserver.notified)
}

type deviceServiceObserver struct {
	notified bool
}

func (o *deviceServiceObserver) NotifyDeviceCreated(_ Device) {
	o.notified = true
}
