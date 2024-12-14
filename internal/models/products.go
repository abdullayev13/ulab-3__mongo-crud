package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Id           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Price        float64            `json:"price" bson:"price"`
	CategoryName string             `json:"category_name" bson:"category_name"`
}
