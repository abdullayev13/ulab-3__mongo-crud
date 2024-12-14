package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

func InitIndies() error {

	coll := db.Collection("products")
	indexName, err := coll.Indexes().CreateOne(ctx,
		mongo.IndexModel{Keys: bson.D{{Key: "category_name", Value: 1}}},
	)
	if err != nil {
		return err
	}
	_ = indexName

	coll = db.Collection("orders")
	indexName, err = coll.Indexes().CreateOne(ctx,
		mongo.IndexModel{Keys: bson.D{{Key: "ordered_by", Value: 1}}},
	)
	if err != nil {
		return err
	}

	coll = db.Collection("order_items")
	indexName, err = coll.Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "order_id", Value: 1}, {Key: "product_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
