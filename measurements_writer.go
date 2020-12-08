package main

import (
	"log"
	"sync"
)

type Measurement struct {
	Id    int
	Value float64
}

type MeasurementsWriter struct {
	mu           sync.Mutex
	measurements <-chan Measurement
	stop         chan bool
	isRunning    bool
}

func (mw *MeasurementsWriter) Start() {
	mw.mu.Lock()
	defer mw.mu.Unlock()

	if mw.isRunning {
		log.Println("measurements writer started already")
		return
	}

	mw.stop = make(chan bool)
	mw.isRunning = true

	go func() {
		for {
			select {
			case measurement := <-mw.measurements:
				log.Printf("measurement device: %v value: %v", measurement.Id, measurement.Value)
			case <-mw.stop:
				log.Println("measurements writing stopped")
				return
			}
		}
	}()
}

func (mw *MeasurementsWriter) Stop() {
	mw.mu.Lock()
	defer mw.mu.Unlock()

	if mw.isRunning {
		log.Println("stopping write measurements")
		close(mw.stop)
		mw.isRunning = false
	}
}
