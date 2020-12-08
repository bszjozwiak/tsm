package main

import (
	"errors"
	"log"
	"sync"
	"time"
)

type TickerService struct {
	mu           sync.Mutex
	ds           *DeviceService
	measurements chan<- Measurement
	stop         chan bool
	isRunning    bool
}

func (ts *TickerService) Start() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.isRunning {
		log.Println("tickers started already")
		return nil
	}

	devices, err := ts.ds.GetAll(0, 0)
	if err != nil {
		return errors.New("failed to start measurements sending")
	}

	ts.stop = make(chan bool)
	ts.isRunning = true

	for _, device := range devices {
		ts.createTickerForDevice(device)
	}

	return nil
}

func (ts *TickerService) Stop() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.isRunning {
		log.Println("stopping send measurements")
		close(ts.stop)
		ts.isRunning = false
	}
}

func (ts *TickerService) NotifyDeviceCreated(device Device) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if !ts.isRunning {
		log.Println("measurements sending not run")
		return
	}

	ts.createTickerForDevice(device)
}

func (ts *TickerService) createTickerForDevice(device Device) {
	ticker := time.NewTicker(time.Second * time.Duration(device.Interval))

	go func(ticker *time.Ticker, deviceId int, value float64) {
		defer log.Printf("ticker for device %v stopped", deviceId)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ts.measurements <- Measurement{Id: deviceId, Value: value}
			case <-ts.stop:
				log.Printf("measurements sending from device %v stopped", deviceId)
				return
			}
		}
	}(ticker, device.Id, device.Value)
}
