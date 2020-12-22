package main

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"log"
	"strconv"
	"time"
)

type Measurement struct {
	Id    int
	Value float64
}

type MeasurementsWriter struct {
	measurements <-chan Measurement
	writeAPI     api.WriteAPIBlocking
}

func (mw *MeasurementsWriter) AsyncStart() {
	go func() {
		for {
			select {
			case measurement := <-mw.measurements:
				point := influxdb2.NewPointWithMeasurement("deviceValues").
					AddTag("deviceId", strconv.Itoa(measurement.Id)).
					AddField("value", measurement.Value).
					SetTime(time.Now().Round(time.Second))

				if err := mw.writeAPI.WritePoint(context.Background(), point); err != nil {
					log.Print(err)
				}
			}
		}
	}()
}
