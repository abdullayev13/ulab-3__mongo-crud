package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	OrderedBy primitive.ObjectID `json:"ordered_by" bson:"ordered_by"`
	Items     []OrderItem        `json:"items" bson:"-"`
	Comment   string             `json:"comment" bson:"comment"`
}

type OrderItem struct {
	OrderId   primitive.ObjectID `json:"order_id" bson:"order_id"`
	ProductId primitive.ObjectID `json:"product_id" bson:"product_id"`
	Quantity  int                `json:"quantity" bson:"quantity"`
}
