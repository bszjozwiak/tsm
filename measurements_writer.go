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

type receiverFactory func() (<-chan Measurement, error)

type MeasurementsWriter struct {
	rf       receiverFactory
	writeAPI api.WriteAPIBlocking
}

func (mw *MeasurementsWriter) AsyncStart() error {
	measurements, err := mw.rf()
	if err != nil {
		return err
	}

	go func() {
		for m := range measurements {
			point := influxdb2.NewPointWithMeasurement("deviceValues").
				AddTag("deviceId", m.Id).
				AddField("value", m.Value).
				SetTime(time.Now().Round(time.Second))

			if writeErr := mw.writeAPI.WritePoint(context.Background(), point); writeErr != nil {
				log.Print(writeErr)
			}
		}
	}()

	return nil
}
