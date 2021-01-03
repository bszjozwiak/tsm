package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestTickerService_Start_SendDeviceMeasurement(t *testing.T) {
	id := primitive.NewObjectID()
	dao := inMemoryDeviceDAO{devices: []Device{{ID: id, Interval: 1, Value: 5}}}
	ds := DeviceService{dao: &dao}
	measurements := make(chan Measurement)
	sendTrigger := make(chan time.Time)

	underTest := TickerService{
		ds:        &ds,
		publisher: func(id string, value float64) error { measurements <- Measurement{Id: id, Value: value}; return nil },
		tf:        func(d time.Duration) <-chan time.Time { return sendTrigger },
	}
	defer underTest.Stop()

	_ = underTest.Start(context.Background())
	sendTrigger <- time.Time{}
	result := <-measurements

	assert.Equal(t, Measurement{Id: id.Hex(), Value: 5}, result)
}
