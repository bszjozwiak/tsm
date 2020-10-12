package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateDeviceWithWrongInterval(t *testing.T) {
	underTest := DeviceService{}

	device := Device{Name: "name", Interval: 0, Value: 1}

	err := underTest.validate(device)

	assert.NotNil(t, err, "Device with Interval less than 1 isn't valid.")
}

func TestValidateDeviceWithoutName(t *testing.T) {
	underTest := DeviceService{}

	device := Device{Interval: 1, Value: 1}

	err := underTest.validate(device)

	assert.NotNil(t, err, "Device without name isn't valid")
}

func TestValidateDeviceWithCorrectData(t *testing.T) {
	underTest := DeviceService{}

	device := Device{Name: "name", Interval: 1, Value: 1}

	err := underTest.validate(device)

	assert.Nil(t, err, "Device with correct data should be validate positive")
}

func TestCreateValidDeviceWithNoDatabaseErrors(t *testing.T) {
	body := strings.NewReader(`{"name":"test name2","interval":1}`)
	req := httptest.NewRequest(http.MethodPost, "/devices", body)
	res := httptest.NewRecorder()
	underTest := DeviceService{Database: &InMemoryDatabase{}}

	underTest.createDevice(res, req)

	result := res.Result()
	responseBody, _ := ioutil.ReadAll(result.Body)
	_ = result.Body.Close()

	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.EqualValues(t, `{"id":0,"name":"test name2","interval":1,"value":0}`, strings.Trim(string(responseBody), "\n"), "")
}

func TestCreateInvalidDevice(t *testing.T) {
	body := strings.NewReader(`{"name":"test name2","interval":-1}`)
	req := httptest.NewRequest(http.MethodPost, "/devices", body)
	res := httptest.NewRecorder()
	underTest := DeviceService{Database: &InMemoryDatabase{}}

	underTest.createDevice(res, req)

	result := res.Result()

	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func TestCreateValidDeviceButErrorWhenSaveInDatabase(t *testing.T) {
	body := strings.NewReader(`{"name":"test name2","interval":1}`)
	req := httptest.NewRequest(http.MethodPost, "/devices", body)
	res := httptest.NewRecorder()
	underTest := DeviceService{Database: &FailingDatabase{}}

	underTest.createDevice(res, req)

	result := res.Result()

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

type FailingDatabase struct {
}

func (db *FailingDatabase) saveDevice(device Device) (Device, error) {
	return device, errors.New("mock error")
}
