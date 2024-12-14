package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Name     string             `json:"name" bson:"name"`
	Bio      string             `json:"bio" bson:"bio"`
}
