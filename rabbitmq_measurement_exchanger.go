package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"strconv"
)

type rabbitMQMeasurementExchanger struct {
	measurements *amqp.Channel
}

func newRabbitMQMeasurementExchanger(url string) (*rabbitMQMeasurementExchanger, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	measurements, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = measurements.ExchangeDeclare(
		"measurements",
		"topic",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return nil, err
	}

	_, err = measurements.QueueDeclare(
		"measurements",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = measurements.QueueBind(
		"measurements",
		"*",
		"measurements",
		false,
		nil)
	if err != nil {
		return nil, err
	}

	return &rabbitMQMeasurementExchanger{measurements: measurements}, nil
}

func (rme *rabbitMQMeasurementExchanger) Publish(id string, value float64) error {
	return rme.measurements.Publish("measurements", id, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: []byte(fmt.Sprintf("%f", value))})
}

func (rme *rabbitMQMeasurementExchanger) CreateReceiver() (<-chan Measurement, error) {
	delivery, err := rme.measurements.Consume(
		"measurements",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	measurements := make(chan Measurement, 10)

	go func() {
		for d := range delivery {
			value, parseValErr := strconv.ParseFloat(string(d.Body), 64)
			if parseValErr != nil {
				log.Print(parseValErr)
				continue
			}

			measurements <- Measurement{Id: d.RoutingKey, Value: value}
		}
	}()

	return measurements, nil
}
