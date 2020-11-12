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

func TestCreateValidDeviceWithNoDatabaseErrors(t *testing.T) {
	body := strings.NewReader(`{"name":"test name2","interval":1}`)
	req := httptest.NewRequest(http.MethodPost, "/devices", body)
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: deviceService{dao: &inMemoryDeviceDAO{}}}

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
	underTest := deviceHTTPHandler{service: deviceService{dao: &inMemoryDeviceDAO{}}}

	underTest.createDevice(res, req)

	result := res.Result()

	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func TestCreateValidDeviceButErrorWhenSaveInDatabase(t *testing.T) {
	body := strings.NewReader(`{"name":"test name2","interval":1}`)
	req := httptest.NewRequest(http.MethodPost, "/devices", body)
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: deviceService{dao: &failingDeviceDAO{}}}

	underTest.createDevice(res, req)

	result := res.Result()

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

type failingDeviceDAO struct {
}

func (db *failingDeviceDAO) Save(device Device) (Device, error) {
	return device, errors.New("mock error")
}
