package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"my_app/internal/config"
)

var db *mongo.Database

func GetDB() *mongo.Database {
	return db
}

func GetColl(coll string) *mongo.Collection {
	return GetDB().Collection(coll)
}

func InitDB() error {
	uri := config.MongoUri

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	db = client.Database("db")

	return nil
}

func CloseDB() error {
	return db.Client().Disconnect(ctx)
}
