package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDeviceDAO struct {
	devices *mongo.Collection
}

func (dao *mongoDeviceDAO) Save(device Device) (Device, error) {
	countDocuments, err := dao.devices.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return Device{}, err
	}

	device.Id = int(countDocuments)
	insertResult, err := dao.devices.InsertOne(context.Background(), device)
	if err != nil {
		return device, err
	}

	var result Device
	if err = dao.devices.FindOne(context.Background(), bson.M{"_id": insertResult.InsertedID}).Decode(&result); err != nil {
		return Device{}, err
	}

	return result, nil
}

func (dao *mongoDeviceDAO) GetByID(id int) (*Device, error) {
	var result Device
	if err := dao.devices.FindOne(context.Background(), bson.M{"id": id}).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func (dao *mongoDeviceDAO) GetAll(limit int, page int) ([]Device, error) {
	if limit < 0 {
		return nil, errors.New("limit can't be negative")
	}

	var filter bson.M

	if limit == 0 {
		filter = bson.M{}
	} else {
		start := limit * page
		end := start + limit
		filter = bson.M{"id": bson.M{"$gte": start, "$lt": end}}
	}

	find, err := dao.devices.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	devices := make([]Device, 0)
	for find.Next(context.Background()) {
		var dev Device
		if err = find.Decode(&dev); err != nil {
			return nil, err
		}

		devices = append(devices, dev)
	}

	return devices, nil
}
