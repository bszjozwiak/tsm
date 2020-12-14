package main

import (
	"errors"
	"log"
	"sync"
	"time"
)

type tickerFactory func(d time.Duration) <-chan time.Time

type TickerService struct {
	mu           sync.Mutex
	ds           *DeviceService
	tf           tickerFactory
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
		log.Print(err)
		return errors.New("failed to start measurements sending")
	}

	ts.stop = make(chan bool)
	ts.isRunning = true

	for _, device := range devices {
		go ts.createTickerForDevice(device)
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

	go ts.createTickerForDevice(device)
}

func (ts *TickerService) createTickerForDevice(device Device) {
	sendTrigger := ts.tf(time.Second * time.Duration(device.Interval))

	func(notifyTime <-chan time.Time, deviceId int, value float64) {
		defer log.Printf("ticker for device %v stopped", deviceId)

		for {
			select {
			case <-notifyTime:
				ts.measurements <- Measurement{Id: deviceId, Value: value}
			case <-ts.stop:
				log.Printf("measurements sending from device %v stopped", deviceId)
				return
			}
		}
	}(sendTrigger, device.Id, device.Value)
}
