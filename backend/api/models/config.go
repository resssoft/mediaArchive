package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Config struct {
	ID      primitive.ObjectID `bson:"_id"`
	Name    string             `bson:"name"`
	Group   string             `bson:"group"`
	Perm    string             `bson:"perm"`
	Data    string             `bson:"data"`
	Comment string             `bson:"comment"`
}
