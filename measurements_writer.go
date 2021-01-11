package main

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"log"
	"time"
)

type Measurement struct {
	Id    string
	Value float64
}

type MeasurementsWriter struct {
	measurements <-chan Measurement
	writeAPI     api.WriteAPIBlocking
}

func (mw *MeasurementsWriter) Start() {
	for measurement := range mw.measurements {
		point := influxdb2.NewPointWithMeasurement("deviceValues").
			AddTag("deviceId", measurement.Id).
			AddField("value", measurement.Value).
			SetTime(time.Now().Round(time.Second))

		if err := mw.writeAPI.WritePoint(context.Background(), point); err != nil {
			log.Print(err)
		}
	}
}
