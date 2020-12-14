package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTickerService_Start_SendDeviceMeasurement(t *testing.T) {
	dao := inMemoryDeviceDAO{devices: []Device{{Id: 0, Interval: 1, Value: 5}}}
	ds := DeviceService{dao: &dao}
	measurements := make(chan Measurement)
	sendTrigger := make(chan time.Time)

	underTest := TickerService{ds: &ds, measurements: measurements, tf: func(d time.Duration) <-chan time.Time { return sendTrigger }}
	defer underTest.Stop()

	_ = underTest.Start()
	sendTrigger <- time.Time{}
	result := <-measurements

	assert.Equal(t, Measurement{Id: 0, Value: 5}, result)
}
