package main

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

type tickerFactory func(d time.Duration) <-chan time.Time
type measurementPublisher func(id string, value float64) error

type TickerService struct {
	mu        sync.Mutex
	ds        *DeviceService
	tf        tickerFactory
	publisher measurementPublisher
	stop      chan bool
	isRunning bool
}

func (ts *TickerService) Start(ctx context.Context) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.isRunning {
		log.Println("tickers started already")
		return nil
	}

	devices, err := ts.ds.GetAll(ctx, 0, 0)
	if err != nil {
		log.Print(err)
		return errors.New("failed to start measurements sending")
	}

	ts.stop = make(chan bool)
	ts.isRunning = true

	for _, device := range devices {
		go func(d Device) {
			ts.createTickerForDevice(d)
		}(device)
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
	defer log.Printf("ticker for device %v stopped", device.ID)

	for {
		select {
		case <-sendTrigger:
			if err := ts.publisher(device.ID.Hex(), device.Value); err != nil {
				log.Print(err)
			}
		case <-ts.stop:
			log.Printf("measurements sending from device %v stopped", device.ID)
			return
		}
	}
}
