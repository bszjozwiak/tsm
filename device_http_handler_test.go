package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	assert.Equal(t, http.StatusCreated, res.Code)
	require.JSONEq(t, `{"id":0,"name":"test name2","interval":1,"value":0}`, res.Body.String())
}

func TestCreateInvalidDevice(t *testing.T) {
	body := strings.NewReader(`{"name":"test name2","interval":-1}`)
	req := httptest.NewRequest(http.MethodPost, "/devices", body)
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: DeviceService{dao: &inMemoryDeviceDAO{}}}

	underTest.createDevice(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestCreateValidDeviceButErrorWhenSaveInDatabase(t *testing.T) {
	body := strings.NewReader(`{"name":"test name2","interval":1}`)
	req := httptest.NewRequest(http.MethodPost, "/devices", body)
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: DeviceService{dao: &failingDeviceDAO{}}}

	underTest.createDevice(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestGetByIDStringPassedInsteadOfNumber(t *testing.T) {
	req := createGetDeviceRequest("string")
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: DeviceService{dao: &inMemoryDeviceDAO{}}}

	underTest.getByID(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestGetByIDDatabaseError(t *testing.T) {
	req := createGetDeviceRequest("0")
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: DeviceService{dao: &failingDeviceDAO{}}}

	underTest.getByID(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestGetByIdDeviceNotFound(t *testing.T) {
	req := createGetDeviceRequest("0")
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: DeviceService{dao: &inMemoryDeviceDAO{}}}

	underTest.getByID(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestGetByIDDeviceExists(t *testing.T) {
	req := createGetDeviceRequest("0")
	res := httptest.NewRecorder()
	device := Device{Id: 0, Name: "device", Interval: 1, Value: 1}
	dao := inMemoryDeviceDAO{}
	_, _ = dao.Save(device)
	underTest := deviceHTTPHandler{service: DeviceService{dao: &dao}}

	underTest.getByID(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	require.JSONEq(t, `{"id":0,"name":"device","interval":1,"value":1}`, res.Body.String())
}

func createGetDeviceRequest(id string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": id,
	})
	return req
}

func TestGetAllBadRequestParameters(t *testing.T) {
	testCases := map[string][]string{
		"page":  {"-1", "string"},
		"limit": {"-1", "string"},
	}

	underTest := deviceHTTPHandler{service: DeviceService{dao: &inMemoryDeviceDAO{}}}

	for param, values := range testCases {
		for _, value := range values {
			req := httptest.NewRequest(http.MethodGet, "/devices", nil)
			query := req.URL.Query()
			query.Add(param, value)
			req.URL.RawQuery = query.Encode()
			res := httptest.NewRecorder()

			underTest.getAll(res, req)

			assert.Equal(t, http.StatusBadRequest, res.Code)
		}
	}
}

func TestGetAllDatabaseError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	res := httptest.NewRecorder()

	underTest := deviceHTTPHandler{service: DeviceService{dao: &failingDeviceDAO{}}}

	underTest.getAll(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestGetAllNoDevicesReturnEmptyList(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	res := httptest.NewRecorder()

	underTest := deviceHTTPHandler{service: DeviceService{dao: &inMemoryDeviceDAO{}}}

	underTest.getAll(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	require.JSONEq(t, "[]", res.Body.String())
}

func TestGetAllReturnRequestedDevices(t *testing.T) {
	devices := []Device{
		{Id: 0, Name: "device 1", Interval: 1, Value: 11},
		{Id: 1, Name: "device 2", Interval: 2, Value: 12},
		{Id: 2, Name: "device 3", Interval: 3, Value: 13},
		{Id: 3, Name: "device 4", Interval: 4, Value: 14},
		{Id: 4, Name: "device 5", Interval: 5, Value: 15},
	}
	dao := inMemoryDeviceDAO{devices: append([]Device(nil), devices...)}

	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	query := req.URL.Query()
	query.Add("limit", "2")
	query.Add("page", "1")
	req.URL.RawQuery = query.Encode()
	res := httptest.NewRecorder()

	underTest := deviceHTTPHandler{service: DeviceService{dao: &dao}}

	underTest.getAll(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	require.JSONEq(t, `[{"id":2,"name":"device 3","interval":3,"value":13},{"id":3,"name":"device 4","interval":4,"value":14}]`, res.Body.String())
}

type failingDeviceDAO struct {
}

func (db *failingDeviceDAO) Save(device Device) (Device, error) {
	return device, errors.New("mock error - failed to create device")
}

func (db *failingDeviceDAO) GetByID(_ int) (*Device, error) {
	return nil, errors.New(fmt.Sprintf("mock error - failed to get device by id"))
}

func (db *failingDeviceDAO) GetAll(_ int, _ int) ([]Device, error) {
	return nil, errors.New(fmt.Sprintf("mock error - failed to get all devices"))
}
