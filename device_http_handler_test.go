package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
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
	underTest := deviceHTTPHandler{service: DeviceService{dao: &inMemoryDeviceDAO{}}}

	underTest.createDevice(res, req)

	result := res.Result()
	responseDevice := getResponseContent(result)

	assert.Equal(t, http.StatusCreated, result.StatusCode)
	assert.EqualValues(t, `{"id":0,"name":"test name2","interval":1,"value":0}`, responseDevice, "")
}

func TestCreateInvalidDevice(t *testing.T) {
	body := strings.NewReader(`{"name":"test name2","interval":-1}`)
	req := httptest.NewRequest(http.MethodPost, "/devices", body)
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: DeviceService{dao: &inMemoryDeviceDAO{}}}

	underTest.createDevice(res, req)

	result := res.Result()

	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func TestCreateValidDeviceButErrorWhenSaveInDatabase(t *testing.T) {
	body := strings.NewReader(`{"name":"test name2","interval":1}`)
	req := httptest.NewRequest(http.MethodPost, "/devices", body)
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: DeviceService{dao: &failingDeviceDAO{}}}

	underTest.createDevice(res, req)

	result := res.Result()

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

func TestGetByIDStringPassedInsteadOfNumber(t *testing.T) {
	req := createGetDeviceRequest("string")
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: DeviceService{dao: &inMemoryDeviceDAO{}}}

	underTest.getByID(res, req)

	result := res.Result()

	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func TestGetByIDDatabaseError(t *testing.T) {
	req := createGetDeviceRequest("0")
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: DeviceService{dao: &failingDeviceDAO{}}}

	underTest.getByID(res, req)

	result := res.Result()

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

func TestGetByIdDeviceNotFound(t *testing.T) {
	req := createGetDeviceRequest("0")
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: DeviceService{dao: &inMemoryDeviceDAO{}}}

	underTest.getByID(res, req)

	result := res.Result()

	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

func TestGetByIDDeviceExists(t *testing.T) {
	req := createGetDeviceRequest("0")
	res := httptest.NewRecorder()
	device := Device{Id: 0, Name: "device", Interval: 1, Value: 1}
	dao := inMemoryDeviceDAO{}
	_, _ = dao.Save(device)
	underTest := deviceHTTPHandler{service: DeviceService{dao: &dao}}

	underTest.getByID(res, req)

	result := res.Result()
	responseDevice := getResponseContent(result)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.EqualValues(t, `{"id":0,"name":"device","interval":1,"value":1}`, responseDevice, "")
}

func createGetDeviceRequest(id string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": id,
	})
	return req
}

func getResponseContent(result *http.Response) string {
	responseBody, _ := ioutil.ReadAll(result.Body)
	_ = result.Body.Close()
	return strings.Trim(string(responseBody), "\n")
}

type failingDeviceDAO struct {
}

func (db *failingDeviceDAO) Save(device Device) (Device, error) {
	return device, errors.New("mock error - fail to create device")
}

func (db *failingDeviceDAO) GetByID(id int) (*Device, error) {
	return nil, errors.New(fmt.Sprintf("mock error - fail to get device with id %v", id))
}
