package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateValidDeviceWithNoDatabaseErrors(t *testing.T) {
	id := primitive.NewObjectID()
	body := strings.NewReader(fmt.Sprintf(`{"id":"%v","name":"test name2","interval":1}`, id.Hex()))
	req := httptest.NewRequest(http.MethodPost, "/devices", body)
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: &DeviceService{dao: &inMemoryDeviceDAO{}}}

	underTest.createDevice(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)
	require.JSONEq(t, fmt.Sprintf(`{"id":"%v","name":"test name2","interval":1,"value":0}`, id.Hex()), res.Body.String())
}

func TestCreateInvalidDevice(t *testing.T) {
	body := strings.NewReader(`{"name":"test name2","interval":-1}`)
	req := httptest.NewRequest(http.MethodPost, "/devices", body)
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: &DeviceService{dao: &inMemoryDeviceDAO{}}}

	underTest.createDevice(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestCreateValidDeviceButErrorWhenSaveInDatabase(t *testing.T) {
	body := strings.NewReader(`{"name":"test name2","interval":1}`)
	req := httptest.NewRequest(http.MethodPost, "/devices", body)
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: &DeviceService{dao: &failingDeviceDAO{}}}

	underTest.createDevice(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestGetByIDDatabaseError(t *testing.T) {
	req := createGetDeviceRequest("0")
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: &DeviceService{dao: &failingDeviceDAO{}}}

	underTest.getByID(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestGetByIdDeviceNotFound(t *testing.T) {
	req := createGetDeviceRequest("0")
	res := httptest.NewRecorder()
	underTest := deviceHTTPHandler{service: &DeviceService{dao: &inMemoryDeviceDAO{}}}

	underTest.getByID(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestGetByIDDeviceExists(t *testing.T) {
	id := primitive.NewObjectID()
	req := createGetDeviceRequest(id.Hex())
	res := httptest.NewRecorder()
	device := Device{ID: id, Name: "device", Interval: 1, Value: 1}
	dao := inMemoryDeviceDAO{}
	_, _ = dao.Save(context.Background(), device)
	underTest := deviceHTTPHandler{service: &DeviceService{dao: &dao}}

	underTest.getByID(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	require.JSONEq(t, fmt.Sprintf(`{"id":"%v","name":"device","interval":1,"value":1}`, id.Hex()), res.Body.String())
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

	underTest := deviceHTTPHandler{service: &DeviceService{dao: &inMemoryDeviceDAO{}}}

	for param, values := range testCases {
		for _, value := range values {
			t.Run(fmt.Sprintf("%v %v", param, value), func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/devices", nil)
				query := req.URL.Query()
				query.Add(param, value)
				req.URL.RawQuery = query.Encode()
				res := httptest.NewRecorder()

				underTest.getAll(res, req)

				assert.Equal(t, http.StatusBadRequest, res.Code)
			})
		}
	}
}

func TestGetAllDatabaseError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	res := httptest.NewRecorder()

	underTest := deviceHTTPHandler{service: &DeviceService{dao: &failingDeviceDAO{}}}

	underTest.getAll(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestGetAllNoDevicesReturnEmptyList(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	res := httptest.NewRecorder()

	underTest := deviceHTTPHandler{service: &DeviceService{dao: &inMemoryDeviceDAO{}}}

	underTest.getAll(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	require.JSONEq(t, "[]", res.Body.String())
}

func TestGetAllReturnRequestedDevices(t *testing.T) {
	devices := []Device{
		{Name: "device 1", Interval: 1, Value: 11},
		{Name: "device 2", Interval: 2, Value: 12},
		{Name: "device 3", Interval: 3, Value: 13},
		{Name: "device 4", Interval: 4, Value: 14},
		{Name: "device 5", Interval: 5, Value: 15},
	}
	dao := inMemoryDeviceDAO{devices: append([]Device(nil), devices...)}

	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	query := req.URL.Query()
	query.Add("limit", "2")
	query.Add("page", "1")
	req.URL.RawQuery = query.Encode()
	res := httptest.NewRecorder()

	underTest := deviceHTTPHandler{service: &DeviceService{dao: &dao}}

	underTest.getAll(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	require.JSONEq(t, fmt.Sprintf(`[{"id":"%[1]v","name":"device 3","interval":3,"value":13},{"id":"%[1]v","name":"device 4","interval":4,"value":14}]`, primitive.NilObjectID.Hex()), res.Body.String())
}

type failingDeviceDAO struct {
}

func (db *failingDeviceDAO) Save(_ context.Context, device Device) (Device, error) {
	return device, errors.New("mock error - failed to create device")
}

func (db *failingDeviceDAO) GetByID(_ context.Context, _ string) (*Device, error) {
	return nil, errors.New(fmt.Sprintf("mock error - failed to get device by id"))
}

func (db *failingDeviceDAO) GetAll(_ context.Context, _ int, _ int) ([]Device, error) {
	return nil, errors.New(fmt.Sprintf("mock error - failed to get all devices"))
}
