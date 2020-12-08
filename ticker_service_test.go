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
	defer close(measurements)

	underTest := TickerService{ds: &ds, measurements: measurements}
	defer underTest.Stop()

	timeout := time.After(3 * time.Second)
	done := make(chan bool)
	defer close(done)

	var result Measurement

	go func() {
		_ = underTest.Start()
		result = <-measurements
		done <- true
	}()

	select {
	case <-timeout:
		t.Fail()
	case <-done:
	}

	assert.Equal(t, Measurement{Id: 0, Value: 5}, result)
}
