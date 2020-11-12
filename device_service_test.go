package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateDeviceWithWrongInterval(t *testing.T) {
	underTest := deviceService{}

	device := Device{Name: "name", Interval: 0, Value: 1}

	err := underTest.validate(device)

	assert.Error(t, err, "Device with Interval less than 1 isn't valid.")
}

func TestValidateDeviceWithoutName(t *testing.T) {
	underTest := deviceService{}

	device := Device{Interval: 1, Value: 1}

	err := underTest.validate(device)

	assert.Error(t, err, "Device without name isn't valid")
}

func TestValidateDeviceWithCorrectData(t *testing.T) {
	underTest := deviceService{}

	device := Device{Name: "name", Interval: 1, Value: 1}

	err := underTest.validate(device)

	assert.NoError(t, err, "Device with correct data should be validate positive")
}
