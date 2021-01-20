package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDeviceDAO struct {
	db *mongo.Database
}

func (dao *mongoDeviceDAO) Save(ctx context.Context, device Device) (Device, error) {
	devices := dao.db.Collection("devices")
	device.ID = primitive.NewObjectID()
	insertResult, err := devices.InsertOne(ctx, device)
	if err != nil {
		return device, err
	}

	var result Device
	if err = devices.FindOne(ctx, bson.M{"_id": insertResult.InsertedID}).Decode(&result); err != nil {
		return Device{}, err
	}

	return result, nil
}

func (dao *mongoDeviceDAO) GetByID(ctx context.Context, id string) (*Device, error) {
	var result Device
	devices := dao.db.Collection("devices")

	searchId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if err := devices.FindOne(ctx, bson.M{"_id": searchId}).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func (dao *mongoDeviceDAO) GetAll(ctx context.Context, limit int, page int) ([]Device, error) {
	if limit < 0 {
		return nil, errors.New("limit can't be negative")
	}

	var filter bson.M

	opts := &options.FindOptions{}
	if limit > 0 {
		opts.SetSkip(int64(limit * page)).
			SetLimit(int64(limit))
	}

	devices := dao.db.Collection("devices")
	find, err := devices.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	result := make([]Device, 0)
	for find.Next(ctx) {
		var dev Device
		if err = find.Decode(&dev); err != nil {
			return nil, err
		}

		result = append(result, dev)
	}

	return result, nil
}
